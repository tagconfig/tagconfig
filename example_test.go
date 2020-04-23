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
	var paires = map[string]string{
		"foo":               "foo",
		"duration":          "3s",
		"userinfo.name":     "n0trace",
		"userinfo.Age":      "20",
		"userinfo.bio":      "don’t be evil",
		"userinfo.male":     "true",
		"userinfo.Follower": "[1,2,101]",
		"userinfo.Followed": `{"3":"hello"}`,
	}

	type Message struct {
		UserInfo struct {
			Name     string `test:"name"`
			Age      uint64
			Male     bool        `test:"male"`
			Bio      interface{} `test:"bio"`
			Follower []int64
			Followed map[int64]string
		} `test:"userinfo"`
		Foo      string `test:"foo"`
		Duration `test:"duration"`
	}
	provider := &TestConfigProvider{configMap: paires}
	decoder := NewDecoder(provider)
	message := new(Message)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message.Foo)
	fmt.Println(message.Duration.Seconds())
	fmt.Println(message.UserInfo.Name)
	fmt.Println(message.UserInfo.Age)
	fmt.Println(message.UserInfo.Male)
	fmt.Println(message.UserInfo.Bio)
	fmt.Println(message.UserInfo.Follower)
	fmt.Println(message.UserInfo.Followed)
	// Output:
	// foo
	// 3
	// n0trace
	// 20
	// true
	// don’t be evil
	// [1 2 101]
	// map[3:hello]
}
