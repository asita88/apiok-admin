.PHONY: build
build:
	go build -o apiok-admin main.go

.PHONY: build-all
build-all: dashboard-build
	@rm -rf html 2>/dev/null || if exist html rmdir /s /q html
	@cp -r apiok-dashboard/html html 2>/dev/null || xcopy /E /I /Y apiok-dashboard\html html
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-linux-arm64

.PHONY: build-linux-amd64
build-linux-amd64:
	@set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=amd64 && go build -o apiok-admin_linux_amd64 main.go || env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apiok-admin_linux_amd64 main.go

.PHONY: build-linux-arm64
build-linux-arm64:
	@set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=arm64 && go build -o apiok-admin_linux_arm64 main.go || env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o apiok-admin_linux_arm64 main.go

.PHONY: dashboard-build
dashboard-build:
	git submodule update --init --recursive
	cd apiok-dashboard && npm install && npm run build

.PHONY: run
run:
	@go run ./main.go

.PHONY: help
help:
	@echo "make build : 仅根据当前平台编译"
	@echo "make build-all : 编译 linux/amd64、linux/arm64（包含前端构建）"
	@echo "make dashboard-build : 构建 apiok-dashboard 前端项目"
	@echo "make run : 直接运行 Go 代码"
