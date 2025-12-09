# APIOK Admin

APIOK Admin 是一个基于 Go 和 Vue.js 的 API 网关管理后台系统，提供可视化的服务、路由、上游、插件和证书管理功能。

## 功能特性

- **服务管理**：支持 HTTP/HTTPS 服务配置，支持多域名、多端口
- **路由管理**：灵活的路由规则配置，支持路径匹配、请求方法过滤
- **上游管理**：支持多种负载均衡算法，健康检查配置
- **插件系统**：丰富的插件生态，包括限流、认证、CORS 等
- **证书管理**：支持 SSL/TLS 证书上传和管理，Let's Encrypt 自动申请
- **LDAP 认证**：支持 LDAP 统一认证登录
- **集群管理**：支持多节点集群配置

## 技术栈

### 后端
- Go 1.23+
- Gin Web Framework
- GORM
- MySQL
- LDAP (可选)

### 前端
- Vue.js
- Ant Design Vue

## 快速开始

### 环境要求

- Go 1.23 或更高版本
- MySQL 5.7+ 或 MySQL 8.0+
- Node.js 16+ (用于前端开发)

### 安装步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd apiok-admin
```

2. **安装依赖**

```bash
go mod download
```

3. **配置数据库**

创建数据库并导入 SQL 文件：

```bash
mysql -u root -p < config/apiok.sql
```

4. **配置文件**

编辑 `config/app.yaml`，配置数据库连接和其他参数：

```yaml
database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  db_name: apiok
  username: apiok
  password: your_password

apiok:
  protocol: http
  ip: 127.0.0.1
  port: 8080
  domain: 127.0.0.1
  secret: your_secret_key
```

5. **运行项目**

```bash
go run main.go
```

或者使用 Makefile：

```bash
make run
```

服务将在 `http://localhost:3000` 启动。

## 配置说明

### 数据库配置

```yaml
database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  db_name: apiok
  username: root
  password: 123456
  max_idel_connections: 10
  max_open_connections: 100
  sql_mode: true
```

### LDAP 配置

```yaml
ldap:
  enabled: true
  host: ldap.example.com
  port: 389
  base_dn: "dc=example,dc=com"
  bind_dn: "cn=admin,dc=example,dc=com"
  bind_password: "password"
  user_filter: "(uid=%s)"
  attributes:
    name: "cn"
    email: "mail"
```

### Let's Encrypt 配置

```yaml
letsencrypt:
  enabled: true
  email: "admin@example.com"
  use_staging: false
  cert_dir: "./certs"
  renew_before_days: 30
```

## 项目结构

```
apiok-admin/
├── app/                    # 应用代码
│   ├── controllers/        # 控制器
│   ├── models/            # 数据模型
│   ├── services/          # 业务逻辑
│   ├── validators/        # 参数验证
│   ├── middlewares/       # 中间件
│   ├── packages/          # 工具包
│   └── utils/             # 工具函数
├── cores/                 # 核心模块
├── config/                # 配置文件
│   ├── app.yaml          # 应用配置
│   └── apiok.sql        # 数据库结构
├── routers/               # 路由定义
├── html/                  # 前端静态文件
├── main.go               # 入口文件
└── go.mod               # Go 模块定义
```

## API 文档

### 认证

所有需要认证的接口都需要在请求头中携带 `auth-token`：

```
auth-token: <your_token>
```

### 主要接口

- `POST /admin/user/login` - 用户登录
- `POST /admin/user/logout` - 用户登出
- `GET /admin/service/list` - 获取服务列表
- `POST /admin/service/add` - 创建服务
- `GET /admin/router/list` - 获取路由列表
- `POST /admin/router/add` - 创建路由
- `GET /admin/upstream/list` - 获取上游列表
- `POST /admin/upstream/add` - 创建上游

更多接口请参考代码中的路由定义。

## 开发

### 构建

```bash
go build -o apiok-admin main.go
```

### 运行测试

```bash
go test ./...
```

### 代码规范

项目遵循 Go 标准代码规范，建议使用 `gofmt` 和 `golint` 进行代码格式化。

## 许可证

查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献

欢迎提交 Issue 和 Pull Request。

## 联系方式

如有问题或建议，请通过 Issue 反馈。

