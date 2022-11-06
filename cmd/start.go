/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xiaoweize/api-demo/apps/host/http"
	"github.com/xiaoweize/api-demo/apps/host/impl"
	"github.com/xiaoweize/api-demo/conf"
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
			fmt.Println(err)
		}
		//创建Host Service实例,同时初始化了全局db实例
		svc := impl.NewHostServiceImpl()
		// 通过host api handler 提供http RestFul接口
		api := http.NewHostHTTPHandler(svc)
		// 提供一个gin router
		g := gin.Default()
		//注册路由
		api.Registry(g)
		//这里也可以直接通过engine实现
		// g.POST("/hosts", func(ctx *gin.Context) {
		// 	ins := host.NewHost()
		// 	//解析用户传递进来的参数到ins实例上
		// 	if err := ctx.Bind(&ins); err != nil {
		// 		fmt.Println("bind faild!")
		// 		response.Failed(ctx.Writer, err)
		// 	} else {
		// 		fmt.Println("bind success!")
		// 		//如果绑定成功，调用业务接口-创建hosts主机  从这里开始传递context
		// 		ins, err = svc.CreateHost(ctx.Request.Context(), ins)
		// 		//创建失败返回错误信息
		// 		if err != nil {
		// 			response.Failed(ctx.Writer, err)
		// 		} else {
		// 			response.Success(ctx.Writer, ins)
		// 		}
		// 	}
		// })

		return g.Run(conf.C().App.HttpAddr())
	},
}

func init() {
	startCmd.PersistentFlags().StringVarP(&confFile, "configfile", "f", "etc/demo.toml", "config file path")
	rootCmd.AddCommand(startCmd)
}
