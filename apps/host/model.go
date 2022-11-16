package host

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/imdario/mergo"
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
	//可以抽象出来的通用属性部分 如host与slb资源共有的字段Id、verndor、region等 存到一张表中，用于资源的全局搜索(搜索到的资源只会显示公共属性)
	//如果需要解锁的资源没有通用属性部分，也可以将其独立属性放这里用于搜索
	*Resource
	//资源独有属性部分 如主机资源有cpu、memory等独有属性 主要用于资源实例表 每种资源类型单独一个表
	*Describe
}

//binding
type Resource struct {
	Id          string            `json:"id"  validate:"required"`    // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                     // 厂商
	Region      string            `json:"region" validate:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                  // 创建时间
	ExpireAt    int64             `json:"expire_at"`                  // 过期时间
	Type        string            `json:"type" validate:"required"`   // 规格
	Name        string            `json:"name" validate:"required"`   // 名称
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
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

//给CreateAt字段注入默认值
func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

// 字段验证
var (
	validate = validator.New()
)

func (h *Host) Validate() error {
	return validate.Struct(h)
}

//对象全量更新
func (h *Host) Put(obj *Host) error {
	if obj.Id != h.Id {
		return fmt.Errorf("id not equal")
	}
	//put全量覆盖 这里替换的是值不是地址
	*h.Resource = *obj.Resource
	*h.Describe = *obj.Describe
	return nil
}

//对象的局部更新
func (h *Host) Patch(obj *Host) error {
	//注意这里指的空字段有2种含义1.没有设置值 2.设置了对应类型的零值
	//不带mergo.WithOverride函数只会将src中的非空字段覆盖到dest中对应字段的空字段上(有值不会覆盖)
	//带mergo.WithOverride会将src中的非空字段覆盖到dest中对应字段上(不管有无值都会覆盖)
	// 比如 obj.A  obj.B  只想修改obj.B该属性 那么就要将obj.B设置成非空
	return mergo.Merge(h, obj, mergo.WithOverride)
}

//主机集的结构体，用于查询主机列表时返回的参数
type HostSet struct {
	//http respond也是要marshal处理
	Total int     `json:"total"`
	Items []*Host `json:"items"`
}

func (s *HostSet) Add(h *Host) {
	s.Items = append(s.Items, h)
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

type QueryHostRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Keywords   string `json:"kws"`
}

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		//默认值
		PageSize:   10,
		PageNumber: 1,
	}
}

func (req *QueryHostRequest) OffSet() int64 {
	return int64((req.PageNumber - 1) * req.PageSize)
}

type DescribeHostRequest struct {
	Id string
}

func NewDescribeHostRequestWithId(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		Id: id,
	}
}

type HTTP_METHOD string

const (
	HTTP_METHOD_PUT   HTTP_METHOD = "put"
	HTTP_METHOD_PATCH HTTP_METHOD = "patch"
)

//用于接收put/patch请求对象
type UpdateHostRequest struct {
	HTTP_METHOD
	*Host
}

func NewPutUpdateHostRequestWithId(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		HTTP_METHOD: HTTP_METHOD_PUT,
		Host:        h,
	}
}

func NewPatchUpdateHostRequestWithId(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		HTTP_METHOD: HTTP_METHOD_PATCH,
		Host:        h,
	}
}

//用于接收patch请求对象
func NewPatchUpdateHostRequest(id string) *UpdateHostRequest {
	return nil
}

type DeleteHostRequest struct {
	//通过Id删除资源
	Id string
}

func NewDeleteHostRequestWithId(id string) *DeleteHostRequest {
	return &DeleteHostRequest{
		Id: id,
	}
}
