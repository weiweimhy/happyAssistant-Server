# protoc 工具安装脚本 (PowerShell)
# 自动安装 protoc 编译和标签注入工具

Write-Host "开始安装 protoc 相关工具..." -ForegroundColor Green

# 检查 Go 是否安装
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "go not found"
    }
    Write-Host "✅ Go 已安装: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ 错误: Go 未安装" -ForegroundColor Red
    Write-Host "请先安装 Go: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# 安装 protoc-gen-go
Write-Host "`n安装 protoc-gen-go..." -ForegroundColor Cyan
$result = go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ protoc-gen-go 安装成功" -ForegroundColor Green
} else {
    Write-Host "❌ protoc-gen-go 安装失败" -ForegroundColor Red
    Write-Host "错误信息: $result" -ForegroundColor Red
}

# 安装 protoc-go-inject-tag
Write-Host "`n安装 protoc-go-inject-tag..." -ForegroundColor Cyan
$result = go install github.com/favadi/protoc-gen-go-inject-tag@latest 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ protoc-go-inject-tag 安装成功" -ForegroundColor Green
} else {
    Write-Host "❌ protoc-go-inject-tag 安装失败" -ForegroundColor Red
    Write-Host "错误信息: $result" -ForegroundColor Red
}

# 检查安装结果
Write-Host "`n检查安装结果..." -ForegroundColor Cyan

$tools = @("protoc-gen-go", "protoc-go-inject-tag")
$allInstalled = $true

foreach ($tool in $tools) {
    try {
        $version = & $tool --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ $tool 已安装: $version" -ForegroundColor Green
        } else {
            Write-Host "❌ $tool 安装失败" -ForegroundColor Red
            $allInstalled = $false
        }
    } catch {
        Write-Host "❌ $tool 未找到" -ForegroundColor Red
        $allInstalled = $false
    }
}

if ($allInstalled) {
    Write-Host "`n🎉 所有工具安装成功！" -ForegroundColor Green
    Write-Host "现在可以运行 protocol_build.ps1 来编译 protobuf 文件" -ForegroundColor Cyan
} else {
    Write-Host "`n⚠️ 部分工具安装失败，请检查错误信息" -ForegroundColor Yellow
    exit 1
} 