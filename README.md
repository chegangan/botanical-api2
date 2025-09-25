# Botanical API

## 目录

- [项目简介](#项目简介)
- [主要功能](#主要功能)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
  - [环境要求](#环境要求)
  - [安装与配置](#安装与配置)
  - [运行项目](#运行项目)
  - [API文档](#api文档)
- [API 概览](#api-概览)
- [贡献](#贡献)
- [许可证](#许可证)

## 项目简介

Botanical API 是一个基于 Go 语言 Gin 框架构建的 RESTful API 服务，专为植物园管理系统设计。项目采用 GORM 作为 ORM 与数据库交互，并遵循了清晰的 **模型-仓库-服务-API** 分层架构。它提供了完整的用户认证、植物信息管理、园区管理、资讯发布等功能。

## 主要功能

- **用户管理**: 支持用户注册、登录、信息查询和头像上传。
- **认证与授权**: 使用 JWT (JSON Web Tokens) 进行无状态认证，并包含管理员权限校验的中间件。
- **植物管理**: 提供植物信息的增删改查，支持按分类、按园区查询。
- **分类与园区**: 管理植物的分类以及所属的园区信息。
- **资讯系统**: 发布和管理植物相关的资讯，并统计浏览量。
- **用户反馈**: 用户可以提交反馈，管理员可以查看所有反馈。
- **API 文档**: 集成 Swagger，自动生成交互式 API 文档。

## 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin)
- **数据库 ORM**: [GORM](https://gorm.io/)
- **数据库**: MySQL
- **配置管理**: [go-ini/ini](https://github.com/go-ini/ini)
- **API 文档**: [Swaggo](https://github.com/swaggo/swag)
- **认证**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)

## 项目结构

```
.
├── cmd/server/               # 应用主入口
│   ├── main.go               # 程序启动文件
│   └── server.go             # 服务器初始化与配置
├── configs/                  # 配置文件目录
│   └── app.ini
├── docs/                     # Swagger 自动生成的文档
├── internal/                 # 内部业务逻辑
│   ├── api/                  # API 层 (路由、处理器、中间件)
│   │   ├── handlers/         # HTTP 请求处理器
│   │   ├── middleware/       # 中间件 (JWT, CORS, Admin)
│   │   └── router.go         # 路由注册
│   ├── dto/                  # 数据传输对象 (Data Transfer Objects)
│   ├── models/               # GORM 数据库模型
│   ├── repository/           # 数据仓库层 (封装数据库操作)
│   └── service/              # 服务层 (封装业务逻辑)
├── pkg/                      # 公共库/工具包
│   ├── app/                  # 应用相关的封装 (响应、分页)
│   ├── e/                    # 错误码定义
│   ├── jwt/                  # JWT 工具
│   ├── setting/              # 配置加载
│   └── utils/                # 通用工具函数
├── go.mod                    # Go 模块依赖
└── README.md                 # 项目说明文档
```

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+

### 安装与配置

1.  **克隆项目**
    ```sh
    git clone <your-repository-url>
    cd botanical-api2
    ```

2.  **安装依赖**
    ```sh
    go mod tidy
    ```

3.  **配置应用**
    - 复制或重命名 `configs/app.ini` 文件。
    - 修改 `app.ini` 中的数据库连接信息 (`[database]`) 和服务器端口 (`[server]`) 等配置。

4.  **数据库**
    - 确保您的 MySQL 服务正在运行。
    - 根据 `app.ini` 中的配置，手动创建数据库。
    - GORM 会在应用首次启动时自动迁移数据表。

### 运行项目

- **开发模式**
  在项目根目录下运行：
  ```sh
  go run cmd/server/main.go
  ```
  服务器默认将在 `http://localhost:8000` 启动。

- **构建可执行文件**
  您可以为特定平台构建项目，例如为 Windows：
  ```sh
  GOOS=windows GOARCH=amd64 go build -o botanical-api.exe cmd/server/main.go cmd/server/server.go
  ```

### API文档

项目启动后，您可以访问 Swagger UI 查看所有 API 接口的详细信息和进行在线测试。

- **Swagger UI 地址**: [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)

## API 概览

API 的基础路径为 `/api/v1`。

| 模块         | 路径                  | 方法   | 描述                       | 需要认证 |
| :----------- | :-------------------- | :----- | :------------------------- | :------- |
| **认证管理** | `/auth/register`      | `POST` | 用户注册                   | 否       |
|              | `/auth/login`         | `POST` | 用户登录                   | 否       |
| **用户管理** | `/users/{id}`         | `GET`  | 获取指定用户信息           | 是       |
|              | `/avatar`             | `POST` | 上传用户头像               | 是       |
| **植物**     | `/plants/{id}`        | `GET`  | 获取植物详情               | 否       |
| **植物分类** | `/plant-classes`      | `GET`  | 获取所有植物分类           | 否       |
| **园区**     | `/parks`              | `GET`  | 获取所有园区列表           | 否       |
| **植物资讯** | `/notices`            | `GET`  | 获取资讯列表               | 否       |
|              | `/notices/{id}`       | `GET`  | 获取资讯详情               | 否       |
| **用户反馈** | `/feedback`           | `POST` | 创建用户反馈               | 是       |
| **管理员**   | `/admin/feedbacks`    | `GET`  | 获取所有用户反馈（管理员） | 是 (管理员) |

*这是一个简化的列表，完整的 API 请参阅 Swagger 文档。*

## 贡献

欢迎任何形式的贡献！如果您有任何建议或发现问题，请随时提交 Issue 或 Pull Request。

## 许可证

本项目采用 MIT 许可证。详情请查看 `LICENSE` 文件。