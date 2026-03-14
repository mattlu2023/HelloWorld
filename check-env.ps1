Write-Host "Environment Check" -ForegroundColor Cyan
Write-Host "=================" -ForegroundColor Cyan
Write-Host ""

# Check Docker
Write-Host "Checking Docker..." -NoNewline
try {
    docker --version
    Write-Host " OK" -ForegroundColor Green
} catch {
    Write-Host " Not installed" -ForegroundColor Red
}

# Check Node.js
Write-Host "Checking Node.js..." -NoNewline
try {
    node --version
    Write-Host " OK" -ForegroundColor Green
} catch {
    Write-Host " Not installed" -ForegroundColor Red
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Install Docker Desktop (recommended)" -ForegroundColor White
Write-Host "  2. Run: .\setup-and-start.ps1" -ForegroundColor White
Write-Host ""
