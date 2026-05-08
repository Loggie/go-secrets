package secrets

import (
	"strings"
	"unicode"

	"github.com/godbus/dbus/v5"
)

func ServicePath() dbus.ObjectPath {
	return dbus.ObjectPath(SecretServicePath)
}

func CollectionPath(name string) dbus.ObjectPath {
	return dbus.ObjectPath(SecretCollectionBasePath + "/" + objectPathSegment(name))
}

func DefaultCollectionPath() dbus.ObjectPath {
	return dbus.ObjectPath(SecretCollectionDefaultPath)
}

func SessionPath(id string) dbus.ObjectPath {
	return dbus.ObjectPath(SecretServicePath + "/session/" + objectPathSegment(id))
}

func PromptPath(id string) dbus.ObjectPath {
	return dbus.ObjectPath(SecretServicePath + "/prompt/" + objectPathSegment(id))
}

func ItemPath(collection dbus.ObjectPath, id string) dbus.ObjectPath {
	return dbus.ObjectPath(string(collection) + "/" + objectPathSegment(id))
}

func objectPathSegment(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "_"
	}

	out := make([]rune, 0, len(v))
	for _, r := range v {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			out = append(out, r)
			continue
		}
		out = append(out, '_')
	}

	return string(out)
}
