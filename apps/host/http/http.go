package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoweize/api-demo/apps/host"
)

//通过实体类将内部的函数和方法通过http协议暴露出去
//需要依赖内部接口的实现
//该struct要实现gin的Handler函数签名 func(*Context)
type Handler struct {
	svc host.Service
}

//host api初始化
//生成handler对象，在main中将具体的service实例传参进去
func NewHostHTTPHandler(svc host.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

//http hanlder的注册
//获取handler实例后开始注册路由
//注册路由方法，要给每个产品的每一个操作都加上路由，这里开始就显得不太优雅了，更好的方式使用IOC
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
}
