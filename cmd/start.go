/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"
	"github.com/xiaoweize/api-demo/apps"
	"github.com/xiaoweize/api-demo/conf"
	"github.com/xiaoweize/api-demo/protocol"

	//注册所有的app模块
	_ "github.com/xiaoweize/api-demo/apps/all"
)

var (
	confFile string
)

// 程序启动的组装
var startCmd = &cobra.Command{
	Use:   "start", //在这里指定command名
	Short: "start demo host-api",
	Long:  `start demo host-api impl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//初始化全局配置
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			//输出加载失败错误，有默认值不影响程序继续运行
			return err
		}

		//初始化全局日志配置
		loadGlobalLogger()

		//创建Host Service实例,同时初始化了全局db实例
		// svc := impl.NewHostServiceImpl()
		//注册HostService实例到Ioc层容器中
		// apps.HostService = impl.NewHostServiceImpl()

		//========将app模块注册和初始化逻辑分开======
		//使用import _ "github.com/xiaoweize/api-demo/apps/all"来管理所有需要注册到Ioc容器层的服务
		//import虽然完成了app的注册，但是还需要初始化
		//如何执行HostService的config方法 apps.HostService是一个接口类型，需要将其断言成实例对象才能使用方法
		//使用apps.Init()初始化所有注册到Ioc容器中的服务，其实就是调用服务中的Config方法来完成初始化 所以每个app模块对象要有Config方法来完成初始化
		apps.InitImpl()

		// g := gin.Default()
		//将所有http handler注册到IOC中
		// apps.InitGin(g)
		//=======http需要获取到app服务模块对外暴露的接口=======
		// 通过host api handler 提供http RestFul接口
		// api := http.NewHostHTTPHandler()
		//api从Ioc层中获取依赖
		//由此api对象和impl.NewHostServiceImpl对象之间的依赖关系通过Ioc容器层解除了
		// api.Config()

		// 提供一个gin router
		// g := gin.Default()
		//路由注册
		// api.Registry(g)

		svc := Newmanager()
		ch := make(chan os.Signal, 1)
		defer close(ch)
		//处理2 3 15信号量 将信号发送到channel中
		//注意9的信号压根收不到直接强制kill掉
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		go svc.WaitStop(ch)
		//grpc在后台启动
		go svc.grpc.Start()

		return svc.Start()
	},
}

//程序启动时需要关注如下问题
// 1. http API, Grpc API 需要启动, 消息总线也需要监听, 比如负责注册与配置,  这些模块都是独立
//    都需要在程序启动时，进行启动, 都写在直接写到start中会膨胀到不易维护
// 2. 服务的优雅关闭怎么办? 外部都会发送一个Terminal(中断)信号给程序, 程序时需要处理这个信号
//    需要实现程序优雅关闭: 由先后顺序 (从外到内完成资源的释放逻辑处理)
//    1. api 层的关闭 (HTTP, GRPC)
//    2. 关闭消息总线
//    3. 关闭数据库连接
//    4. 如果使用了注册中心, 还要在注册中心完成注销操作
//    5. 退出完毕
//通过manage 用于管理所有需要启动的服务如http/grpc  位于跟目录下protocol目录 用于对外暴露的协议通常有http grpc
type manager struct {
	http *protocol.HttpService
	grpc *protocol.GrpcService
	l    logger.Logger
}

func Newmanager() *manager {
	return &manager{
		protocol.NewHttpService(),
		protocol.NewGrpcService(),
		zap.L().Named("CLI"),
	}
}

func (m *manager) Start() error {
	return m.http.Start()
}

//处理来自外部的终端信号
//键盘Ctrl + C 会发送INT(2)信号量
//键盘Ctrl + \会发送QUIT(3)信号量
//kll PID命令会发送TERM(15)信号量给程序
func (m *manager) WaitStop(ch <-chan os.Signal) {
	for v := range ch {
		//可以针对不同的信号量做处理，这里全部按关闭处理
		switch v {
		default:
			m.l.Infof("received signal:%s", v)
			//先关内部的grpc调用
			m.grpc.Stop()
			//再关外部的http
			m.http.Stop()
		}
	}
}

//全局Logger对象初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string    //Logger初始化消息
		level      zap.Level //日志级别
	)

	// 加载全局配置conf里面的日志配置来初始化全局Logger对象，注意要先初始化全局conf配置
	lc := conf.C().Log
	// 解析并设置日志lever级别——日志Level级别：DebugLevel: "debug",InfoLevel:  "info",WarnLevel:  "warning",
	// 					ErrorLevel: "error",FatalLevel: "fatal",PanicLevel: "panic"
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		//如果设置了不支持的级别，就将级别置为InfoLevel
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	// 使用默认配置初始化Logger的全局配置
	//默认配置:日志级别InfoLevel,日志输出到文件并按照日志大小轮转等
	//有时候默认配置也满足基本需求
	zapConfig := zap.DefaultConfig()

	//===在默认配置上面修改成用户的自定义配置===
	//配置日志的Level级别
	zapConfig.Level = level

	//配置程序启动时是否需要日志的轮转即是否生成新的日志文件
	zapConfig.Files.RotateOnStartup = false

	// 配置日志的输出方式
	switch lc.To {
	//使用枚举的方式来定义输出，方便阅读代码
	case conf.ToStdout:
		zapConfig.ToStderr = true // 把日志打印到控制台
		zapConfig.ToFiles = false // 关闭日志写文件
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"  //日志文件名
		zapConfig.Files.Path = lc.PathDir //日志文件路径
	}

	// 配置日志的输出格式为json
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}

	// 把配置应用到全局Logger
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}

	//打印日志初始化信息
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	startCmd.PersistentFlags().StringVarP(&confFile, "configfile", "f", "etc/demo.toml", "config file path")
	rootCmd.AddCommand(startCmd)
}
