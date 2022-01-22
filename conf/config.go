package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"sync"
)

type tomlConfig struct {
	Name       string `toml:"app_name"`
	Port       int
	Mode       string
	StaticPath string
	DB         tomlDatabase `toml:"database"`
	JWT        tomlJWT      `toml:"jwt"`
}

type tomlDatabase struct {
	Host     string
	Name     string
	User     string
	Password string
}

type tomlJWT struct {
	SecretKey string
}

const (
	PROD             Env = "production"
	DEV              Env = "development"
	defaultSecretKey     = "chang'emoonadsfwerf"
)

type Env string

func (e Env) String() string {
	return string(e)
}

var (
	conf        *tomlConfig
	once        sync.Once
	ProjectName = "ChangE"
	Version     = "0.0.1"
)

func GetSecretKey() string {
	if Config().JWT.SecretKey == "" {
		return defaultSecretKey
	}
	return Config().JWT.SecretKey
}

func (c *tomlConfig) Print() {
	fmt.Println("===================配置文件================")
	fmt.Printf("NAME=%s\n", c.Name)
	fmt.Printf("PORT=%d\n", c.Port)
	fmt.Printf("MODE=%s\n", c.Mode)
	fmt.Printf("STATIC_PATH=%s\n", c.StaticPath)

	fmt.Println("[DATABASE]")
	fmt.Printf("HOST=%s\n", c.DB.Host)
	fmt.Printf("NAME=%s\n", c.DB.Name)
	fmt.Printf("USER=%s\n", c.DB.User)

	fmt.Println("[JWT]")
	fmt.Printf("SECRET_KEY=%s\n", c.JWT.SecretKey)

	fmt.Println("==========================================")
}

// Config 单例模式
func Config() *tomlConfig {
	// exec only once
	once.Do(func() {
		if _, err := toml.DecodeFile("./conf/nine.conf", &conf); err != nil {
			panic(err)
		}

		conf.Print()
	})

	return conf
}
