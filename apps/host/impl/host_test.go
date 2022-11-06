package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/stretchr/testify/assert"
	"github.com/xiaoweize/api-demo/apps/host"
	"github.com/xiaoweize/api-demo/apps/host/impl"
	"github.com/xiaoweize/api-demo/conf"
)

var (
	//定义一个满足该接口的实例
	service host.Service
)

func init() {
	//加载配置
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		panic(err)
	}
	//设置全局logger为debug级别输出到stderr
	zap.DevelopmentSetup()
	//host service的具体实现,初始化之前要加载全局配置，否则会引发panic
	service = impl.NewHostServiceImpl()
}

func TestCreateHost(t *testing.T) {
	should := assert.New(t)
	ins := host.NewHost()
	ins.Name = "test"
	ins.Id = "001"
	ins.Region = "cn-hangzhou"
	ins.Type = "sm1"
	ins.CPU = 2
	ins.Memory = 2048
	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}
