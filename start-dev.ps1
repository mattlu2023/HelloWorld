# 本地开发环境启动脚本（不使用 Docker）
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "  启动本地开发环境" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Go
Write-Host "[1/5] 检查 Go 环境..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "  ✓ $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Go 未安装！请先运行：.\install-go-mysql.ps1" -ForegroundColor Red
    exit 1
}

# 检查 MySQL
Write-Host "[2/5] 检查 MySQL 环境..." -ForegroundColor Yellow
try {
    $mysqlVersion = mysql --version
    Write-Host "  ✓ $mysqlVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ MySQL 未安装！请先运行：.\install-go-mysql.ps1" -ForegroundColor Red
    exit 1
}

# 检查 Node.js
Write-Host "[3/5] 检查 Node.js 环境..." -ForegroundColor Yellow
try {
    $nodeVersion = node --version
    Write-Host "  ✓ Node.js $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Node.js 未安装！" -ForegroundColor Red
    Write-Host "  请安装 Node.js: https://nodejs.org/" -ForegroundColor Yellow
    exit 1
}

# 初始化数据库
Write-Host "[4/5] 初始化 MySQL 数据库..." -ForegroundColor Yellow
try {
    Write-Host "  创建数据库并导入数据..." -ForegroundColor Cyan
    mysql -u root -proot123 < database\init.sql 2>$null
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ 数据库初始化成功" -ForegroundColor Green
    } else {
        # 尝试不带密码连接
        Write-Host "  尝试无密码连接..." -ForegroundColor Yellow
        mysql -u root < database\init.sql 2>$null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ 数据库初始化成功（无密码）" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ 数据库初始化失败，请手动执行：" -ForegroundColor Yellow
            Write-Host "    mysql -u root -p < database\init.sql" -ForegroundColor Cyan
        }
    }
} catch {
    Write-Host "  ⚠ 数据库初始化异常：$_" -ForegroundColor Yellow
    Write-Host "  请确保 MySQL 服务已启动" -ForegroundColor Yellow
}

# 启动服务
Write-Host "[5/5] 启动开发服务..." -ForegroundColor Yellow
Write-Host ""

# 启动 Go 后端
Write-Host "  启动 Go 后端服务..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd backend-go; go run main.go"
Start-Sleep -Seconds 2

# 启动 Node 中间层
Write-Host "  启动 Node.js 中间层..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd middleware-node; npm start"
Start-Sleep -Seconds 2

# 启动 Vue 前端
Write-Host "  启动 Vue 前端服务..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd frontend-vue; npm run dev"

Write-Host ""
Write-Host "====================================" -ForegroundColor Green
Write-Host "  所有服务已启动！" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Green
Write-Host ""
Write-Host "服务访问地址：" -ForegroundColor Cyan
Write-Host "  🌐 前端界面：http://localhost:3000" -ForegroundColor White
Write-Host "  🔌 中间层：http://localhost:3001" -ForegroundColor White
Write-Host "  ⚙️  Go 后端：http://localhost:8080" -ForegroundColor White
Write-Host "  🗄️  MySQL: localhost:3306" -ForegroundColor White
Write-Host ""
Write-Host "数据库连接信息：" -ForegroundColor Cyan
Write-Host "  主机：localhost" -ForegroundColor White
Write-Host "  端口：3306" -ForegroundColor White
Write-Host "  数据库：ad_bi_system" -ForegroundColor White
Write-Host "  用户名：root" -ForegroundColor White
Write-Host "  密码：root123（或无密码）" -ForegroundColor White
Write-Host ""
Write-Host "默认登录账号：" -ForegroundColor Cyan
Write-Host "  用户名：admin" -ForegroundColor White
Write-Host "  密码：admin123" -ForegroundColor White
Write-Host ""
Write-Host "按任意键关闭此窗口，服务将在独立窗口中运行" -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
