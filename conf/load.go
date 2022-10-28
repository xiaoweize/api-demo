package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

//将配置文件映射成config对象

//从toml格式文件中加载配置
func LoadConfigFromToml(filepath string) error {
	//初始化全局config实例
	config = NewDefaultConfig()
	_, err := toml.DecodeFile(filepath, config)
	//如果配置文件加载失败就使用默认配置
	if err != nil {
		return fmt.Errorf("load config from filepath->%s, error->%s", filepath, err)
	}
	return nil
}

//从环境变量加载配置
func LoadConfigFromEnv() error {
	//初始化全局config实例
	config = NewDefaultConfig()
	err := env.Parse(config)
	//如果环境变量加载失败就使用默认配置
	if err != nil {
		return err
	}
	return nil

}
