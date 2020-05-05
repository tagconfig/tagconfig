package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shima-park/agollo"
	"github.com/tagconfig/tagconfig"
	"github.com/tagconfig/tagconfig/apollo"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalTagConfig(m map[string]string) (err error) {
	d.Duration, err = time.ParseDuration(m["."])
	return
}

// Config 配置结构体
// 请参考 http://106.54.227.205/config.html?#/appid=tagconfig apollo:admin
type Config struct {
	ServiceName string   `apollo:"serviceName"`
	Duration    Duration `apollo:"duration"`
	UserInfo    struct {
		Name      string `apollo:"name"`
		Age       int64  `apollo:"age"`
		Email     string `apollo:"email"`
		Phone     string
		Followers []int64 `apollo:"followers,json"`
	} `apollo:"userinfo"`
	Friend struct {
		Trade struct {
			Amount float64 `apollo:"amount"`
		} `apollo:"trade"`
	} `apollo:"development.common-mysql:friend"`
}

func main() {
	config := new(Config)
	var token = ""
	var configServerURL = ""
	var cluster = ""
	var appid = "tagconfig"

	var opts = []agollo.Option{
		agollo.AutoFetchOnCacheMiss(),
		agollo.Cluster(cluster),
		agollo.WithApolloClient(agollo.NewApolloClient(agollo.WithDoer(&TokenDoer{Token: token}))),
	}

	c, err := agollo.New(configServerURL, appid, opts...)
	if err != nil {
		panic(err)
	}

	client := apollo.NewClient(c, appid, []string{"application", "development.common-mysql"})
	fmt.Println(client.Paires())
	decoder := tagconfig.NewDecoder(client)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(config)
	fmt.Println(string(bs))
}
