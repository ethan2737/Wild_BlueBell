# Wild Bluebell 项目 Makefile
# 用于自动化构建、运行和清理等常见开发任务
#
# ==================== 使用场景 ====================
# 本 Makefile 主要用于以下场景：
#   1. 打包 Linux 部署产物 - make build 生成 Linux amd64 二进制
#   2. CI/CD 流水线 - 自动化构建和发布
#   3. 团队协作 - 统一构建命令，避免不同系统的构建差异
#
# 本地 Windows 开发通常不需要使用，直接用以下命令即可：
#   go run main.go conf/config.yaml
#
# ==================== 使用方法 ====================
#   make          - 格式化代码并构建二进制文件
#   make build    - 仅构建 Linux 二进制文件
#   make run      - 直接运行主程序（本地开发用 go run 即可）
#   make gotool   - 运行代码检查工具 (fmt + vet)
#   make clean    - 清理构建产物
#   make help     - 显示帮助信息


# PHONY 声明的是伪目标（phony targets），告诉 Make 该名称不对应实际文件。
.PHONY: all build run gotool clean help

# BINARY 是一个变量定义，用于存储输出的二进制文件名 (生成的二进制文件名)
BINARY="wild_bluebell"

# 默认目标：先运行代码格式化/检查，再构建
all: gotool build

# 构建目标：编译 Go 代码为 Linux 64位二进制文件
#   - CGO_ENABLED=0    禁用 CGO，静态编译
#   - GOOS=linux        目标操作系统为 Linux
#   - GOARCH=amd64      目标架构为 x86_64
#   - -ldflags "-s -w"  移除调试信息和符号表，减小二进制体积
#   - 输出到 ./bin/ 目录
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/${BINARY}

# 运行目标：直接启动主程序，需要指定配置文件
run:
	@go run ./main.go conf/config.yaml

# gotool 目标：运行 Go 代码质量工具
#   - go fmt   格式化代码
#   - go vet   检查代码潜在问题
gotool:
	go fmt ./
	go vet ./

# clean 目标：清理构建产物
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

# help 目标：显示 Makefile 使用帮助信息
help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
