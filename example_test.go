package tagconfig

import (
	"fmt"
	"log"
)

func Example_decode() {
	var paires = map[string]string{
		"userinfo.name": "n0trace",
		"userinfo.Age":  "20",
		"userinfo.bio":  "don’t be evil",
		"foo":           "foo",
	}
	type Message struct {
		UserInfo struct {
			Name string `test:"name"`
			Age  int64
			Bio  interface{} `test:"bio"`
		} `test:"userinfo"`
		Foo string `test:"foo"`
	}
	provider := &TestConfigProvider{configMap: paires}
	decoder := NewDecoder(provider)
	message := new(Message)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message.Foo)
	fmt.Println(message.UserInfo.Name)
	fmt.Println(message.UserInfo.Age)
	fmt.Println(message.UserInfo.Bio)
	// Output:
	// foo
	// n0trace
	// 20
	// don’t be evil
}
