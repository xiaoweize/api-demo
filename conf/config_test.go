package conf_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiaoweize/api-demo/conf"
)

func TestLoadConfigFromToml(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../etc/demo.toml")
	//即使配置文件加载失败也能得到默认配置
	fmt.Println(conf.C().App.Host)
	//此时conf.C().MySQL.GetDB()可以正常初始化全局db实例
	fmt.Println(conf.C().MySQL.GetDB())

	if should.NoError(err) {
		should.Equal("api_demo", conf.C().MySQL.Database)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	should := assert.New(t)
	//设置环境变量
	// os.Setenv("MYSQL_DATABASE", "unit_test")
	err := conf.LoadConfigFromEnv()
	//即使环境变量不存在也能得到默认配置
	fmt.Println(conf.C().MySQL.Database)
	if should.NoError(err) {
		should.Equal("unit_test", conf.C().MySQL.Database)
	}
}

func TestGetDB(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../etc/demo.toml")
	if should.NoError(err) {
		conf.C().MySQL.GetDB()
	}
}
