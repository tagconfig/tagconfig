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
		{"userinfo.Male", "true"},
		{"userinfo.follower", "[1,2,101]"},
		{"userinfo.followed.1.name", "hello"},
		{"userinfo.followed.2.Male", "true"},
		{"userinfo.followed.2.Age", "22"},
		{"userinfo.history.x", "beijing"},
		{"userinfo.history.y", "tianjin"},
	}
	type User struct {
		Name     string `tagconfig:"name"`
		Age      uint64
		Male     bool
		Bio      interface{}       `tagconfig:"bio"`
		Follower []int64           `tagconfig:"follower,json"`
		Followed map[int64]User    `tagconfig:"followed"`
		History  map[string]string `tagconfig:"history"`
	}

	type Message struct {
		UserInfo User   `tagconfig:"userinfo"`
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
	fmt.Println(message.UserInfo.Followed[1].Name)
	fmt.Println(message.UserInfo.Followed[2].Male)
	fmt.Println(message.UserInfo.Followed[2].Age)
	fmt.Println(message.UserInfo.History["x"])
	fmt.Println(message.UserInfo.History["y"])
	// Output:
	// foo
	// n0trace
	// 20
	// 3
	// true
	// don’t be evil
	// [1 2 101]
	// hello
	// true
	// 22
	// beijing
	// tianjin
}
