package secrets

import "github.com/godbus/dbus/v5"

type Secret struct {
	Path        dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string
}
