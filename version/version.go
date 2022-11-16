package version

import (
	"fmt"
)

const (
	// ServiceName 服务名称
	ServiceName = "api-demo"
)

//这里的版本参数 通过宏, 编译时注入, 一般的编译型语言都支持  即执行make命令时注入进去 注意没有注入GIT_TAG参数
var (
	GIT_TAG    string
	GIT_COMMIT string
	GIT_BRANCH string
	BUILD_TIME string
	GO_VERSION string
)

// FullVersion show the version info
func FullVersion() string {

	version := fmt.Sprintf("Version   : %s\nBuild Time: %s\nGit Branch: %s\nGit Commit: %s\nGo Version: %s\n", GIT_TAG, BUILD_TIME, GIT_BRANCH, GIT_COMMIT, GO_VERSION)
	return version
}

// Short 版本缩写
func Short() string {
	return fmt.Sprintf("%s[%s %s]", GIT_TAG, BUILD_TIME, GIT_COMMIT)
}
