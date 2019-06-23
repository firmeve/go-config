package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	directory = "./testdata/"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		t.Fatalf("Config create error. Detail :%s", err.Error())
	}

	fmt.Printf("%T", config)
}

func TestConfig_Set(t *testing.T) {
	tests := []struct {
		file  string
		key   string
		value string
	}{
		{"app", "x", "x"},
		{"app", "s1.x", "s1x"},
		{"app", "s1.z.y", "s1xy"},
		{"new", "x", "x"},
	}

	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	for _, test := range tests {
		//fmt.Println(test.file,test.key,test.value,i)
		err = config.Set(test.file+"."+test.key, test.value)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			t.Fail()
		}
	}
}

// 正常数据测试
func TestConfig_Get(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	//// 最后一个可能还不是key
	//z, err := config.Get("app.ggg.z")
	////fmt.Println(z.(*ini.Section))
	//if _, ok := z.(*ini.Section); ok {
	//	fmt.Println("SSSSSSSSSSSSSSS")
	//}
	//fmt.Printf("%#v", z.(*ini.Section))

	tests := []struct {
		file  string
		key   string
		value string
	}{
		{"app", "t_key", "t_value"},
		{"app", "t1.t2", "t2_value"},
		{"new", "nt.nt2.nt3", "nt3_value"},
	}

	for _, test := range tests {
		err := config.Set(test.file+"."+test.key, test.value)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			t.Fail()
		}
	}

	// 当一个值的时候，返回ini.File完整对象
	res, err := config.Get(`app`)
	if _, ok := res.(*ini.File); !ok {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	// 当是.的时候，返回默认section的key值
	res1, err := config.Get(`app.t_key`)
	fmt.Printf("\n%s\n",res1.(*ini.Key).Value())
	assert.Equal(t, `t_value`, res1.(*ini.Key).Value())

	// 当是3个值的时候，返回指定section的key值
	res2, err := config.Get(`app.t1.t2`)
	fmt.Printf("\n%s\n",res2.(*ini.Key).Value())
	assert.Equal(t, `t2_value`, res2.(*ini.Key).Value())

	// 当是4个以上值的时候，返回指定section子section的key值
	res3, err := config.Get(`new.nt.nt2.nt3`)
	fmt.Printf("\n%s\n",res3.(*ini.Key).Value())
	assert.Equal(t, `nt3_value`, res3.(*ini.Key).Value())

	// 优先查找section的拼接
	res4, err := config.Get(`new.nt.nt2`)
	fmt.Printf("\n%T\n",res4)
	if _,ok := res4.(*ini.Section); !ok {
		t.Fail()
	}
}

// default 非正常数据测试
func TestConfig_GetDefault(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	// 当一个值的时候，返回ini.File完整对象
	res := config.GetDefault(`ssss`,`def`)

	assert.Equal(t,`def`,res.(string))
}

func TestConfig_All(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	for _,v := range config.All() {
		fmt.Printf("%v",v)
	}
}