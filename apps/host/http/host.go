package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
	"github.com/xiaoweize/api-demo/apps/host"
)

//用于暴露host service方法 实现了HandlerFunc方法 CreateHost
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
	//处理成功后将ins返回给http api的调用方
	response.Success(c.Writer, ins)
}

func (h *Handler) queryHost(c *gin.Context) {
	//使用gin的querystring解析
	req := host.NewQueryHostRequest()
	pss := c.Query("page_size")
	pns := c.Query("page_number")
	if pss != "" {
		req.PageSize, _ = strconv.Atoi(pss)
	}
	if pns != "" {
		req.PageNumber, _ = strconv.Atoi(pns)
	}
	req.Keywords = c.Query("kws")

	set, err := h.svc.QueryHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
	} else {
		response.Success(c.Writer, set)
	}
}

func (h *Handler) describeHost(c *gin.Context) {
	//注意这里不需要判断路径参数，如果路径参数不存在就是搜索主机列表
	// id := c.Param("id")
	// if id == "" {
	// 	response.Failed(c.Writer, fmt.Errorf("Id cannot be empty!"))
	// 	return
	// }
	req := host.NewDescribeHostRequestWithId(c.Param("id"))
	ins, err := h.svc.DescribeHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
	} else {
		response.Success(c.Writer, ins)
	}
}

func (h *Handler) deleteHost(c *gin.Context) {
	req := host.NewDeleteHostRequestWithId(c.Param("id"))
	ins, err := h.svc.DeleteHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
	} else {
		response.Success(c.Writer, ins)
	}
}

func (h *Handler) putHost(c *gin.Context) {
	//querystring解析
	req := host.NewPutUpdateHostRequestWithId(c.Param("id"))
	//body json解析
	if err := c.Bind(&req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	//body中的id可能会将路径参数重的id重写掉，这里重新赋下值
	req.Id = c.Param("id")

	//将接收到的请求交给UpdateHost处理，UpdateHost处理逻辑如下
	// 1.从数据库中获取要更新的对象到内存中
	// 2.调用ins.put方法全量替换从1中获取的对象 注意请求中没有携带的字段会以空值覆盖已存在的对象
	// 3.进行validate验证，确保设置了required属性的字段必须带值
	// 4.将通过验证的req对象更新到数据库中，注意更新的update语句设置了仅允许更新的字段，所以要前端配合仅允许更新的字段
	//总结 put更新必须带上required字段值，然后携带update语句中的字段就进行更新，没有携带的就不更新
	if ins, err := h.svc.UpdateHost(c.Request.Context(), req); err != nil {
		response.Failed(c.Writer, err)
	} else {
		response.Success(c.Writer, ins)
	}

}

func (h *Handler) patchHost(c *gin.Context) {
	//querystring解析 将路径参数中的id对象值传给req
	req := host.NewPatchUpdateHostRequestWithId(c.Param("id"))
	//body json解析 请求中有字段的就填充对应的值 无字段的就是类型的零值
	if err := c.Bind(&req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	//body中的id可能会将路径参数重的id重写掉，这里重新赋下值
	req.Id = c.Param("id")

	//将接收到的请求交给UpdateHost处理 UpdateHost处理的逻辑如下
	//  1.从数据库中获取要更新的对象到内存中
	//	2.ins.Patch方法，将请求字段中的非空字段覆盖到已有对象对应的字段上(非空指的是没带字段或带了零值的字段)
	//	3.进行validate验证，其实patch这里不需要验证，因为更新的对象字段肯定有值，主要是给put方法用的
	//  4.更新数据库.注意数据库里面设置了仅允许指定字段更新，如果添加了额外的字段，会返回假象的修改 这里存在bug逻辑需要修复(应该让前端控制住修改范围)
	// 总结:patch可以修改局部字段 但携带的字段必须为非空值
	ins, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
	} else {
		response.Success(c.Writer, ins)
	}
}
