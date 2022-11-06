package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/xiaoweize/api-demo/apps/host"
	"github.com/xiaoweize/api-demo/conf"
)

//接口实现的静态检查->vscode检查
var _ host.Service = (*HostServiceImpl)(nil)

type HostServiceImpl struct {
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
