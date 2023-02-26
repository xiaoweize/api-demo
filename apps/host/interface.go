package host

import (
	"context"
)

//定义业务支持的API接口操作，有如下API接口
//POST： /hosts/ 新增主机
//GET: /hosts/ 查询主机列表
//GET: /hosts/:id/ 查询主机详情
//PATCH: /hosts/:id 主机更新
//DELETE: /hosts/:id 删除主机
type Service interface {
	//新增主机，返回新增的主机实例对象
	CreateHost(context.Context, *Host) (*Host, error)
	//查询主机列表，返回给前端做主机列表页面
	//注意这里不要返回[]*Host,重新创建一个 HostSet 返回，保证了返回参数的一致性: 结构体指针
	QueryHost(context.Context, *QueryHostRequest) (*HostSet, error)
	//查询主机详情，返回给前端做主机详情页
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	//主机更新，返回更新后的主机信息给前端
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	//删除主机也需要返回Host对象，前端用于展示当前删除的Host信息
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
}
