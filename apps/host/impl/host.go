package impl

import (
	"context"

	"github.com/xiaoweize/api-demo/apps/host"
)

//业务处理层(Controller层 )
//结构体HostServiceImpl方法 实现了host.Service接口
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	//打印info级别日志
	// i.l.Info("create host")
	// //Named创建子logger 分别创建create/test子logger对应的debug/info级别日志
	// i.l.Named("Create").Debug("create host")
	// i.l.Named("test").Info("create host")
	// //格式化打印debug日志
	// i.l.Debugf("create host %s", ins.Name)
	// //携带额外的meta数据，常用于Trace系统
	// i.l.With(logger.NewAny("request-id", "req01")).Debug("create host with meta kv")

	//校验Host结构体字段合法性使用github.com/go-playground/validator
	var err error
	if err = ins.Validate(); err != nil {
		return nil, err
	}
	//默认值填充
	ins.InjectDefault()
	//由dao模块将对象转换为数据库数据 
	return ins, i.save(ctx, ins)
}

func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error) {
	return nil, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.QueryHostRequest) (*host.Host, error) {
	return nil, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	return nil, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	return nil, nil
}
