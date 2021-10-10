package state

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/umbracle/eth-indexer/schema"
)

type PostgresqlState struct {
	db *sqlx.DB
}

func newState(path string) (*PostgresqlState, error) {
	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, err
	}
	return newStateWithDB(db)
}

func newStateWithDB(db *sqlx.DB) (*PostgresqlState, error) {
	s := &PostgresqlState{
		db: db,
	}
	return s, nil
}

func (s *PostgresqlState) getSchema(table string) (*schema.Table, error) {
	// read data from db and cache
	return &schema.Table{}, nil
}

func (s *PostgresqlState) GetObjs(q *Query) ([]*Obj, error) {
	sch, err := s.getSchema(q.Table)
	if err != nil {
		return nil, err
	}

	res := []*Obj{}

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

func (s *PostgresqlState) GetObj(table string, keys map[string]string) (*Obj, error) {
	sch, err := s.getSchema(table)
	if err != nil {
		return nil, err
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

func (s *PostgresqlState) decodeObj(rows *sql.Rows, table *schema.Table) (*Obj, error) {
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

	obj := &Obj{
		Data: map[string]string{},
	}
	for i := 0; i < len(cols); i++ {
		if values[i] == nil {
			continue
		}
		val := values[i]

		switch table.Fields[i].Type {
		case schema.TypeAddress:
			// TODO: Validate address
			obj.Data[cols[i]] = val.(string)

		case schema.TypeUint:
			raw := string(val.([]byte))
			// make sure its an int
			if _, err := strconv.Atoi(raw); err != nil {
				return nil, err
			}
			obj.Data[cols[i]] = raw

		case schema.TypeDecimal:
			raw := string(val.([]byte))

			// make sure its a float
			f := new(schema.Float)
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

func (s *PostgresqlState) UpsertTable(t *schema.Table) error {
	// create the ddl
	ddl := buildDDL(t)

	if _, err := s.db.Exec(ddl); err != nil {
		return err
	}
	return nil
}

func (s *PostgresqlState) ApplyDiff(obj []*schema.Diff, apply bool) error {

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

func buildDDL(t *schema.Table) string {
	idFields := []string{}
	fieldNames := []string{}
	for _, f := range t.Fields {
		var typ string
		switch f.Type {
		case schema.TypeAddress:
			typ = "text"
		case schema.TypeUint:
			typ = "numeric"
		case schema.TypeDecimal:
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
