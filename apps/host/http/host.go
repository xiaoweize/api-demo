package http

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
	"github.com/xiaoweize/api-demo/apps/host"
)

//用于暴露host service方法CreateHost
func (h *Handler) createHost(c *gin.Context) {
	//创建Host默认值实例 用于解析request数据
	ins := host.NewHost()
	//解析用户传递进来的参数到ins实例上
	if err := c.Bind(ins); err != nil {
		//绑定失败返回400响应码，以及json类型响应，这里用的是封装后的http库
		response.Failed(c.Writer, err)
		return
	}
	//成功获取到客户端发送的ins数据后 开始调用业务接口  从这里开始传递context
	//当request取消后，context会传递到mysql将事务取消掉
	ins, err := h.svc.CreateHost(c.Request.Context(), ins)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	//处理成功后将ins返回给前端
	response.Success(c.Writer, ins)
}
