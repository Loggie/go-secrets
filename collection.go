package secrets

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

type Collection struct {
	*Object
}

func NewCollection(
	svc *Service,
	path dbus.ObjectPath,
) *Collection {
	return &Collection{
		Object: NewObject(svc.conn, path),
	}
}

func (c *Collection) Items() ([]*Item, error) {
	v, err := c.GetProperty(SecretCollectionPropertyItems)
	if err != nil {
		return nil, err
	}

	objects, ok := v.Value().([]dbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf(
			"%w: in Collection.Items expected []dbus.ObjectPath, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	items := make([]*Item, 0, len(objects))
	for _, item := range objects {
		items = append(items, NewItem(c.conn, item))
	}

	return items, nil
}

func (c *Collection) Label() (string, error) {
	v, err := c.GetProperty(SecretCollectionPropertyLabel)

	if err != nil {
		return "", err
	}

	val, ok := v.Value().(string)
	if !ok {
		return "", fmt.Errorf(
			"%w: in Collection.Label expected string, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return val, nil
}

func (c *Collection) SetLabel(label string) error {
	return c.SetProperty(SecretCollectionPropertyLabel, label)
}

func (c *Collection) Locked() (bool, error) {
	v, err := c.GetProperty(SecretCollectionPropertyLocked)

	if err != nil {
		return false, err
	}

	val, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf(
			"%w: in Collection.Locked expected bool, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return val, nil
}

func (c *Collection) Created() (time.Time, error) {
	v, err := c.GetProperty(SecretCollectionPropertyCreated)

	if err != nil {
		return time.Time{}, err
	}

	val, ok := v.Value().(uint64)
	if !ok {
		return time.Time{}, fmt.Errorf(
			"%w: in Collection.Created expected uint64, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return time.Unix(int64(val), 0), nil //gosec:disable G115
}

func (c *Collection) Modified() (time.Time, error) {
	v, err := c.GetProperty(SecretCollectionPropertyModified)
	if err != nil {
		return time.Time{}, err
	}

	val, ok := v.Value().(uint64)
	if !ok {
		return time.Time{}, fmt.Errorf(
			"%w: in Collection.Modified expected uint64, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	return time.Unix(int64(val), 0), nil //gosec:disable G115
}

func (c *Collection) Delete() (Prompt, error) {
	var promptPath dbus.ObjectPath

	err := c.Call(SecretCollectionMethodDelete, 0).Store(&promptPath)
	if err != nil {
		return Prompt{}, err
	}

	return NewPrompt(c.conn, promptPath), nil
}

func (c *Collection) SearchItems(
	attrs Attributes,
) ([]*Item, error) {
	var itemPaths []dbus.ObjectPath

	err := c.Call(SecretCollectionMethodSearchItems, 0,
		attrs,
	).Store(&itemPaths)

	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0, len(itemPaths))
	for _, item := range itemPaths {
		items = append(items, NewItem(c.conn, item))
	}

	return items, nil
}

func (c *Collection) CreateItem(
	label string,
	attrs Attributes,
	secret Secret,
	replace bool,
) (*Item, error) {
	var itemPath dbus.ObjectPath
	var promptPath dbus.ObjectPath

	err := c.Call(SecretCollectionMethodCreateItem, 0, map[string]any{
		SecretItemPropertyLabel:      label,
		SecretItemPropertyAttributes: attrs,
	}, secret, replace,
	).Store(
		&itemPath, &promptPath,
	)

	if err != nil {
		return nil, err
	}

	item := NewItem(c.conn, itemPath)

	if promptPath != "/" {
		prompt := NewPrompt(c.conn, promptPath)
		result, errPrompt := prompt.Prompt()
		if errPrompt != nil {
			return nil, errPrompt
		}

		createdPath, ok := result.(dbus.ObjectPath)
		if !ok {
			return nil, fmt.Errorf(
				"%w: in Collection.CreateItem expected dbus.ObjectPath, got %T",
				ErrUnexpectedType,
				result,
			)
		}
		itemPath = createdPath
	}

	return item, nil
}
