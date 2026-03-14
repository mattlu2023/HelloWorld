# 广告 BI 系统 - 快速入门指南

本指南将帮助你在本地搭建完整的开发和运行环境。

---

## 📋 目录

1. [环境要求](#环境要求)
2. [方式一：使用 Docker（推荐）](#方式一使用-docker推荐)
3. [方式二：本地安装 Go 和 MySQL](#方式二本地安装-go-和-mysql)
4. [常见问题](#常见问题)
5. [下一步](#下一步)

---

## 🛠️ 环境要求

### 必需
- ✅ **Node.js 18+** （运行前端和中间层）
- ✅ **npm** （Node.js 包管理器）

### 可选（推荐）
- 🐳 **Docker Desktop** （一键启动所有服务）
- 🔷 **Go 1.21+** （本地开发 Go 后端）
- 🗄️ **MySQL 8.0+** （本地数据库开发）

---

## 方式一：使用 Docker（推荐 ⭐⭐⭐⭐⭐）

这是最简单的方式，无需手动安装 Go 和 MySQL。

### 步骤

#### 1. 安装 Docker Desktop

**下载地址：** https://www.docker.com/products/docker-desktop

1. 下载并安装 Docker Desktop
2. 启动 Docker Desktop
3. 等待 Docker 引擎启动完成

#### 2. 一键启动所有服务

在项目根目录执行：

```powershell
.\setup-and-start.ps1
```

脚本会自动：
- ✅ 检查 Docker 环境
- ✅ 检查端口占用
- ✅ 创建环境配置
- ✅ 启动 MySQL、Go 后端、Node 中间层、Vue 前端

#### 3. 访问服务

启动成功后，访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| 🌐 前端 | http://localhost:3000 | 用户界面 |
| 🔌 中间层 | http://localhost:3001 | API 代理 |
| ⚙️ Go 后端 | http://localhost:8080 | REST API |
| 🗄️ MySQL | localhost:3306 | 数据库 |

**默认登录账号：**
- 用户名：`admin`
- 密码：`admin123`

### Docker 常用命令

```powershell
# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 重启某个服务
docker compose restart backend

# 停止所有服务
docker compose down

# 停止并删除数据（谨慎使用）
docker compose down -v
```

---

## 方式二：本地安装 Go 和 MySQL

适合需要本地调试和开发的场景。

### 步骤

#### 1. 自动安装 Go 和 MySQL

```powershell
.\install-go-mysql.ps1
```

这个脚本会：
- 自动下载并安装 Go
- 下载并启动 MySQL 安装程序
- 配置环境变量

#### 2. 手动安装（如果自动安装失败）

**安装 Go：**
1. 访问 https://golang.org/dl/
2. 下载 Windows 安装程序
3. 运行安装程序
4. 验证：`go version`

**安装 MySQL：**
1. 访问 https://dev.mysql.com/downloads/installer/
2. 下载 MySQL Installer
3. 运行安装程序
4. 选择 "Server only" 或 "Developer Default"
5. 设置 root 密码（建议：`root123`）
6. 验证：`mysql --version`

**安装 Node.js：**
1. 访问 https://nodejs.org/
2. 下载 LTS 版本
3. 运行安装程序
4. 验证：`node --version`

#### 3. 初始化数据库

```powershell
mysql -u root -p < database\init.sql
```

输入 MySQL root 密码（默认：`root123` 或无密码）

#### 4. 启动开发环境

```powershell
.\start-dev.ps1
```

这个脚本会：
- ✅ 检查所有环境
- ✅ 初始化数据库
- ✅ 启动 Go 后端
- ✅ 启动 Node 中间层
- ✅ 启动 Vue 前端

---

## 🔍 环境检查

运行环境检查脚本：

```powershell
.\check-env.ps1
```

会显示所有已安装的工具和版本信息。

---

## ❓ 常见问题

### Q1: Docker 启动失败，提示端口被占用

**解决方案：**
1. 检查端口占用：
   ```powershell
   netstat -ano | findstr :3306
   netstat -ano | findstr :8080
   netstat -ano | findstr :3000
   ```
2. 关闭占用端口的程序
3. 或者修改 `docker-compose.yml` 中的端口映射

### Q2: MySQL 连接失败

**解决方案：**
1. 检查 MySQL 服务是否启动
2. 确认用户名密码正确
3. 尝试无密码连接：`mysql -u root`
4. 查看 MySQL 错误日志

### Q3: Go 后端编译失败

**解决方案：**
```powershell
cd backend-go
go mod download
go build -o main.exe
```

检查错误信息，确保所有依赖已下载。

### Q4: 前端页面空白或无法加载

**解决方案：**
1. 检查浏览器控制台错误
2. 确认中间层和后端已启动
3. 清除浏览器缓存
4. 检查跨域配置

### Q5: 数据库初始化失败

**手动初始化：**
```sql
mysql -u root -p
CREATE DATABASE ad_bi_system;
USE ad_bi_system;
source database/init.sql;
```

---

## 📊 项目结构

```
HelloWorld/
├── database/              # 数据库脚本
│   └── init.sql          # MySQL 初始化（含测试数据）
├── backend-go/           # Go 后端
│   ├── main.go
│   ├── handlers/         # API 处理器
│   └── middleware/       # 中间件
├── middleware-node/      # Node.js 中间层
│   └── server.js
├── frontend-vue/         # Vue 3 前端
│   └── src/
│       ├── views/        # 页面组件
│       └── api/          # API 接口
├── docker-compose.yml    # Docker 配置
├── setup-and-start.ps1   # Docker 一键启动
├── install-go-mysql.ps1  # 安装 Go 和 MySQL
├── start-dev.ps1         # 本地开发启动
└── check-env.ps1         # 环境检查
```

---

## 🎯 下一步

环境搭建完成后：

1. **访问前端界面**
   - http://localhost:3000
   - 使用 admin/admin123 登录

2. **查看 API 文档**
   - 参考 README.md 中的 API 接口说明

3. **开始开发**
   - 修改前端代码（实时热更新）
   - 修改 Go 后端代码（需重启服务）
   - 查看数据库数据

4. **学习代码结构**
   - 阅读各模块的 README.md
   - 查看示例代码和注释

---

## 📞 获取帮助

如有问题，请：
1. 查看项目 README.md
2. 检查 Docker 日志：`docker compose logs`
3. 查看各服务的错误信息

祝你使用愉快！🎉
