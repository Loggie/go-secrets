package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

func TestServicePath(t *testing.T) {
	t.Parallel()

	require.Equal(t, dbus.ObjectPath(secrets.SecretServicePath), secrets.ServicePath())
}

func TestDefaultCollectionPath(t *testing.T) {
	t.Parallel()

	require.Equal(t, dbus.ObjectPath(secrets.SecretCollectionDefaultPath), secrets.DefaultCollectionPath())
}

func TestCollectionPath(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		dbus.ObjectPath(secrets.SecretCollectionBasePath+"/new_secret_name"),
		secrets.CollectionPath("new secret-name"),
	)
}

func TestSessionPath(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		dbus.ObjectPath(secrets.SecretServicePath+"/session/my_session"),
		secrets.SessionPath("my session"),
	)
}

func TestPromptPath(t *testing.T) {
	t.Parallel()

	require.Equal(t, dbus.ObjectPath(secrets.SecretServicePath+"/prompt/my_prompt"), secrets.PromptPath("my prompt"))
}

func TestItemPath(t *testing.T) {
	t.Parallel()

	collection := secrets.CollectionPath("my_collection")
	require.Equal(t, dbus.ObjectPath(string(collection)+"/item_1"), secrets.ItemPath(collection, "item-1"))
}

func TestCollectionPathEmptyFallsBack(t *testing.T) {
	t.Parallel()

	require.Equal(t, dbus.ObjectPath(secrets.SecretCollectionBasePath+"/_"), secrets.CollectionPath("   "))
}
