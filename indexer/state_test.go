package indexer

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/umbracle/eth-indexer/indexer/proto"
	"github.com/umbracle/eth-indexer/sdk"
	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
	"github.com/umbracle/go-web3"
)

func setupPostgresql(t *testing.T) (*sqlx.DB, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_HOST_AUTH_METHOD=trust"})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	endpoint := fmt.Sprintf("postgres://postgres@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp"))

	// wait for the db to be running
	var db *sqlx.DB
	if err := pool.Retry(func() error {
		db, err = sqlx.Open("postgres", endpoint)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	}
	return db, cleanup
}

func TestState_SchemaDDL(t *testing.T) {
	db, close := setupPostgresql(t)
	defer close()

	s, err := newStateWithDB(db)
	assert.NoError(t, err)

	tb := &sdk.Table{
		Name: "a",
		Fields: []*sdk.Field{
			{
				Name: "a",
				Type: sdk.TypeAddress,
				ID:   true,
			},
			{
				Name: "b",
				Type: sdk.TypeAddress,
				ID:   true,
			},
			{
				Name: "c",
				Type: sdk.TypeAddress,
			},
		},
	}

	// validate that we can create the table
	assert.NoError(t, s.UpsertTable(tb))
}

func TestState_Diff(t *testing.T) {
	db, close := setupPostgresql(t)
	defer close()

	s, err := newStateWithDB(db)
	assert.NoError(t, err)

	tb := &sdk.Table{
		Name: "tname",
		Fields: []*sdk.Field{
			{
				Name: "a",
				Type: sdk.TypeAddress,
			},
			{
				Name: "b",
				Type: sdk.TypeUint,
			},
		},
	}
	assert.NoError(t, s.UpsertTable(tb))

	diff := []*protosdk.Diff{
		{
			Creation: true,
			Table:    "tname",
			Vals: map[string]string{
				"a": "b",
				"b": "1",
			},
		},
		{
			Table: "tname",
			Keys: map[string]string{
				"a": "b",
			},
			Vals: map[string]string{
				"b": "2",
			},
		},
	}
	assert.NoError(t, s.ApplyDiff(diff, true))
}

func TestState_Track(t *testing.T) {
	db, close := setupPostgresql(t)
	defer close()

	s, err := newStateWithDB(db)
	assert.NoError(t, err)

	t0 := &proto.Track{
		Name:       "track0",
		StartBlock: 2525,
	}
	assert.NoError(t, s.UpsertTrack(t0))

	t1, err := s.GetTrackByName("track0")
	assert.NoError(t, err)

	assert.Equal(t, t0, t1)

	// update block
	assert.NoError(t, s.UpdateTrackSyncBlock("track0", 1000, web3.Hash{0x1}))

	t2, err := s.GetTrackByName("track0")
	assert.NoError(t, err)
	assert.Equal(t, t2.LastBlockNum, uint64(1000))
}
