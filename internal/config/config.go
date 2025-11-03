package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DEV  = "dev"
	PROD = "prod"
)

type AppConfig struct {
	Db
	Jwt
}

type Db struct{}

type Jwt struct{}

func Load(conf any, filenames ...string) error {
	// 收集解析得配置k:v
	obj := make(map[string]string)

	for _, file := range filenames {
		c, err := readFile(file)
		if err != nil {
			panic(fmt.Sprintf("读取解析配置[%s]错误: %q", file, err))
		}
		for k, v := range c {
			obj[k] = v
		}
	}

	err := mapToStruct(obj, conf)
	if err != nil {
		panic(fmt.Sprintf("配置转结构体错误: %s", err))
	}

	return nil
}

// readFile 读取配置，使用 godotenv.Parse 解析配置
func readFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return godotenv.Parse(file)
}

func mapToStruct(obj map[string]string, conf any) error {
	var tagName = "env" // 配置字段指定映射的tag

	sType := reflect.TypeOf(conf).Elem()
	sValue := reflect.ValueOf(conf).Elem()

	for i := 0; i < sType.NumField(); i++ {
		ft := sType.Field(i)
		fv := sValue.Field(i)

		// 判断是否是嵌套结构体，进行递归处理
		// 嵌套结构体或结构体指针
		if fv.Kind() == reflect.Struct && fv.CanSet() {
			if err := mapToStruct(obj, fv.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if (fv.Kind() == reflect.Ptr) && fv.CanSet() {
			if err := mapToStruct(obj, fv.Interface()); err != nil {
				return err
			}
			continue
		}

		// 查找obj对应的key
		mapKey := ft.Name
		if key, ok := ft.Tag.Lookup(tagName); ok && strings.Trim(key, "") != "" {
			mapKey = key
		}

		// 获取value
		mapVal, ok := obj[mapKey]
		if !ok {
			continue
		}

		// 赋值
		handle, ok := defaultBuiltInParsers[fv.Kind()]
		if ok {
			val, err := handle(mapVal)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(val).Convert(ft.Type))
		}
	}

	// TODO: 非内置类型

	return nil
}
