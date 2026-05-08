package secrets

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

type Prompt struct {
	*Object
}

func NewPrompt(
	conn *dbus.Conn,
	path dbus.ObjectPath,
) Prompt {
	return Prompt{
		Object: NewObject(conn, path),
	}
}

func (p *Prompt) Prompt() (any, error) {
	parentWindowID := ""

	c := make(chan *dbus.Signal, 1)
	defer close(c)

	p.conn.Signal(c)
	defer p.conn.RemoveSignal(c)

	err := p.Call(SecretPromptMethodPrompt, 0, parentWindowID).Err
	if err != nil {
		return nil, err
	}

	err = p.conn.AddMatchSignal(
		dbus.WithMatchInterface(SecretPromptInterface),
		dbus.WithMatchMember("Completed"),
	)
	if err != nil {
		return nil, err
	}

	res := <-c

	if res.Path != p.Path() {
		return nil, errors.New("received signal from unexpected path")
	}

	v, ok := res.Body[1].(dbus.Variant)
	if !ok {
		return nil, fmt.Errorf(
			"%w: in Prompt.Prompt expected dbus.Variant, got %T",
			ErrUnexpectedType,
			res.Body[1],
		)
	}

	return v.Value(), nil
}

func (p *Prompt) Dismiss() error {
	return p.Call(SecretPromptMethodDismiss, 0).Err
}
