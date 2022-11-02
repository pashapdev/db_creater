package repo_test

import (
	"context"
	"testing"

	dbCreater "github.com/pashapdev/db_creater"
	"github.com/pashapdev/db_creater/examples/repo"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	const (
		user     = "postgres"
		password = "postgres"
		address  = "localhost"
		port     = 5432
		db       = "db_test"
	)
	ctx := context.Background()

	creater := dbCreater.New(user, password, address, db, port)
	testDb, err := creater.CreateWithMigration("file://./migrations/")
	require.NoError(t, err)
	defer creater.Drop(testDb)

	r, err := repo.New(user, password, address, testDb, port)
	require.NoError(t, err)
	defer r.Close(ctx)

	entities, err := r.SelectContent(ctx)
	require.NoError(t, err)
	require.Len(t, entities, 0)

	testData := []repo.TestEntity{
		{
			Content: "Content1",
		},
		{
			Content: "Content2",
		},
		{
			Content: "Content3",
		},
	}
	err = r.InsertContent(ctx, testData)
	require.NoError(t, err)

	entities, err = r.SelectContent(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, entities, testData)
}
