package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/apps/host"
)

//通过实体类将内部的函数和方法通过http协议暴露出去
//需要依赖内部接口的实现
//该struct要实现gin的Handler函数签名 func(*Context)
type Handler struct {
	svc host.Service
}

var handler = &Handler{}

//host api初始化 面向接口
func NewHostHTTPHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Config() {
	// if apps.HostService == nil {
	// 	panic("Dependence Host Service Required!")
	// }
	//从Ioc容器层里面获取HostService的实例对象
	// h.svc = apps.HostService

	//使用断言方式从Ioc层获取host.Service实例对象
	h.svc = apps.GetImpl(host.AppName).(host.Service)
}

//http hanlder的注册
func (h *Handler) Registry(r gin.IRouter) {
	//新增主机
	r.POST("/hosts", h.createHost)
	//查询主机列表
	r.GET("/hosts", h.queryHost)
	//查询主机详情,如果不接:id就是查询主机列表
	r.GET("hosts/:id", h.describeHost)
	//全量更新
	r.PUT("hosts/:id", h.putHost)
	//部分更新
	r.PATCH("hosts/:id", h.patchHost)
	//删除主机
	r.DELETE("hosts/:id", h.deleteHost)
}

//http名称
func (h *Handler) Name() string {
	return host.AppName
}

//通过init方法 完成http handler的注册

func init() {
	apps.RegistryGin(handler)
}
