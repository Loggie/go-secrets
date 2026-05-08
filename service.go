package secrets

import (
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

type Service struct {
	*Object
}

func NewService() (*Service, error) {
	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		return nil, errors.New("failed to connect to dbus")
	}

	return &Service{
		Object: NewObject(conn, SecretServicePath),
	}, nil
}

func (s *Service) Collections() ([]*Collection, error) {
	v, err := s.GetProperty(SecretServicePropertyCollections)
	if err != nil {
		return nil, err
	}

	collections, ok := v.Value().([]dbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf(
			"%w: in Service.Collections expected []dbus.ObjectPath, got %T",
			ErrUnexpectedType,
			v.Value(),
		)
	}

	var collectionList []*Collection
	for _, collection := range collections {
		collectionList = append(collectionList, NewCollection(s, collection))
	}

	return collectionList, nil
}

func (s *Service) OpenSession(
	algorithm string,
	args string,
) (*Session, error) {
	var output dbus.Variant
	var sessionPath dbus.ObjectPath

	err := s.Call(SecretServiceMethodOpenSession, 0,
		algorithm,
		dbus.MakeVariant(args),
	).Store(
		&output, &sessionPath,
	)

	if err != nil {
		return nil, err
	}

	return NewSession(s.conn, sessionPath), nil
}

func (s *Service) CreateCollection(
	label string,
	alias string,
) (*Collection, error) {
	var collectionPath dbus.ObjectPath
	var promptPath dbus.ObjectPath

	properties := map[string]any{
		SecretCollectionPropertyLabel:    label,
		SecretCollectionPropertyCreated:  uint64(time.Now().Unix()), //gosec:disable G115
		SecretCollectionPropertyModified: uint64(time.Now().Unix()), //gosec:disable G115
	}

	err := s.Call(SecretServiceMethodCreateCollection, 0,
		properties, alias,
	).Store(&collectionPath, &promptPath)

	if err != nil {
		return nil, err
	}

	if promptPath != "/" {
		prompt := NewPrompt(s.conn, promptPath)
		result, errPrompt := prompt.Prompt()
		if errPrompt != nil {
			return nil, errPrompt
		}

		createdPath, ok := result.(dbus.ObjectPath)
		if !ok {
			return nil, fmt.Errorf(
				"%w: in Service.CreateCollection expected dbus.ObjectPath, got %T",
				ErrUnexpectedType,
				result,
			)
		}
		collectionPath = createdPath
	}

	return NewCollection(s, collectionPath), nil
}

func (s *Service) SearchItems(
	attrs Attributes,
) ([]*Item, []*Item, error) {
	var unlockedPaths []dbus.ObjectPath
	var lockedPaths []dbus.ObjectPath

	err := s.Call(SecretServiceMethodSearchItems, 0,
		attrs,
	).Store(&unlockedPaths, &lockedPaths)

	if err != nil {
		return nil, nil, err
	}

	unlocked := make([]*Item, 0, len(unlockedPaths))
	for _, path := range unlockedPaths {
		unlocked = append(unlocked, NewItem(s.conn, path))
	}

	locked := make([]*Item, 0, len(lockedPaths))
	for _, path := range lockedPaths {
		locked = append(locked, NewItem(s.conn, path))
	}

	return unlocked, locked, nil
}

func (s *Service) Unlock(
	objects []LockableObject,
) error {
	var unlocked []dbus.ObjectPath
	var promptPath dbus.ObjectPath

	paths := make([]dbus.ObjectPath, 0, len(objects))
	for _, object := range objects {
		paths = append(paths, object.Path())
	}

	err := s.Call(SecretServiceMethodUnlock, 0,
		paths,
	).Store(
		&unlocked, &promptPath,
	)

	if err != nil {
		return err
	}

	if promptPath != "/" {
		prompt := NewPrompt(s.conn, promptPath)
		_, errPrompt := prompt.Prompt()
		if errPrompt != nil {
			return errPrompt
		}
	}

	return nil
}

func (s *Service) Lock(
	objects []LockableObject,
) error {
	var locked []dbus.ObjectPath
	var promptPath dbus.ObjectPath

	paths := make([]dbus.ObjectPath, 0, len(objects))
	for _, object := range objects {
		paths = append(paths, object.Path())
	}

	err := s.Call(SecretServiceMethodLock, 0,
		paths,
	).Store(
		&locked, &promptPath,
	)

	if err != nil {
		return err
	}

	if promptPath != "/" {
		prompt := NewPrompt(s.conn, promptPath)
		_, errPrompt := prompt.Prompt()
		if errPrompt != nil {
			return errPrompt
		}
	}

	return err
}

func (s *Service) GetSecrets(
	items []*Item,
	session *Session,
) (map[string]*Secret, error) {
	var itemPaths []dbus.ObjectPath
	for _, item := range items {
		itemPaths = append(itemPaths, item.Path())
	}

	var secrets map[dbus.ObjectPath]*Secret

	err := s.Call(SecretServiceMethodGetSecrets, 0,
		itemPaths, session.Path(),
	).Store(&secrets)

	if err != nil {
		return nil, err
	}

	result := make(map[string]*Secret)
	for itemPath, secret := range secrets {
		result[string(itemPath)] = secret
	}
	return result, nil
}

func (s *Service) ReadAlias(name string) (*Collection, error) {
	var collectionPath dbus.ObjectPath

	err := s.Call(SecretServiceMethodReadAlias, 0,
		name,
	).Store(&collectionPath)

	if err != nil {
		return nil, err
	}

	return NewCollection(s, collectionPath), nil
}

func (s *Service) SetAlias(
	collection Collection,
	alias string,
) error {
	err := s.Call(SecretServiceMethodSetAlias, 0,
		alias,
		collection.Path(),
	).Err

	if err != nil {
		return err
	}

	return nil
}
