package protocol

import (
	"net"
	"os"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/conf"
	"google.golang.org/grpc"
)

type GrpcService struct {
	svc *grpc.Server
	l   logger.Logger
}

func NewGrpcService() *GrpcService {
	return &GrpcService{
		svc: grpc.NewServer(),
		l:   zap.L().Named("GrpcService"),
	}
}

func (g *GrpcService) Start() {
	//初始化并注册grpc服务
	apps.InitGrpc(g.svc)
	//打印已加载的服务
	appNames := apps.LoadedGrpcApps()
	g.l.Infof("Load Grpc Apps:%s", appNames)

	//打开socket监听
	lis, err := net.Listen("tcp", conf.C().App.GrpcAddr())
	if err != nil {
		g.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}
	//创建grpc监听
	g.l.Infof("GRPC 服务监听地址: %s", conf.C().App.GrpcAddr())
	if err := g.svc.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			g.l.Info("service is stopped")
			os.Exit(1)
		}

		g.l.Errorf("start grpc service error, %s", err.Error())
		return
	}

}

func (g *GrpcService) Stop() {
	g.svc.GracefulStop()
}
