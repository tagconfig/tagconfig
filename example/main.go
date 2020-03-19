package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/n0trace/tagconfig"
	"github.com/n0trace/tagconfig/apollo"
)

type LibraryConfig struct {
	IPBlackList  []string           `apollo:"ip_black_list"`
	Oauth2Config map[int64][]string `apollo:"oauth2"` //{"<appid>":[<appid_key>,<appid_secret>]}
}

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

type Config struct {
	Product      string        `apollo:"product"`
	LibaryConfig LibraryConfig `apollo:"library"`
	Config3D     struct {
		Config2D struct {
			Foo       string      `apollo:"foo"`
			Bar       int64       `apollo:"bar"`
			Interface interface{} `apollo:"interface"`
		} `apollo:"config2d"` //2维配置
	} `apollo:"config3d"` //3维配置
	Duration Duration `apollo:"duration"`
}

func main() {
	config := new(Config)
	client := apollo.MustClient("example", []string{"application", "example-common"})
	decoder := tagconfig.NewDecoder(client)
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	bs, _ := json.Marshal(config)
	fmt.Println(string(bs))
}
