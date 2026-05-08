package secrets

import "github.com/godbus/dbus/v5"

type Object struct {
	dbus.BusObject

	conn *dbus.Conn
}

type LockableObject interface {
	Locked() (bool, error)
	Path() dbus.ObjectPath
}

func NewObject(
	conn *dbus.Conn,
	path dbus.ObjectPath,
) *Object {
	return &Object{
		conn:      conn,
		BusObject: conn.Object(SecretService, path),
	}
}

func (c Object) Path() dbus.ObjectPath {
	return c.BusObject.Path()
}
