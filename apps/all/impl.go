package all

import (
	//管理所有app模块的注册 执行模块中的init方法 匿名导入的方式将app模块注册到Ioc容器层中
	_ "github.com/xiaoweize/api-demo/apps/book/impl"
	_ "github.com/xiaoweize/api-demo/apps/host/http"
	_ "github.com/xiaoweize/api-demo/apps/host/impl"
)
