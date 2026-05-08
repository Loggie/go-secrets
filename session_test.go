package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	path := dbus.ObjectPath("/org/freedesktop/secrets/session/test")
	session := secrets.NewSession(conn, path)
	require.NotNil(t, session)
	require.Equal(t, path, session.Path())
}

func TestSessionClose(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	session := requireOpenSession(t, svc)
	require.NoError(t, session.Close())
}

func TestSessionCloseInvalidPathReturnsError(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	session := secrets.NewSession(conn, dbus.ObjectPath("/org/freedesktop/secrets/session/does-not-exist"))
	err = session.Close()
	require.Error(t, err)
}
