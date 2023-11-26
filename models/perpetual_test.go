package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPerpetual(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	path := "/tmp/test.sqlite3"
	store, err := OpenDataSQLite3Store(path)
	require.Nil(err)
	defer os.Remove(path)

	p, err := store.CreatePerpetual(ctx, "bitcoin")
	require.Nil(err)
	require.NotNil(p)
	p, err = store.ReadPerpetual(ctx, "bitcoin")
	require.Nil(err)
	require.NotNil(p)
}
