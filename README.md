# 广告分析 BI 系统

一个功能完整的广告数据分析 BI 系统，包含数据看板、广告投放管理、用户行为分析和报表导出功能。

## 🏗️ 技术架构

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐      ┌──────────┐
│  Vue 3      │ ───► │  Node.js     │ ───► │  Go + Gin   │ ───► │  MySQL   │
│  前端界面   │      │  中间层      │      │  后端 API    │      │  数据库  │
└─────────────┘      └──────────────┘      └─────────────┘      └──────────┘
```

### 技术栈

- **前端**: Vue 3 + Vite + Element Plus + ECharts
- **中间层**: Node.js + Express
- **后端**: Go + Gin Framework
- **数据库**: MySQL 8.0+

## 📁 项目结构

```
ad-bi-system/
├── database/              # 数据库脚本
│   └── init.sql          # MySQL 初始化脚本
├── backend-go/           # Go 后端服务
│   ├── main.go
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   └── models/
├── middleware-node/      # Node.js 中间层
│   ├── server.js
│   └── package.json
├── frontend-vue/         # Vue 前端
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── api/
│   │   └── assets/
│   └── package.json
└── README.md
```

## 🚀 快速开始

### 1. 数据库配置

```bash
# 登录 MySQL
mysql -u root -p

# 执行初始化脚本
source database/init.sql
```

### 2. 启动 Go 后端

```bash
cd backend-go

# 安装依赖
go mod download

# 配置环境变量（可选）
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=ad_bi_system
export PORT=8080

# 运行服务
go run main.go
```

Go 后端将在 `http://localhost:8080` 启动

### 3. 启动 Node.js 中间层

```bash
cd middleware-node

# 安装依赖
npm install

# 创建环境配置
cp .env.example .env

# 运行服务
npm start
```

中间层将在 `http://localhost:3001` 启动

### 4. 启动 Vue 前端

```bash
cd frontend-vue

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端将在 `http://localhost:3000` 启动

## 📊 功能模块

### ✅ 已实现功能

1. **用户认证**
   - 登录/注册
   - JWT Token 认证
   - 路由守卫

2. **数据看板**
   - 核心指标展示（曝光、点击、转化、收入）
   - 每日趋势图表
   - 转化漏斗可视化

3. **广告活动管理**
   - 活动列表查询
   - 创建/编辑/删除活动
   - 活动状态管理

4. **广告单元管理**
   - 广告单元 CRUD
   - 多种广告类型支持
   - 投放位置管理

5. **数据统计 API**
   - 概览数据统计
   - 每日趋势分析
   - 转化漏斗数据
   - 活动维度统计

### 🔧 待完善功能

- [ ] 广告活动完整的 CRUD 界面
- [ ] 广告单元管理界面
- [ ] 用户行为分析详情页
- [ ] 报表导出功能（Excel/PDF）
- [ ] 实时数据更新
- [ ] 权限管理系统
- [ ] 数据导入功能

## 🔌 API 接口文档

### 认证接口

```bash
POST /api/v1/login
{
  "username": "admin",
  "password": "admin123"
}

POST /api/v1/register
{
  "username": "newuser",
  "password": "password123"
}
```

### 广告活动接口

```bash
GET /api/v1/campaigns           # 获取活动列表
GET /api/v1/campaigns/:id       # 获取活动详情
POST /api/v1/campaigns          # 创建活动
PUT /api/v1/campaigns/:id       # 更新活动
DELETE /api/v1/campaigns/:id    # 删除活动
```

### 数据统计接口

```bash
GET /api/v1/stats/overview      # 数据概览
GET /api/v1/stats/daily-trend   # 每日趋势
GET /api/v1/stats/funnel        # 转化漏斗
GET /api/v1/stats/campaign/:id  # 活动统计
```

## 📝 默认账号

- 用户名：`admin`
- 密码：`admin123`

## ⚙️ 环境配置

### Go 后端环境变量

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ad_bi_system
PORT=8080
```

### Node.js 中间层环境变量

```bash
PORT=3001
GO_BACKEND_URL=http://localhost:8080
NODE_ENV=development
```

### Vue 前端配置

在 `vite.config.js` 中配置代理：

```javascript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:3001',
      changeOrigin: true
    }
  }
}
```

## 🎨 界面预览

### 数据看板
- 核心指标卡片（曝光量、点击量、转化量、收入）
- 每日趋势折线图
- 转化漏斗图

### 广告活动管理
- 活动列表表格
- 活动创建/编辑表单
- 活动状态筛选

### 广告单元管理
- 广告单元列表
- 单元类型筛选
- 投放位置管理

## 🔒 安全说明

1. **生产环境必须修改默认密码**
2. **启用 HTTPS**
3. **配置 CORS 白名单**
4. **实现请求限流**
5. **添加日志记录**
6. **定期备份数据库**

## 📦 部署建议

### Docker 部署（推荐）

创建 `docker-compose.yml`:

```yaml
version: '3'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: ad_bi_system
  
  backend:
    build: ./backend-go
    ports:
      - "8080:8080"
    depends_on:
      - mysql
  
  frontend:
    build: ./frontend-vue
    ports:
      - "3000:3000"
```

### 生产环境部署

1. **前端**: 构建静态文件，使用 Nginx 托管
2. **后端**: 编译为二进制文件，使用 systemd 或 Supervisor 管理
3. **数据库**: 配置主从复制，定期备份
4. **监控**: 集成 Prometheus + Grafana

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 📧 联系方式

如有问题或建议，请通过 GitHub Issues 联系。
