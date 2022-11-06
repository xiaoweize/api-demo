package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

//如何把配置映射成config对象

//从toml格式文件中加载配置
func LoadConfigFromToml(filepath string) error {
	//初始化全局config实例，此时conf.C().MySQL.GetDB()可以正常初始化全局db实例
	config = NewDefaultConfig()
	//从文件中解析配置
	_, err := toml.DecodeFile(filepath, config)
	//解析失败返回err错误，config使用默认值
	if err != nil {
		return fmt.Errorf("load config from filepath->%s, error->%s", filepath, err)
	}
	return nil
}

//从环境变量加载配置
func LoadConfigFromEnv() error {
	//初始化全局config实例,此时conf.C().MySQL.GetDB()可以正常初始化全局db实例
	config = NewDefaultConfig()
	//从环境变量加载配置，加载失败使用默认配置
	return env.Parse(config)
}
