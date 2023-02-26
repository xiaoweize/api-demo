package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/apps/host"
	"github.com/xiaoweize/api-demo/conf"
)

//接口实现的静态检查->vscode检查
// var _ host.Service = (*HostServiceImpl)(nil)

var (
	impl = &HostServiceImpl{}
)

type HostServiceImpl struct {
	//业务模块相关功能都放在这里如日志、数据库配置等
	//Logger程序日志接口, 用于适配多种第三方日志插件
	l  logger.Logger
	db *sql.DB
}

//HostServiceImpl结构体实例构造函数
//注意调用该函数之前，保证全局对象config已经初始化，否则会发生panic
func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		//Host service 服务的子 Logger
		l: zap.L().Named("Host"),
		//未初始化全局config对象，会发生panic
		db: conf.C().MySQL.GetDB(),
	}
}

//之前是手动将HostService注册到Ioc容器
//apps.HostService = impl.NewHostServiceImpl()
//可以通过import匿名导入 自动执行注册逻辑
//import _ app  该app模块就注册到了Ioc层
func init() {
	//为什么没有apps.HostService = NewHostServiceImpl() 这样会造成config.C()并未准备好而发生panic
	//这样将对象的注册和对象的初始化这两个逻辑分开
	//这是app模块的注册逻辑 将impl这个实例(没有初始化)注册到Ioc容器中 同时实例化了Ioc容器中的HostService
	apps.RegistryImpl(impl)
}

//这是对象的初始化逻辑，只要保证全局config和全局logger初始化就能正常执行此初始化逻辑
func (i *HostServiceImpl) Config() {
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.GetDB()
}

//app模块/服务的名称 用于在Ioc层中的注册
func (i *HostServiceImpl) Name() string {
	return host.AppName
}
