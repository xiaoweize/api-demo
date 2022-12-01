#注意报错内容跟每一行执行有关系
#-----定义变量名------
#!!!注意每行后面不要有空格不然也会被解析到变量里面保存
PROJECT_NAME=api-demo#项目名称  
MAIN_FILE=main.go#入口文件
PKG := "github.com/xiaoweize/$(PROJECT_NAME)"
MOD_DIR := $(shell go env GOMODCACHE)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
##用于生成version信息
BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "${PKG}/version"
#脚手架生成的应用依赖的protobuf
MCUBE_MODULE := "github.com/infraboard/mcube"
MCUBE_VERSION :=$(shell go list -m ${MCUBE_MODULE} | cut -d' ' -f2)
MCUBE_PKG_PATH := ${MOD_DIR}/${MCUBE_MODULE}@${MCUBE_VERSION}

.PHONY: all dep lint vet test test-coverage build clean

all: build


#make dep相当于执行go mod tidy
# @符号表示不输出具体命令
dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt


install: ## Install depence go package
	@go install github.com/infraboard/mcube/cmd/mcube@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#安装grpc插件
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#安装标签注入插件
	@go install github.com/favadi/protoc-go-inject-tag@latest

#将依赖的公共库拷贝到本项目
pb: ## Copy mcube protobuf files to common/pb
	@mkdir -pv common/pb/github.com/infraboard/mcube/pb
	@cp -r ${MCUBE_PKG_PATH}/pb/* common/pb/github.com/infraboard/mcube/pb
	@sudo rm -rf common/pb/github.com/infraboard/mcube/pb/*/*.go


gen: ## Init Service
# 编译成golang文件
	@protoc -I=. -I=common/pb --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto
	@go fmt ./...
# 标签注入
	@protoc-go-inject-tag -input=apps/*/*.pb.go
# 为枚举类型添加方法，会生成一个新的go文件
	@mcube generate enum -p -m apps/*/*.pb.go

#编译生成文件,先执行dep指令,注意go build后面的参数注入到version文件的变量中 没有注入GIT_TAG参数 
build: dep ## Build the binary file
	@go build -ldflags "-s -w" -ldflags " -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/$(PROJECT_NAME) $(MAIN_FILE)

#跨平台编译
linux: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/demo-api $(MAIN_FILE)

run: # Run Develop server
	@go run $(MAIN_FILE) start -f etc/demo.toml

clean: ## Remove previous build
	@rm -f dist/*


#make help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'