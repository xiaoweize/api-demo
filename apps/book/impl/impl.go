package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/apps/book"
	"github.com/xiaoweize/api-demo/conf"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

//grpc接口的实现类 要满足外层grpc生成的ServiceServer定义
type service struct {
	db *sql.DB

	log logger.Logger
	//必须嵌套 用于实现ServiceServer
	book.UnimplementedServiceServer
}

//初始化配置
func (s *service) Config() {
	s.db = conf.C().MySQL.GetDB()
	s.log = zap.L().Named(s.Name())
}

func (s *service) Name() string {
	return book.AppName
}

//注册
func (s *service) Registry(server *grpc.Server) {
	book.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpc(svr)
}
