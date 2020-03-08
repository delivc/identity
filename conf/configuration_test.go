package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	defer os.Clearenv()
	os.Exit(m.Run())
}

func TestGlobal(t *testing.T) {
	os.Setenv("DELIVC_DB_DRIVER", "pgsql")
	os.Setenv("DELIVC_DATABASE_URL", "fake")
	os.Setenv("DELIVC_OPERATOR_TOKEN", "token")
	os.Setenv("DELIVC_API_REQUEST_ID_HEADER", "X-Request-ID")
	gc, err := LoadGlobal("")
	require.NoError(t, err)
	require.NotNil(t, gc)
	assert.Equal(t, "X-Request-ID", gc.API.RequestIDHeader)
}
