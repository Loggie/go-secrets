package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

func TestNewPrompt(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	path := secrets.PromptPath("test")
	prompt := secrets.NewPrompt(conn, path)
	require.Equal(t, path, prompt.Path())
}

func TestPromptPromptInvalidPathReturnsError(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	prompt := secrets.NewPrompt(conn, secrets.PromptPath("does-not-exist"))
	_, err = prompt.Prompt()
	require.Error(t, err)
}

func TestPromptDismissInvalidPathReturnsError(t *testing.T) {
	t.Parallel()

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	prompt := secrets.NewPrompt(conn, secrets.PromptPath("does-not-exist"))
	err = prompt.Dismiss()
	require.Error(t, err)
}
