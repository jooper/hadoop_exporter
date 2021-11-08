package Utiles

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//配置文件中字母要小写，结构体属性首字母要大写

type Myconf struct {
	Ip                 string
	Port               string
	FaasApiUrl         string
	StartSendTime      string
	SendMaxCountPerDay int
	Devices            []Device
	WarnFrequency      int
	SendFrequency      int
}
type Device struct {
	DevId string
	Nodes []Node
}
type Node struct {
	PkId     string
	BkId     string
	Index    string
	MinValue float32
	MaxValue float32
	DataType string
}

func Yml() Myconf {
	data, _ := ioutil.ReadFile("config.yml")
	t := Myconf{}
	yaml.Unmarshal(data, &t)
	return t
}

func ReadYml() {
	data, _ := ioutil.ReadFile("config.yml")
	//fmt.Println(string(data))
	t := Myconf{}
	//把yaml形式的字符串解析成struct类型
	yaml.Unmarshal(data, &t)

	fmt.Println(t.Ip)

	//fmt.Println("初始数据", t)
	//if t.Ipport=="" {
	//	fmt.Println("配置文件设置错误")
	//	return
	//}

	//把struct形式的字符串解析成yaml类型
	//d, _ := yaml.Marshal(&t)
	//fmt.Println("看看 :", string(d))
}
