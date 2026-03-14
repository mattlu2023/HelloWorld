# Go 和 MySQL 本地环境安装脚本
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "  安装 Go 和 MySQL 开发环境" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# 检测系统架构
$architecture = (Get-CimInstance Win32_Processor).AddressWidth
$is64Bit = $architecture -eq 64

Write-Host "[1/4] 系统架构：$("$(if ($is64Bit) {'64 位'} else {'32 位'})")" -ForegroundColor Green

# 安装 Go
Write-Host ""
Write-Host "[2/4] 安装 Go 语言环境..." -ForegroundColor Yellow

$goUrl = if ($is64Bit) {
    "https://go.dev/dl/go1.21.5.windows-amd64.msi"
} else {
    "https://go.dev/dl/go1.21.5.windows-386.msi"
}

$goInstaller = "$env:TEMP\go-installer.msi"

try {
    Write-Host "  下载 Go 安装程序..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller -UseBasicParsing
    
    Write-Host "  安装 Go..." -ForegroundColor Cyan
    Start-Process msiexec.exe -Wait -ArgumentList "/i `"$goInstaller`" /quiet /norestart"
    
    # 刷新环境变量
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    Write-Host "  ✓ Go 安装完成" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Go 安装失败：$_" -ForegroundColor Red
    Write-Host "  请手动下载：https://go.dev/dl/" -ForegroundColor Yellow
}

# 安装 MySQL
Write-Host ""
Write-Host "[3/4] 安装 MySQL..." -ForegroundColor Yellow

$mysqlInstallerUrl = "https://dev.mysql.com/get/Downloads/MySQLInstaller/mysql-installer-community-8.0.35.0.msi"
$mysqlInstaller = "$env:TEMP\mysql-installer.msi"

try {
    Write-Host "  下载 MySQL 安装程序..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $mysqlInstallerUrl -OutFile $mysqlInstaller -UseBasicParsing
    
    Write-Host "  启动 MySQL 安装向导..." -ForegroundColor Cyan
    Write-Host "  提示：请选择 'Server only' 或 'Developer Default'" -ForegroundColor Yellow
    Start-Process msiexec.exe -Wait -ArgumentList "/i `"$mysqlInstaller`""
    
    Write-Host "  ✓ MySQL 安装程序已启动" -ForegroundColor Green
} catch {
    Write-Host "  ✗ MySQL 下载失败：$_" -ForegroundColor Red
    Write-Host "  请手动下载：https://dev.mysql.com/downloads/installer/" -ForegroundColor Yellow
}

# 清理安装文件
Write-Host ""
Write-Host "[4/4] 清理安装文件..." -ForegroundColor Yellow
Remove-Item $goInstaller -ErrorAction SilentlyContinue
Write-Host "  ✓ 清理完成" -ForegroundColor Green

Write-Host ""
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "  安装完成！" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "请重启终端后验证安装：" -ForegroundColor Yellow
Write-Host "  go version" -ForegroundColor Cyan
Write-Host "  mysql --version" -ForegroundColor Cyan
Write-Host ""
Write-Host "然后运行启动脚本：" -ForegroundColor Yellow
Write-Host "  .\start-dev.ps1" -ForegroundColor Cyan
Write-Host ""
