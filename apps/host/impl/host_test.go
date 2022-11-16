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

func TestDescribeHost(t *testing.T) {
	should := assert.New(t)
	req := host.NewDescribeHostRequestWithId("kq")
	ins, err := service.DescribeHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func TestQueryHost(t *testing.T) {
	should := assert.New(t)
	req := host.NewQueryHostRequest()
	req.Keywords = "kq"
	set, err := service.QueryHost(context.Background(), req)
	if should.NoError(err) {
		for i := range set.Items {
			fmt.Println(set.Items[i].Id)
		}
	}

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

func TestPatchUpdateHost(t *testing.T) {
	should := assert.New(t)
	//测试模拟requst请求数据 生成指定id字段的host
	req := host.NewPatchUpdateHostRequestWithId("223")
	req.Name = ""
	//调用后端接口处理请求数据
	//patch方法只覆盖请求的值字段，所以能够成功覆盖
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Name)
	}
}

func TestPutUpdateHost(t *testing.T) {
	should := assert.New(t)
	//测试模拟requst请求数据
	req := host.NewPutUpdateHostRequestWithId("223")
	req.Region = "hangzhou"
	req.Type = "large"
	req.Name = "测试主机"
	//调用后端接口处理请求数据
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Region)
	}
}
