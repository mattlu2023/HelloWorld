# 广告分析 BI 系统 - Go 后端

## 项目说明

这是广告分析 BI 系统的后端服务，使用 Go + Gin Framework 开发。

## 功能特性

- ✅ 用户认证（登录/注册）
- ✅ 广告活动管理（CRUD）
- ✅ 广告单元管理（CRUD）
- ✅ 数据统计接口（概览、趋势、漏斗）
- ✅ 用户行为分析
- ✅ 报表导出

## 环境要求

- Go 1.21+
- MySQL 8.0+

## 安装步骤

### 1. 安装依赖

```bash
cd backend-go
go mod download
```

### 2. 配置环境变量

创建 `.env` 文件或设置环境变量：

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ad_bi_system
PORT=8080
```

### 3. 初始化数据库

```bash
mysql -u root -p < database/init.sql
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## API 接口

### 认证接口
- `POST /api/v1/login` - 用户登录
- `POST /api/v1/register` - 用户注册

### 广告活动
- `GET /api/v1/campaigns` - 获取活动列表
- `GET /api/v1/campaigns/:id` - 获取活动详情
- `POST /api/v1/campaigns` - 创建活动
- `PUT /api/v1/campaigns/:id` - 更新活动
- `DELETE /api/v1/campaigns/:id` - 删除活动

### 广告单元
- `GET /api/v1/ad-units` - 获取单元列表
- `GET /api/v1/ad-units/:id` - 获取单元详情
- `POST /api/v1/ad-units` - 创建单元
- `PUT /api/v1/ad-units/:id` - 更新单元
- `DELETE /api/v1/ad-units/:id` - 删除单元

### 数据统计
- `GET /api/v1/stats/overview` - 数据概览
- `GET /api/v1/stats/daily-trend` - 每日趋势
- `GET /api/v1/stats/funnel` - 转化漏斗
- `GET /api/v1/stats/campaign/:id` - 活动统计

### 用户行为
- `GET /api/v1/user-actions` - 用户行为列表
- `GET /api/v1/user-actions/analysis` - 行为分析

### 报表
- `GET /api/v1/reports/export` - 导出报表

## 请求示例

### 登录
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 获取活动列表
```bash
curl -X GET http://localhost:8080/api/v1/campaigns \
  -H "Authorization: Bearer your_token"
```

## 项目结构

```
backend-go/
├── main.go              # 入口文件
├── config/
│   └── database.go      # 数据库配置
├── handlers/
│   ├── auth.go          # 认证处理器
│   ├── campaign.go      # 活动处理器
│   ├── adunit.go        # 广告单元处理器
│   └── stats.go         # 统计处理器
├── middleware/
│   └── auth.go          # 认证中间件
└── models/              # 数据模型
```

## 开发说明

### 添加新的 Handler

1. 在 `handlers/` 目录创建新的处理器文件
2. 在 `main.go` 中注册路由
3. 如需认证，使用 `auth.Use(middleware.AuthMiddleware())`

### 数据库操作

使用标准库 `database/sql`，通过 `config.GetDB()` 获取连接

## 注意事项

- 生产环境请使用 JWT 进行认证
- 密码应该使用 bcrypt 等加密存储
- 建议添加请求限流和日志记录
- 数据库连接池需要优化配置
