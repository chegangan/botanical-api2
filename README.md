# README.md

# Botanical API

## 项目简介

Botanical API 是一个基于 Go 语言的 RESTful API，使用 Gin 框架和 Gorm ORM 进行数据库操作。该项目实现了用户模型的增删改查接口，并采用模型-仓储-服务-API 架构。

## 项目结构

```
botanical-api2
├── conf
│   └── app.ini              # 配置文件，包含应用的基本设置，如运行模式、数据库连接信息等。
├── middleware
│   ├── cors.go              # 处理跨域请求的中间件。
│   └── jwt.go               # 处理JWT认证的中间件。
├── models
│   ├── init.go              # 初始化数据库连接和模型。
│   └── user.go              # 用户模型，定义用户的结构体及其相关的数据库操作。
├── pkg
│   ├── app
│   │   ├── request.go       # 定义请求的结构体和验证逻辑。
│   │   └── response.go      # 定义响应的结构体和格式化逻辑。
│   ├── setting
│   │   └── setting.go       # 加载和解析配置文件的逻辑。
│   ├── util
│   │   └── pagination.go     # 处理分页逻辑的工具函数。
│   └── e
│       ├── code.go          # 定义错误代码的常量。
│       └── msg.go           # 定义错误消息的常量。
├── repository
│   └── user_repository.go    # 用户数据访问层，封装与数据库的交互逻辑。
├── routers
│   ├── api
│   │   ├── v1
│   │   │   └── user.go       # 用户相关的API路由，定义用户的增删改查接口。
│   │   └── auth.go           # 认证相关的API路由。
│   └── router.go             # 设置所有路由的入口。
├── service
│   └── user_service.go       # 用户服务层，包含业务逻辑，调用数据访问层的方法。
├── main.go                   # 应用的入口文件，启动Gin服务器并加载路由。
└── README.md                 # 项目的文档，包含使用说明和其他相关信息。
```

## 使用说明

1. **配置文件**: 修改 `conf/app.ini` 文件以设置应用的基本配置，包括数据库连接信息和运行模式。

2. **数据库**: 确保已安装并运行 MySQL 数据库，并根据配置文件中的设置创建数据库。

3. **运行项目**: 使用以下命令启动项目：
   ```
   go run main.go
   ```

4. **API 接口**:
   - **创建用户**: `POST /users`
   - **获取用户**: `GET /users/:id`
   - **更新用户**: `PUT /users/:id`
   - **删除用户**: `DELETE /users/:id`

## 依赖

- Go 1.16+
- Gin
- Gorm
- MySQL

## 贡献

欢迎任何形式的贡献！请提交问题或拉取请求。

## 许可证

本项目采用 MIT 许可证，详细信息请查看 LICENSE 文件。