package secrets

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type Item struct {
	*Object
}

func NewItem(
	conn *dbus.Conn,
	path dbus.ObjectPath,
) *Item {
	return &Item{
		Object: NewObject(conn, path),
	}
}

func (i *Item) Locked() (bool, error) {
	v, err := i.GetProperty(SecretItemPropertyLocked)

	if err != nil {
		return false, err
	}

	val, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf(
			"%w: in Item.Locked expected bool, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return val, nil
}

func (i *Item) Attributes() (Attributes, error) {
	v, err := i.GetProperty(SecretItemPropertyAttributes)

	if err != nil {
		return nil, err
	}

	switch attrs := v.Value().(type) {
	case Attributes:
		return attrs, nil
	case map[string]string:
		return Attributes(attrs), nil
	default:
		return nil, fmt.Errorf(
			"%w: in Item.Attributes expected Attributes or map[string]string, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}
}

func (i *Item) Label() (string, error) {
	v, err := i.GetProperty(SecretItemPropertyLabel)

	if err != nil {
		return "", err
	}

	val, ok := v.Value().(string)
	if !ok {
		return "", fmt.Errorf(
			"%w: in Item.Label expected string, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return val, nil
}

func (i *Item) Created() (time.Time, error) {
	v, err := i.GetProperty(SecretItemPropertyCreated)

	if err != nil {
		return time.Time{}, err
	}

	val, ok := v.Value().(uint64)
	if !ok {
		return time.Time{}, fmt.Errorf(
			"%w: in Item.Created expected uint64, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return time.Unix(int64(val), 0), nil //gosec:disable G115
}

func (i *Item) Modified() (time.Time, error) {
	v, err := i.GetProperty(SecretItemPropertyModified)

	if err != nil {
		return time.Time{}, err
	}

	val, ok := v.Value().(uint64)
	if !ok {
		return time.Time{}, fmt.Errorf(
			"%w: in Item.Modified expected uint64, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return time.Unix(int64(val), 0), nil //gosec:disable G115
}

func (i *Item) Delete() (Prompt, error) {
	var promptPath dbus.ObjectPath

	err := i.Call(SecretItemMethodDelete, 0).Store(&promptPath)
	if err != nil {
		return Prompt{}, err
	}

	return NewPrompt(i.conn, promptPath), nil
}

func (i *Item) GetSecret(session Session) (Secret, error) {
	var secret Secret

	err := i.Call(SecretItemMethodGetSecret, 0,
		session.Path(),
	).Store(
		&secret,
	)

	if err != nil {
		if strings.Contains(err.Error(), "org.freedesktop.Secret.Error.IsLocked") {
			return Secret{}, ErrSecretLocked
		}
		return Secret{}, err
	}

	return secret, nil
}

func (i *Item) SetSecret(secret Secret) error {
	err := i.Call(SecretItemMethodSetSecret, 0, secret).Err

	if err != nil {
		return err
	}

	return nil
}
