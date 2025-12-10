# 🚀 go-blog 项目文档

一个基于 Go 语言的博客系统，采用分层架构设计，支持用户认证、文章管理、评论功能，具备完善的配置、错误处理和日志系统。

---

## 📂 项目结构

```bash
go-blog/
├── cmd/                    # 主程序入口
│   └── main.go
├── config/                 # 配置模块
│   └── database.go         # 数据库配置
├── errors/                 # 错误码与错误处理
│   └── errors.go
├── logger/                 # 日志模块
│   ├── logger.go           # 日志接口
│   └── zap_logger.go       # Zap 日志实现
├── middleware/             # Gin 中间件
│   ├── auth.go             # 认证中间件
│   └── logger.go           # 请求日志中间件
├── models/                 # 数据模型
│   ├── comment.go          # 评论模型
│   ├── post.go             # 文章模型
│   └── user.go             # 用户模型
├── routers/                # 路由模块
│   └── routers.go          # 路由注册
├── handlers/                 # 业务逻辑层
│   ├── auth.go             # 认证逻辑
│   ├── comment.go          # 评论逻辑
│   └── post.go             # 文章逻辑
├── utils/                  # 工具类
│   ├── jwt.go              # JWT 生成与解析
│   ├── page.go             # 分页工具
│   ├── response.go         # 统一响应格式
│   └── validationField.go  # 字段验证工具
├── .env                    # 环境变量配置
└── README.md               # 项目说明
```

## ⚙️ 技术栈
### Web 框架：Gin（高性能 HTTP 框架）
### ORM：GORM（数据库 ORM 工具）
### 日志：Zap（结构化日志库）
### 认证：JWT（基于 utils/jwt.go 实现）
### 配置管理：.env + config/database.go
### 中间件：Auth、Logger
### 分页工具：utils/page.go（支持标准分页参数处理）
### 错误处理：自定义错误码与统一响应

## 启动服务
```bash
go run cmd/main.go
```

## 📡 核心功能
### ✅ 用户认证
#### 登录 / 注册（server/auth.go）
#### JWT 生成与验证（utils/jwt.go）
#### 认证中间件（middleware/auth.go）
### ✅ 文章管理
#### 文章增删改查（server/post.go）
#### 分页列表（utils/page.go）
#### 响应格式统一（utils/response.go）
### ✅ 评论功能
#### 评论发布与查询（server/comment.go）
#### 关联文章与用户（models/comment.go）
### ✅ 安全与日志
#### CORS 跨域支持（middleware/cors.go）
#### 请求日志记录（middleware/logger.go + logger/zap_logger.go）
#### 错误码统一管理（errors/errors.go）