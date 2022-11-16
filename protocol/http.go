package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/conf"
)

func NewHttpService() *HttpService {
	//new gin router 并没有加载handler
	r := gin.Default()

	//配置httpServer 使用标准库中的http.Server
	server := &http.Server{
		//下面5个参数都建议配置上
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HttpAddr(),
		//把所有请求都交给gin处理
		Handler: r,
	}
	return &HttpService{
		server: server,
		//分配全局子Logger
		l: zap.L().Named("HttpService"),
		r: r,
	}

}

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

//HTTP服务的启动
func (s *HttpService) Start() error {
	//加载Handler 把所有模块的Handler注册给gin router
	apps.InitGin(s.r)
	//打印已加载好gin的apps信息
	appNames := apps.LoadedGinApps()
	s.l.Infof("loaded gin apps%v", appNames)

	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			//排除正常关闭的情况
			s.l.Info("Http Server Success Stopped!")
		} else {
			return fmt.Errorf("Start Http Server Error:%v", err)
		}
	}
	return nil
}

//HTTP服务的停止 gracefully优雅的关闭Http Server
func (h *HttpService) Stop() {
	h.l.Info("Start gracefully shuts down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() //退出时发送上content cancel

	//关闭超过30秒也会强制退出
	if err := h.server.Shutdown(ctx); err != nil {
		h.l.Warnf("ShoutDown Http Server err:%s ", err)
	}
}
