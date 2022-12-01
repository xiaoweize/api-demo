package apps

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

//Ioc容器层 管理所有服务的实例 所有的依赖层
//凡是有接口申明的对象都放在这里面

//1.Host Service实例必须注册过来，HostService才会有具体的实例，服务启动的时候完成注册

//2.HTTP暴露模块依赖Ioc里面的HostService
var (
	//如果有40个app模块 要定义40个app模块变量，怎么办
	//使用空接口any+断言进行抽象
	// HostService host.Service
	//维护当前所有app模块名与接口映射
	implApps = map[string]ImplService{}
	ginApps  = map[string]GinService{}
	grpcApps = map[string]GrpcService{}
)

//定义所有app模块接口 注意与host.ImplService的区别 后者是单个app模块的接口 前者是所有app模块的接口
type ImplService interface {
	Name() string
	Config()
}

//注册app模块时使用
func RegistryImpl(svc ImplService) {
	//将app模块接口注册到svcs的map中
	if _, ok := implApps[svc.Name()]; ok {
		//如果已经注册就panic
		panic(fmt.Errorf("%s has registried!", svc.Name()))
	}
	implApps[svc.Name()] = svc
}

//初始化注册到Ioc容器里面所有服务
func InitImpl() {
	for _, v := range implApps {
		v.Config()
	}
}

//get的是一个impl的实例 从维护的map里面去拿
//返回空接口，使用时由使用方进行断言
func GetImpl(name string) any {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}
	return nil
}

//注册由gin编写的handler
type GinService interface {
	Registry(r gin.IRouter)
	Name() string
	Config()
}

func RegistryGin(svc GinService) {
	//将gin编写的模块接口注册到svcs的map中
	if _, ok := ginApps[svc.Name()]; ok {
		//如果已经注册就panic
		panic(fmt.Errorf("%s has registried!", svc.Name()))
	}
	ginApps[svc.Name()] = svc
}

//http对象的初始化
func InitGin(r gin.IRouter) {
	//先初始化好所有的对象
	for _, v := range ginApps {
		v.Config()
	}
	//再完成http handler的注册
	for _, v := range ginApps {
		v.Registry(r)
	}
}

//获取已加载完成的gin apps
func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}
	return
}

//grpc service
type GrpcService interface {
	Registry(r *grpc.Server)
	Name() string
	Config()
}

//grpc service注册
func RegistryGrpc(svc GrpcService) {
	//将gin编写的模块接口注册到svcs的map中
	if _, ok := grpcApps[svc.Name()]; ok {
		//如果已经注册就panic
		panic(fmt.Errorf("%s has registried!", svc.Name()))
	}
	grpcApps[svc.Name()] = svc
}

//grpc初始化
func InitGrpc(g *grpc.Server) {
	//先初始化好所有的对象
	for _, v := range grpcApps {
		v.Config()
	}
	//再完成grpc的注册
	for _, v := range grpcApps {
		v.Registry(g)
	}
}

//获取已加载完成的grpc service
func LoadedGrpcApps() (names []string) {
	for k := range grpcApps {
		names = append(names, k)
	}
	return
}
