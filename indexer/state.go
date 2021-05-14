package indexer

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/umbracle/eth-indexer/indexer/proto"
	"github.com/umbracle/eth-indexer/sdk"
	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
	"github.com/umbracle/go-web3"
)

type State struct {
	db *sqlx.DB
	i  *Server
}

func newState(path string) (*State, error) {
	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, err
	}
	return newStateWithDB(db)
}

func newStateWithDB(db *sqlx.DB) (*State, error) {
	s := &State{
		db: db,
	}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *State) migrate() error {
	for _, migration := range AssetNames() {
		if _, err := s.db.Exec(string(MustAsset(migration))); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) GetTrackByName(name string) (*proto.Track, error) {
	var track proto.Track
	if err := s.db.Get(&track, "SELECT * FROM tracks WHERE name = $1", name); err != nil {
		return nil, err
	}
	return &track, nil
}

func (s *State) GetTracks() ([]proto.Track, error) {
	var tracks []proto.Track
	if err := s.db.Select(&tracks, "SELECT * FROM tracks"); err != nil {
		return nil, err
	}
	return tracks, nil
}

func (s *State) UpsertTrack(t *proto.Track) error {
	// safe check
	if _, err := filterConfigFromTracker(t); err != nil {
		return err
	}
	_, err := s.db.NamedExec("INSERT INTO tracks (name, to_addr, from_addr, topic, startblock, lastblocknum, lastblockhash, synced) VALUES (:name, :to_addr, :from_addr, :topic, :startblock, 0, '', false) ON CONFLICT DO NOTHING", t)
	if err != nil {
		return err
	}
	return nil
}

func (s *State) UpdateTrackSyncBlock(name string, blockNum uint64, blockHash web3.Hash) error {
	if _, err := s.db.Exec("UPDATE tracks SET lastblockhash = $1, lastblocknum = $2 WHERE name = $3", blockHash.String(), blockNum, name); err != nil {
		return err
	}
	return nil
}

func (s *State) UpdateTrackSynced(name string, synced bool) error {
	if _, err := s.db.Exec("UPDATE tracks SET synced = $1 WHERE name = $2", synced, name); err != nil {
		return err
	}
	return nil
}

const (
	AscOrder  = "asc"
	DescOrder = "desc"
)

type WhereCond string

const (
	WhereCondEqual = "equal"
)

type QueryWhere struct {
	Key   string
	Val   string
	Where WhereCond
}

type Query struct {
	Table   string
	First   uint64
	Skip    uint64
	OrderBy string
	Order   string
	Where   []QueryWhere
}

type ResObj struct {
	Data map[string]string
}

func (s *State) GetObjs(q *Query) ([]*ResObj, error) {
	sch, ok := s.i.schemas[q.Table]
	if !ok {
		return nil, fmt.Errorf("table %s not found", q.Table)
	}

	res := []*ResObj{}

	query := "SELECT * FROM " + q.Table
	if q.First != 0 {
		query += " LIMIT " + strconv.Itoa(int(q.First))
	}
	if q.Skip != 0 {
		query += " OFFSET " + strconv.Itoa(int(q.Skip))
	}
	if q.OrderBy != "" {
		query += " ORDER BY " + q.OrderBy
	}
	if q.Order != AscOrder {
		query += " DESC"
	}
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		obj, err := s.decodeObj(rows, sch)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func (s *State) GetObj2(table string, keys map[string]string) (*ResObj, error) {
	sch, ok := s.i.schemas[table]
	if !ok {
		return nil, fmt.Errorf("table not found %s", table)
	}

	kv := []string{}
	for k, v := range keys {
		kv = append(kv, fmt.Sprintf("%s = '%s'", k, v))
	}
	query := "SELECT * FROM " + table + " WHERE " + strings.Join(kv, " AND ")

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	defer rows.Close()

	obj, err := s.decodeObj(rows, sch)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *State) GetObj(table, k, v string) (*ResObj, error) {
	sch, ok := s.i.schemas[table]
	if !ok {
		return nil, fmt.Errorf("table %s not found", table)
	}

	rows, err := s.db.Query("SELECT * FROM "+table+" WHERE "+k+" = $1", v)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	obj, err := s.decodeObj(rows, sch)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *State) decodeObj(rows *sql.Rows, table *sdk.Table) (*ResObj, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(interface{})
	}
	if err := rows.Scan(values...); err != nil {
		return nil, err
	}
	for i := range cols {
		values[i] = *(values[i].(*interface{}))
	}

	obj := &ResObj{
		Data: map[string]string{},
	}
	for i := 0; i < len(cols); i++ {
		if values[i] == nil {
			continue
		}
		val := values[i]

		switch table.Fields[i].Type {
		case sdk.TypeAddress:
			// TODO: Validate address
			obj.Data[cols[i]] = val.(string)

		case sdk.TypeUint:
			raw := string(val.([]byte))
			// make sure its an int
			if _, err := strconv.Atoi(raw); err != nil {
				return nil, err
			}
			obj.Data[cols[i]] = raw

		case sdk.TypeDecimal:
			raw := string(val.([]byte))

			// make sure its a float
			f := new(sdk.Float)
			if !f.SetString(raw) {
				return nil, fmt.Errorf("incorrect float")
			}
			obj.Data[cols[i]] = raw

		default:
			return nil, fmt.Errorf("type not found")
		}
	}
	return obj, nil
}

func (s *State) UpsertTable(t *sdk.Table) error {
	// create the ddl
	ddl := buildDDL(t)

	if _, err := s.db.Exec(ddl); err != nil {
		return err
	}
	return nil
}

func (s *State) ApplyDiff(obj []*protosdk.Diff, apply bool) error {

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, diff := range obj {
		var query string
		if diff.Creation {
			// insert op
			names := []string{}
			vals := []string{}

			for k, v := range diff.Keys {
				names = append(names, k)
				vals = append(vals, "'"+v+"'")
			}
			for k, v := range diff.Vals {
				names = append(names, k)
				vals = append(vals, "'"+v+"'")
			}
			query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", diff.Table, strings.Join(names, ", "), strings.Join(vals, ", "))
		} else {
			// update op
			vals := []string{}
			for k, v := range diff.Vals {
				vals = append(vals, fmt.Sprintf("%s = '%s'", k, v))
			}
			where := []string{}
			for k, v := range diff.Keys {
				where = append(where, fmt.Sprintf("%s = '%s'", k, v))
			}
			query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", diff.Table, strings.Join(vals, ", "), strings.Join(where, " AND "))
		}

		if apply {
			if _, err := txn.Exec(query); err != nil {
				return err
			}
		}
	}

	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func buildDDL(t *sdk.Table) string {
	idFields := []string{}
	fieldNames := []string{}
	for _, f := range t.Fields {
		var typ string
		switch f.Type {
		case sdk.TypeAddress:
			typ = "text"
		case sdk.TypeUint:
			typ = "numeric"
		case sdk.TypeDecimal:
			typ = "decimal"
		default:
			panic(fmt.Sprintf("Not found: %d", f.Type))
		}
		if f.ID {
			idFields = append(idFields, f.Name)
		}
		fieldNames = append(fieldNames, fmt.Sprintf("%s %s", f.Name, typ))
	}

	// add the id fields
	if len(idFields) != 0 {
		fieldNames = append(fieldNames, fmt.Sprintf("UNIQUE (%s)", strings.Join(idFields, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", t.Name, strings.Join(fieldNames, ", "))
}
