package apollo

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/n0trace/tagconfig"
	"github.com/shima-park/agollo"
)

var (
	defaultNamespace = "application"
)

type Client struct {
	agollo.Agollo
	namespaces []string
}

func MustClient(appid string, namespaces []string) *Client {
	c, err := GetClient(appid, namespaces)
	if err != nil {
		panic(err)
	}
	return c
}

func GetClient(appid string, namespaces []string) (client *Client, err error) {
	client = new(Client)
	client.namespaces = namespaces
	var (
		token           = os.Getenv("CONFIG_CENTER_TOKEN")
		configServerURL = os.Getenv("CONFIG_CENTER_URL")
		cluster         = os.Getenv("RUNTIME_CLUSTER")
	)

	if token == "" {
		err = errors.New("token empty")
	}

	if configServerURL == "" {
		err = errors.New("config server empty")
	}

	if cluster == "" {
		err = errors.New("cluster empty")
	}

	if err != nil {
		return
	}

	var opts = []agollo.Option{
		agollo.AutoFetchOnCacheMiss(),
		agollo.Cluster(cluster),
		agollo.WithApolloClient(agollo.NewApolloClient()),
	}

	client.Agollo, err = agollo.New(configServerURL, appid, opts...)
	return
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
