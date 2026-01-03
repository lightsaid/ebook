package config

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DEV  = "dev"
	PROD = "prod"
)

// Load 加载配置到conf上，filenames 支持多个配置文件，
// 对应相同字段，以后者为准，解析过程中遇到错误回直接panic，因此可以忽略返回的错误
func Load(conf any, filenames ...string) error {
	// 收集解析得配置k:v
	obj := make(map[string]string)

	for _, file := range filenames {
		c, err := readFile(file)
		if err != nil {
			panic(fmt.Sprintf("读取解析配置[%s]错误: %q", file, err))
		}
		maps.Copy(obj, c)
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

// mapToStruct 将解析得到的配置文件key:value转换为struct
func mapToStruct(obj map[string]string, conf any) error {
	sType := reflect.TypeOf(conf).Elem()
	sValue := reflect.ValueOf(conf).Elem()

	for i := 0; i < sType.NumField(); i++ {
		ft := sType.Field(i)
		fv := sValue.Field(i)

		// --- 生成 env key ---
		mapKey := ft.Name
		if tagVal, ok := ft.Tag.Lookup("env"); ok && strings.TrimSpace(tagVal) != "" {
			mapKey = tagVal
		}

		// 获取 env 对应值
		valStr, ok := obj[mapKey]
		hasValue := ok

		// ======= 1. 嵌套结构体 ========
		if fv.Kind() == reflect.Struct {
			// 仅递归用户自定义 struct
			if fv.Type().PkgPath() != "" {
				if err := mapToStruct(obj, fv.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		// ======= 2. 指针类型 (struct pointer + base type pointer) ========
		if fv.Kind() == reflect.Pointer {
			// 分配指针
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}

			elem := fv.Elem()

			if elem.Kind() == reflect.Struct {
				// 递归填充结构体
				if err := mapToStruct(obj, elem.Addr().Interface()); err != nil {
					return err
				}
				continue
			}

			// 基本类型指针，则继续按基础类型解析
			fv = elem
		}

		// 如果当前字段没有 env 值，不赋值
		if !hasValue {
			continue
		}

		if handler, ok := customParsers[fv.Type().String()]; ok {
			raw, err := handler(valStr)
			if err != nil {
				return fmt.Errorf("parse %s failed: %w", mapKey, err)
			}
			fv.Set(reflect.ValueOf(raw).Convert(fv.Type()))
			continue
		}

		// ======= 3. 内置类型解析 ========
		if handler, ok := defaultBuiltInParsers[fv.Kind()]; ok {
			raw, err := handler(valStr)
			if err != nil {
				return fmt.Errorf("parse %s failed: %w", mapKey, err)
			}
			fv.Set(reflect.ValueOf(raw).Convert(fv.Type()))
			continue
		}

		// ======= 4. 非内置类型 (alias type) ========
		underlyingKind := fv.Type().Kind()
		if handler, ok := defaultBuiltInParsers[underlyingKind]; ok {
			raw, err := handler(valStr)
			if err != nil {
				return fmt.Errorf("parse custom type %s failed: %w", mapKey, err)
			}

			// 转换为目标的自定义类型
			fv.Set(reflect.ValueOf(raw).Convert(fv.Type()))
			continue
		}

		return fmt.Errorf("unsupported field type: %s (%s)", ft.Name, fv.Type())
	}

	return nil
}
