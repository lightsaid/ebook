package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Mode string
	Port int
	Foo  bool `env:"FOO"`
	Cat  Cat
	Dog  *Dog // TODO: 指针类型
}

type Cat struct {
	CatAge float32
}

type Dog struct {
	DogName string
}

func TestLoadConfig(t *testing.T) {
	conf := new(TestConfig)
	// conf.Dog = new(Dog)
	err := Load(conf, "./test.env")
	require.NoError(t, err)
	by, _ := json.MarshalIndent(conf, "", "\t")
	fmt.Println(string(by))
}
