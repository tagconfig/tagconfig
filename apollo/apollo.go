package apollo

import (
	"reflect"
	"strings"

	"github.com/shima-park/agollo"
	"github.com/tagconfig/tagconfig"
)

var (
	defaultNamespace = "application"
)

type Client struct {
	agollo.Agollo
	namespaces []string
}

// NewClient returns a client
func NewClient(c agollo.Agollo, appid string, namespaces []string) *Client {
	return &Client{
		Agollo:     c,
		namespaces: namespaces,
	}
}

func (c *Client) Paires() (paires []tagconfig.Paire, err error) {
	for i := 0; i < len(c.namespaces); i++ {
		properties := c.GetNameSpace(c.namespaces[i])
		for k, v := range properties {
			paires = append(paires, tagconfig.Paire{
				Namespace: c.namespaces[i],
				Key:       k,
				Value:     v.(string),
			})
		}
	}
	return
}

func (c *Client) FieldInfo(field reflect.StructField) (namespace string, key string) {
	namespace = defaultNamespace
	tag, hasTag := field.Tag.Lookup("apollo")
	fieldName := field.Name
	if !hasTag {
		return namespace, fieldName
	}

	parts := strings.Split(tag, ":")
	switch len(parts) {
	case 1:
		key = parts[0]
	case 2:
		namespace, key = parts[0], parts[1]
	default:
	}
	return
}
