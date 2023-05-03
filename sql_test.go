package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQL(t *testing.T) {
	raw, err := ioutil.ReadFile("sql.json")
	require.NoError(t, err)

	type testData struct {
		Feature string
		ID      string
		SQL     []string
	}

	var testSet []testData
	require.NoError(t, json.Unmarshal(raw, &testSet))

	opts := immudb.DefaultOptions().WithAddress("localhost").WithPort(3322)
	client := immudb.NewClient().WithOptions(opts)
	require.NoError(t, client.OpenSession(context.TODO(), []byte("immudb"), []byte("immudb"), "defaultdb"))
	defer client.CloseSession(context.TODO())

	fmt.Printf("Test set count: %d\n", len(testSet))
	for _, td := range testSet {
		t.Run(td.Feature, func(t *testing.T) {
			td := td
			for _, sql := range td.SQL {
				_, err := client.SQLExec(context.TODO(), sql, nil)
				assert.NoError(t, err)
			}
		})
	}
}
