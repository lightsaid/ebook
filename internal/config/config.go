package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lightsaid/ebook/pkg/logger"
)

// 环境参数，提供其他包使用
const (
	Env_Dev  = "dev"
	Env_Prod = "prod"
)

// AppConfig 应用程序配置
type AppConfig struct {
	Env          string `json:"env"`          // 环境参数：dev | prod
	Port         int    `json:"port"`         // 服务端口
	DSN          string `json:"dsn"`          // 数据库链接
	MaxOpenConns int    `json:"maxOpenConns"` // 数据库最大链接数
	MaxIdleConns int    `json:"maxIdleConns"` // 数据库链接最大空闲数
	MaxIdleTime  string `json:"maxIdleTime"`  // 数据库链接最大空闲时间
}

// MaxIdleTimeToDuration 将 MaxIdleTime 转换成 time.Duration 返回，如果转换出错，返回默认值
func (app *AppConfig) MaxIdleTimeToDuration() time.Duration {
	dur, err := time.ParseDuration(app.MaxIdleTime)
	if err != nil {
		logger.ErrorfoLog.Println("time.ParseDuration(app.MaxIdleTime) failed: " + err.Error())
		return time.Minute * 5
	}
	return dur
}

// LoadAppConfig 根据配置文件路径加载配置
func LoadAppConfig(path string) (cfg AppConfig, err error) {
	var buf []byte
	buf, err = os.ReadFile(path)
	if err != nil {
		return
	}
	if err = json.Unmarshal(buf, &cfg); err != nil {
		return
	}

	return
}

func (a *AppConfig) Println() {
	if a.Env == Env_Dev {
		buf, _ := json.MarshalIndent(a, "", " ")
		fmt.Println(string(buf))
	}
}
