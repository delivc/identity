package models

import (
	"testing"

	"github.com/delivc/identity/conf"
	"github.com/delivc/identity/storage"
	"github.com/delivc/identity/storage/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const modelsTestConfig = "../hack/test.env"

type UserTestSuite struct {
	suite.Suite
	db *storage.Connection
}

func (ts *UserTestSuite) SetupTest() {
	TruncateAll(ts.db)
}

func TestUser(t *testing.T) {
	globalConfig, err := conf.LoadGlobal(modelsTestConfig)
	require.NoError(t, err)

	conn, err := test.SetupDBConnection(globalConfig)
	require.NoError(t, err)

	ts := &UserTestSuite{
		db: conn,
	}
	defer ts.db.Close()

	suite.Run(t, ts)
}
