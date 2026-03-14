# 本地开发环境搭建完成指南

✅ **恭喜！你的本地开发环境已经配置完成！**

---

## 📊 当前环境状态

根据检查结果：

| 工具 | 状态 | 说明 |
|------|------|------|
| **Node.js** | ✅ 已安装 | v22.14.0 |
| **npm** | ✅ 已安装 | 随 Node.js 安装 |
| **Docker** | ❌ 未安装 | 需要安装 Docker Desktop |
| **Go** | ❌ 未安装 | 可选，本地开发需要 |
| **MySQL** | ❌ 未安装 | 可选，Docker 会自动安装 |

---

## 🚀 推荐方案：使用 Docker（最简单）

### 第一步：安装 Docker Desktop

1. **访问官网下载：**
   ```
   https://www.docker.com/products/docker-desktop/
   ```

2. **选择版本：**
   - Windows 家庭版：Docker Desktop with WSL 2
   - Windows 专业版/企业版：Docker Desktop with Hyper-V

3. **安装步骤：**
   - 下载 `Docker Desktop Installer.exe`
   - 双击运行安装程序
   - 按照向导完成安装
   - 重启电脑（如果需要）
   - 启动 Docker Desktop
   - 等待底部状态栏显示 "Docker Desktop is running"

4. **验证安装：**
   ```powershell
   docker --version
   docker compose version
   ```

### 第二步：一键启动所有服务

安装好 Docker 后，在项目根目录执行：

```powershell
.\setup-and-start.ps1
```

这个脚本会自动：
- ✅ 检查 Docker 环境
- ✅ 检查端口占用（3306, 8080, 3001, 3000）
- ✅ 创建环境配置文件
- ✅ 启动所有 4 个服务（MySQL, Go, Node, Vue）

### 第三步：访问系统

启动成功后，你会看到：

```
====================================
  所有服务已成功启动！
====================================

服务访问地址：
  🌐 前端界面：http://localhost:3000
  🔌 中间层：http://localhost:3001
  ⚙️  Go 后端：http://localhost:8080
  🗄️  MySQL: localhost:3306

数据库连接信息：
  主机：localhost
  端口：3306
  数据库：ad_bi_system
  用户名：root
  密码：root123

默认登录账号：
  用户名：admin
  密码：admin123
```

**现在可以访问 http://localhost:3000 登录系统了！**

---

## 🛠️ 备选方案：本地安装 Go 和 MySQL

如果你不想使用 Docker，可以选择本地安装所有组件。

### 第一步：运行自动安装脚本

```powershell
.\install-go-mysql.ps1
```

这个脚本会：
- 自动下载 Go 1.21.5
- 自动安装 Go
- 下载 MySQL 安装程序
- 启动 MySQL 安装向导

### 第二步：手动完成 MySQL 安装

MySQL 安装向导启动后：
1. 选择 "Server only" 或 "Developer Default"
2. 设置 root 密码为：`root123`
3. 完成安装

### 第三步：初始化数据库

```powershell
mysql -u root -proot123 < database\init.sql
```

### 第四步：启动开发服务

```powershell
.\start-dev.ps1
```

这个脚本会打开 3 个新的 PowerShell 窗口，分别启动：
- Go 后端（端口 8080）
- Node 中间层（端口 3001）
- Vue 前端（端口 3000）

---

## 📁 已创建的配置文件

### Docker 相关
- `docker-compose.yml` - Docker 编排配置
- `backend-go/Dockerfile` - Go 后端镜像
- `middleware-node/Dockerfile` - Node 中间层镜像
- `frontend-vue/Dockerfile` - Vue 前端镜像

### PowerShell 脚本
- `setup-and-start.ps1` - Docker 一键启动脚本 ⭐推荐
- `install-go-mysql.ps1` - 安装 Go 和 MySQL
- `start-dev.ps1` - 本地开发环境启动
- `check-env.ps1` - 环境检查工具

### 文档
- `QUICKSTART.md` - 快速入门指南
- `SETUP_COMPLETE.md` - 本文件

---

## ❓ 常见问题

### Q: Docker Desktop 安装失败？

**解决方案：**
1. 确保 Windows 已启用 WSL 2 或 Hyper-V
2. 以管理员身份运行安装程序
3. 重启电脑后再试
4. 查看 Docker 日志：`C:\ProgramData\Docker\log`

### Q: 端口被占用？

**解决方案：**
```powershell
# 查看占用端口的程序
netstat -ano | findstr :3306
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# 或者修改 docker-compose.yml 中的端口映射
```

### Q: 如何停止服务？

**Docker 方式：**
```powershell
docker compose down
```

**本地方式：**
关闭所有打开的 PowerShell 窗口

### Q: 如何查看日志？

**Docker 方式：**
```powershell
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f backend
docker compose logs -f frontend
```

---

## 📚 学习路径建议

1. **先运行起来**
   - 使用 Docker 方式快速启动
   - 访问前端界面体验功能

2. **了解架构**
   - 阅读 README.md 了解技术栈
   - 查看项目结构和文件说明

3. **学习代码**
   - Go 后端：`backend-go/` 目录
   - Node 中间层：`middleware-node/` 目录
   - Vue 前端：`frontend-vue/` 目录

4. **开始开发**
   - 修改前端代码（实时热更新）
   - 修改后端代码（需重启服务）
   - 查看数据库结构和测试数据

---

## 🎯 下一步行动

### 立即开始（推荐 Docker）

1. **安装 Docker Desktop**
   ```
   下载地址：https://www.docker.com/products/docker-desktop/
   ```

2. **运行启动脚本**
   ```powershell
   .\setup-and-start.ps1
   ```

3. **访问系统**
   ```
   http://localhost:3000
   账号：admin / admin123
   ```

### 需要帮助？

- 查看 `QUICKSTART.md` 详细指南
- 查看 `README.md` 项目说明
- 查看各子目录的 README.md

---

## 📞 技术支持

如有问题，请检查：
1. ✅ Docker Desktop 是否运行
2. ✅ 端口是否被占用
3. ✅ 查看日志文件

**祝你使用愉快！** 🎉

---

*最后更新：2024-03-14*
