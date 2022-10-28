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
func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		//Host service 服务的日志Logger 用zap实现
		l: zap.L().Named("Host"),
		//获取sql.DB,注意在之前要加载全局配置
		db: conf.C().MySQL.GetDB(),
	}
}
