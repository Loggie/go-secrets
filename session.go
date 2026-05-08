package secrets

import "github.com/godbus/dbus/v5"

type Session struct {
	*Object
}

func NewSession(
	conn *dbus.Conn,
	path dbus.ObjectPath,
) *Session {
	return &Session{
		Object: NewObject(conn, path),
	}
}

func (s *Session) Close() error {
	return s.Call(SecretSessionMethodClose, 0).Err
}
