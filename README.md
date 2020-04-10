# tagconfig

```sh
cd example
go run .
```

## 用法


```go
import (
	"github.com/shima-park/agollo"
	"github.com/n0trace/tagconfig"
	"github.com/n0trace/tagconfig/apollo"
)
//配置结构体
config := new(Config)
//new一个agollo客户端(第三方)
c, err := agollo.New(configServerURL, appid, opts...)
if err != nil {
	panic(err)
}
//用第三方客户端初始化一个配置provider
client := apollo.NewClient(c, appid, []string{"application", "development.common-mysql"})
//获得一个apollo配置解析器
decoder := tagconfig.NewDecoder(client)
//把配置解析到config
err := decoder.Decode(&config)
```

## 功能

从apollo拉取properties配置类型的config,并把这些配置scan到一个go的结构体,会把key按照"."分割并解析为嵌套结构，具体的使用如下:

```golang
type Config struct {
	Foo      string `apollo:"foo"` //读取application命名空间下的key为foo的配置
	Bar      string //读取application命名空间下的key为Bar的配置
	Config3D struct {
		Config2D struct {
			Foo       string      //读取n1命名空间下的key为config3d.config2d.foo的配置
			Bar       int64       //读取n1命名空间下的key为config3d.config2d.bar的配置
			Interface interface{} `apollo:"interface"`
		} `apollo:"config2d"`
	} `apollo:"n1:config3d"`
}
```

> 注意事项

1.当子结构为slice或是map,会当作json解析

## 自定义解析器

如果需要特殊解析，实现下面的方法即可，

```go
type Config struct{
    Foo struct{
        Bar string
    }
}
func (foo *Foo)UnmarshalTagConfig(m map[string]string) (err error) {
    //m是一个map
    //m["."]可以匹配到apollo设置为Config.Foo的kv
    //m["Bar"]可以匹配到apollo设置为Config.Foo.Bar的kv
}
```
也可以参考[example](/example)