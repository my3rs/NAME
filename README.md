# NAME - Go博客内容管理系统

一个基于Go语言和Iris框架开发的博客内容管理系统，提供完整的REST API接口，支持博客文章、页面、用户管理、评论系统和文件上传等功能。

## 功能特性

- 📝 **内容管理**: 支持博客文章、静态页面和微博客（digu）
- 👥 **用户系统**: 用户注册、登录、角色权限管理
- 💬 **评论系统**: 支持内容评论和回复
- 🏷️ **标签分类**: 灵活的标签和分类管理
- 📁 **文件上传**: 支持附件上传和管理
- 🔐 **JWT认证**: 基于RSA密钥的JWT token认证
- 🌐 **RESTful API**: 完整的REST API接口
- 📱 **跨域支持**: 配置CORS中间件
- 🗃️ **多数据库**: 支持SQLite（开发）和PostgreSQL（生产）
- 🔧 **配置管理**: 灵活的TOML配置文件

## 技术栈

- **后端框架**: Iris v12
- **数据库ORM**: GORM
- **身份认证**: JWT (自定义fork)
- **配置管理**: Viper
- **Markdown处理**: Blackfriday v2
- **HTML安全**: Bluemonday
- **测试框架**: Testify
- **数据库**: SQLite, PostgreSQL
- **容器化**: Docker

## 快速开始

### 环境要求

- Go 1.23+
- SQLite3 (开发环境)
- PostgreSQL (生产环境，可选)

### 安装和运行

1. **克隆项目**
```bash
git clone <repository-url>
cd NAME
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置应用**
```bash
# 复制配置文件模板
cp conf/example.toml name.toml

# 编辑配置文件
vim name.toml
```

4. **运行应用**
```bash
# 开发模式
go run main.go

# 或者构建后运行
go build -o NAME main.go
./NAME
```

5. **访问应用**
```
http://localhost:8000
```

### Docker部署

```bash
# 构建镜像
docker build -t name .

# 运行容器
docker run -p 8000:8000 name
```

## 配置说明

### 配置文件优先级

系统按以下顺序查找配置文件：

1. 环境变量 (前缀: `NAME_`)
2. `./name.toml`
3. `./config/name.toml`
4. `./bin/config/name.toml`
5. `$HOME/.name/name.toml`

### 主要配置项

```toml
# 服务器配置
ROOT_URL = "http://localhost"
PORT = 8000
MODE = "development"  # development 或 production
STATIC_PATH = "/static"
DATA_PATH = "/path/to/data"

# 数据库配置
[database]
HOST = "localhost"
PORT = "5432"
NAME = "name"
USER = "postgres"
PASSWORD = "your_password"

# JWT配置
[jwt]
ACCESS_TOKEN_MAX_AGE = 3600     # 访问令牌过期时间（秒）
REFRESH_TOKEN_MAX_AGE = 86400   # 刷新令牌过期时间（秒）
SECRET_KEY = "your_secret_key"
PRIVATE_KEY = "/path/to/rsa_private_key.pem"
PUBLIC_KEY = "/path/to/rsa_public_key.pem"
```

## API接口

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/refresh` - 刷新令牌

### 内容管理
- `GET /api/v1/posts` - 获取文章列表（公开）
- `POST /api/v1/posts` - 创建文章（需认证）
- `PUT /api/v1/posts/:id` - 更新文章（需认证）
- `DELETE /api/v1/posts/:id` - 删除文章（需认证）

### 页面管理
- `GET /api/v1/pages` - 获取页面列表
- `POST /api/v1/pages` - 创建页面（需认证）

### 分类和标签
- `GET /api/v1/categories` - 获取分类列表
- `GET /api/v1/tags` - 获取标签列表

### 评论系统
- `GET /api/v1/comments` - 获取评论列表
- `POST /api/v1/comments` - 创建评论

### 文件上传
- `POST /api/v1/attachments` - 上传文件（需认证）

### 用户管理
- `GET /api/v1/users` - 获取用户列表（需认证）
- `POST /api/v1/users` - 创建用户（需认证）

### 系统设置
- `GET /api/v1/settings` - 获取系统设置
- `GET /api/v1/meta` - 获取博客元数据

## 项目结构

```
NAME/
├── main.go              # 应用入口
├── conf/                # 配置管理
│   └── example.toml     # 配置文件模板
├── controller/          # 控制器层
├── service/            # 业务逻辑层
├── model/              # 数据模型
├── database/           # 数据库连接
├── middleware/         # 中间件
├── route/              # 路由定义
├── utils/              # 工具函数
├── customerror/        # 自定义错误
├── dict/               # 常量定义
├── test/               # 测试文件
├── bin/                # 编译产物和运行时数据
└── web/                # 静态资源
```

## 开发说明

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./test/model/
go test ./test/service/
```

### 内容类型

系统支持三种内容类型：

- **post**: 博客文章
- **page**: 静态页面  
- **digu**: 微博客/短动态

### 权限系统

- **admin**: 管理员权限，可以管理所有内容
- **user**: 普通用户权限，可以创建和管理自己的内容

### 安全特性

- JWT基于RSA密钥签名
- HTML内容自动清理和安全化
- CORS跨域请求支持
- 基于角色的访问控制
- 密码bcrypt加密存储

## 贡献指南

1. Fork本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交Issue: [GitHub Issues](https://github.com/your-username/NAME/issues)
- 邮箱: your-email@example.com

---

感谢使用NAME博客内容管理系统！