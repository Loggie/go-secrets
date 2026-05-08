package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

func TestNewObject(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		if errClose := conn.Close(); errClose != nil {
			t.Fatal(errClose)
		}
	}()

	path := secrets.CollectionPath("test")
	obj := secrets.NewObject(conn, path)
	require.NotNil(t, obj)
	require.Equal(t, path, obj.Path())
}
