.PHONY: build
build:
	go build -o apiok-admin main.go

.PHONY: build-all
build-all: dashboard-build
	@rm -rf html 2>/dev/null || cmd /c "if exist html rmdir /s /q html" 2>nul || true
	@cp -r apiok-dashboard/html html 2>/dev/null || cmd /c "xcopy /E /I /Y apiok-dashboard\html html" 2>nul || true
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

.PHONY: clean
clean:
	@rm -f apiok-admin apiok-admin.exe apiok-admin_linux_amd64 apiok-admin_linux_arm64 2>/dev/null || true
	@rm -rf html 2>/dev/null || cmd /c "if exist html rmdir /s /q html" 2>nul || true
	@echo "已清理构建产物"

.PHONY: install
install: build
	@if [ -f apiok-admin ]; then \
		if [ "$$(uname -s)" = "Linux" ] || [ "$$(uname -s)" = "Darwin" ]; then \
			sudo mkdir -p /opt/apiok-admin && \
			sudo cp apiok-admin /opt/apiok-admin/apiok-admin && \
			sudo chmod +x /opt/apiok-admin/apiok-admin && \
			sudo cp -r config /opt/apiok-admin/ && \
			echo "已安装到 /opt/apiok-admin/apiok-admin"; \
			echo "配置文件已复制到 /opt/apiok-admin/config/"; \
		else \
			echo "Windows 平台请手动复制 apiok-admin.exe 到 PATH 目录"; \
		fi \
	else \
		echo "请先运行 make build 构建项目"; \
	fi

.PHONY: install-supervisor
install-supervisor:
	@if [ "$$(uname -s)" = "Linux" ]; then \
		if [ -f supervisor/apiok-admin.ini ]; then \
			sudo cp supervisor/apiok-admin.ini /etc/supervisord.d/apiok-admin.ini && \
			sudo chmod 644 /etc/supervisord.d/apiok-admin.ini && \
			echo "supervisor配置文件已安装到 /etc/supervisord.d/apiok-admin.ini"; \
			echo "请运行以下命令重新加载supervisor配置:"; \
			echo "  sudo supervisorctl reread"; \
			echo "  sudo supervisorctl update"; \
		else \
			echo "supervisor配置文件不存在: supervisor/apiok-admin.ini"; \
		fi \
	else \
		echo "supervisor安装仅支持Linux平台"; \
	fi

.PHONY: help
help:
	@echo "make build : 仅根据当前平台编译"
	@echo "make build-all : 编译 linux/amd64、linux/arm64（包含前端构建）"
	@echo "make dashboard-build : 构建 apiok-dashboard 前端项目"
	@echo "make run : 直接运行 Go 代码"
	@echo "make clean : 清理构建产物"
	@echo "make install : 安装到 /opt/apiok-admin（Linux/macOS）"
	@echo "make install-supervisor : 安装supervisor配置文件到 /etc/supervisord.d/"
