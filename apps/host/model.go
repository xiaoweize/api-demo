package host

import (
	"time"
)

//定义资源厂家
type Vendor int

const (
	//默认自建机房 枚举的默认值0
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

//定义数据模型,用于数据库存储，资源对象由Resource、Describe字段组成
type Host struct {
	//可以抽象出来的通用属性部分 如host与slb资源共有的字段Id、verndor、region等 存到一张表中，用于资源的搜索(搜索到的资源只会显示公共属性)
	*Resource
	//资源独有属性部分 如主机资源有cpu、memory等独有属性 主要用于资源实例表 每种资源类型单独一个表
	*Describe
}

type Resource struct {
	Id          string            `json:"id"  binding:"required"`     // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                     // 厂商
	Region      string            `json:"region"  binding:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                  // 创建时间
	ExpireAt    int64             `json:"expire_at"`                  // 过期时间
	Type        string            `json:"type"  binding:"required"`   // 规格
	Name        string            `json:"name"  binding:"required"`   // 名称
	Description string            `json:"description"`                // 描述
	Status      string            `json:"status"`                     // 服务商中的状态
	Tags        map[string]string `json:"tags"`                       // 标签
	UpdateAt    int64             `json:"update_at"`                  // 更新时间
	SyncAt      int64             `json:"sync_at"`                    // 同步时间
	Account     string            `json:"accout"`                     // 资源的所属账号
	PublicIP    string            `json:"public_ip"`                  // 公网IP
	PrivateIP   string            `json:"private_ip"`                 // 内网IP
}

type Describe struct {
	CPU          int    `json:"cpu" binding:"required"`    // 核数
	Memory       int    `json:"memory" binding:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                // GPU数量
	GPUSpec      string `json:"gpu_spec"`                  // GPU类型
	OSType       string `json:"os_type"`                   // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                   // 操作系统名称
	SerialNumber string `json:"serial_number"`             // 序列号
}

//给CreateAt字段注入默认值
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
	//仅允许修改资源独有属性
	*Describe
}

type DeleteHostRequest struct {
	//通过Id删除资源
	Id string
}
