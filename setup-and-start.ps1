# 广告 BI 系统 - Windows 环境一键启动脚本
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "  广告 BI 系统 - 环境检查与启动" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Docker
Write-Host "[1/5] 检查 Docker 环境..." -ForegroundColor Yellow
try {
    $dockerVersion = docker --version
    Write-Host "  ✓ Docker 已安装：$dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Docker 未安装！" -ForegroundColor Red
    Write-Host ""
    Write-Host "请先安装 Docker Desktop:" -ForegroundColor Yellow
    Write-Host "  下载地址：https://www.docker.com/products/docker-desktop" -ForegroundColor Cyan
    Write-Host ""
    exit 1
}

# 检查 Docker Compose
Write-Host "[2/5] 检查 Docker Compose..." -ForegroundColor Yellow
try {
    $composeVersion = docker compose version
    Write-Host "  ✓ Docker Compose 已安装：$composeVersion" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Docker Compose 未安装！" -ForegroundColor Red
    exit 1
}

# 检查端口占用
Write-Host "[3/5] 检查端口占用情况..." -ForegroundColor Yellow
$ports = @(3306, 3001, 8080, 3000)
$portAvailable = $true

foreach ($port in $ports) {
    $connection = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    if ($connection) {
        Write-Host "  ✗ 端口 $port 被占用" -ForegroundColor Red
        $portAvailable = $false
    } else {
        Write-Host "  ✓ 端口 $port 可用" -ForegroundColor Green
    }
}

if (-not $portAvailable) {
    Write-Host ""
    Write-Host "请关闭占用端口的程序后重试" -ForegroundColor Yellow
    exit 1
}

# 创建 .env 文件
Write-Host "[4/5] 创建环境配置文件..." -ForegroundColor Yellow
if (!(Test-Path "middleware-node\.env")) {
    @"
PORT=3001
GO_BACKEND_URL=http://backend:8080
NODE_ENV=development
"@ | Out-File -FilePath "middleware-node\.env" -Encoding utf8
    Write-Host "  ✓ 已创建 middleware-node\.env" -ForegroundColor Green
} else {
    Write-Host "  ✓ middleware-node\.env 已存在" -ForegroundColor Green
}

# 启动 Docker Compose
Write-Host "[5/5] 启动所有服务..." -ForegroundColor Yellow
Write-Host ""
docker compose up -d --build

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "====================================" -ForegroundColor Green
    Write-Host "  所有服务已成功启动！" -ForegroundColor Green
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
    Write-Host "  密码：root123" -ForegroundColor White
    Write-Host ""
    Write-Host "默认登录账号：" -ForegroundColor Cyan
    Write-Host "  用户名：admin" -ForegroundColor White
    Write-Host "  密码：admin123" -ForegroundColor White
    Write-Host ""
    Write-Host "查看日志：" -ForegroundColor Yellow
    Write-Host "  docker compose logs -f" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "停止服务：" -ForegroundColor Yellow
    Write-Host "  docker compose down" -ForegroundColor Cyan
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "====================================" -ForegroundColor Red
    Write-Host "  服务启动失败！" -ForegroundColor Red
    Write-Host "====================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "请检查：" -ForegroundColor Yellow
    Write-Host "  1. Docker Desktop 是否已启动" -ForegroundColor White
    Write-Host "  2. 端口是否被占用" -ForegroundColor White
    Write-Host "  3. 查看错误日志：" -ForegroundColor White
    Write-Host "     docker compose logs" -ForegroundColor Cyan
    Write-Host ""
}
