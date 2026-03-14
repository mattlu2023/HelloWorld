# 表单提交系统

一个简单实用的表单提交和查看系统，支持用户提交信息并实时查看所有提交记录。

## 🚀 功能特性

- ✅ 表单提交：姓名、性别、邮箱、留言内容
- ✅ 实时查看：展示所有已提交的表单记录
- ✅ 数据持久化：使用 SQLite 数据库存储数据
- ✅ 美观界面：渐变色设计和响应式布局
- ✅ 无需登录：开放提交和查看

## 🛠️ 技术栈

**后端:**
- Node.js
- Express.js
- SQLite3

**前端:**
- HTML5
- CSS3
- JavaScript (原生)

## 📦 安装步骤

### 1. 克隆项目
```bash
git clone <你的仓库地址>
cd HelloWorld
```

### 2. 安装依赖
```bash
npm install
```

### 3. 启动服务器
```bash
npm start
```

### 4. 访问网站
打开浏览器访问：http://localhost:3000

## 📁 项目结构

```
HelloWorld/
├── server.js              # 后端服务器
├── public/
│   └── index.html         # 前端页面
├── database.sqlite        # SQLite 数据库 (运行时自动创建)
├── package.json           # 项目配置
├── .gitignore            # Git 忽略文件
└── README.md             # 项目说明
```

## 🔌 API 接口

### 提交表单
```
POST /api/submit
Content-Type: application/json

{
  "name": "张三",
  "gender": "男",
  "email": "zhangsan@example.com",
  "message": "这是留言内容"
}
```

### 获取所有提交
```
GET /api/submissions
```

### 获取单条提交
```
GET /api/submissions/:id
```

## 📝 数据库结构

**表名:** `submissions`

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自动递增 |
| name | TEXT | 姓名 |
| gender | TEXT | 性别 |
| email | TEXT | 邮箱 |
| message | TEXT | 留言内容 |
| created_at | DATETIME | 提交时间 |

## ⚙️ 配置说明

- **端口号:** 3000 (可在 `server.js` 中修改)
- **数据库文件:** `database.sqlite` (自动创建)
- **静态文件目录:** `public/`

## 🔒 安全提示

⚠️ **注意:** 当前版本没有身份验证和权限控制，适合内部测试或学习使用。

如需在生产环境使用，建议添加:
- 用户认证系统
- 输入验证和 XSS 防护
- CSRF 保护
- 数据备份机制
- HTTPS 加密

## 📄 许可证

ISC

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📧 联系方式

如有问题或建议，请通过 GitHub Issues 联系。
