package tagconfig

import (
	"fmt"
	"log"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalTagConfig(m map[string]string) (err error) {
	d.Duration, err = time.ParseDuration(m["."])
	return
}

func Example_decode() {
	var paires = [][]string{
		{"foo", "foo", "common"},
		{"duration", "3s"},
		{"userinfo.name", "n0trace"},
		{"userinfo.Age", "20"},
		{"userinfo.bio", "don’t be evil"},
		{"userinfo.male", "true"},
		{"userinfo.Follower", "[1,2,101]"},
		{"userinfo.Followed", `{"3":"hello"}`},
	}

	type Message struct {
		UserInfo struct {
			Name     string `tagconfig:"name"`
			Age      uint64
			Male     bool        `tagconfig:"male"`
			Bio      interface{} `tagconfig:"bio"`
			Follower []int64
			Followed map[int64]string
		} `tagconfig:"userinfo"`
		Foo      string `tagconfig:"common:foo"`
		Duration `tagconfig:"duration"`
	}
	provider := &TestConfigProvider{configs: paires}
	decoder := NewDecoder(provider)
	message := new(Message)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message.Foo)
	fmt.Println(message.UserInfo.Name)
	fmt.Println(message.UserInfo.Age)
	fmt.Println(message.Duration.Seconds())
	fmt.Println(message.UserInfo.Male)
	fmt.Println(message.UserInfo.Bio)
	fmt.Println(message.UserInfo.Follower)
	fmt.Println(message.UserInfo.Followed)
	// Output:
	// foo
	// n0trace
	// 20
	// 3
	// true
	// don’t be evil
	// [1 2 101]
	// map[3:hello]
}
