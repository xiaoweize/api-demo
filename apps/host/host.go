package host

import (
	"context"
	"time"

	"github.com/go-playground/validator"
)

//定义资源厂家
type Vendor int

const (
	//自建机房 枚举
	PrivateIDC Vendor = iota
	//腾讯云
	Txyun
	//阿里云
	Aliyun
)

//用于创建Host资源实例
func NewHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

//定义数据模型,用于数据库存储，资源对象主要由Resource、Describe字段组成
type Host struct {
	//资源公共属性部分 如host与slb资源共有的字段Id、verndor、region等 主要用于资源索引表
	*Resource
	//资源独有属性部分 如主机资源有cpu、memory等独有属性 主要用于资源实例表 每种资源类型单独一个表
	*Describe
}

type Resource struct {
	Id          string            `json:"id"  validate:"required"`     // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                      // 厂商
	Region      string            `json:"region"  validate:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                   // 创建时间
	ExpireAt    int64             `json:"expire_at"`                   // 过期时间
	Type        string            `json:"type"  validate:"required"`   // 规格
	Name        string            `json:"name"  validate:"required"`   // 名称
	Description string            `json:"description"`                 // 描述
	Status      string            `json:"status"`                      // 服务商中的状态
	Tags        map[string]string `json:"tags"`                        // 标签
	UpdateAt    int64             `json:"update_at"`                   // 更新时间
	SyncAt      int64             `json:"sync_at"`                     // 同步时间
	Account     string            `json:"accout"`                      // 资源的所属账号
	PublicIP    string            `json:"public_ip"`                   // 公网IP
	PrivateIP   string            `json:"private_ip"`                  // 内网IP
}

type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

//定义业务支持的API接口操作，有如下API接口
//POST： /hosts/ 新增主机
//GET: /hosts/ 查询主机列表
//GET: /hosts/:id/ 查询主机详情
//PATCH: /hosts/:id 主机更新
//DELETE: /hosts/:id 删除主机
type Service interface {
	//新增主机，返回新增的主机实例对象
	CreateHost(context.Context, *Host) (*Host, error)
	//查询主机列表，返回主机列表，
	//注意这里没有返回[]*Host,而是重新创建了一个 主机集 的结构体HostSet来返回，里面带有主机数量属性，用于前端列表页的展示操作如换页等
	QueryHost(context.Context, *QueryHostRequest) (*HostSet, error)
	//查询主机详情，返回查询的主机信息给前端
	DescribeHost(context.Context, *QueryHostRequest) (*Host, error)
	//主机更新，返回更新后的主机信息给前端
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	//删除主机也需要返回Host对象，用于前端打印当前的删除Host的相关信息
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
}

var (
	validate = validator.New()
)

//主机字段校验方法
func (h *Host) Validate() error {
	return validate.Struct(h)
}

//给Host字段注入默认值
func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

//主机集的结构体，用于查询主机列表时返回的参数
type HostSet struct {
	Items []*Host
	Total int
}

type QueryHostRequest struct{}

type UpdateHostRequest struct {
	*Describe
}

type DeleteHostRequest struct{}
