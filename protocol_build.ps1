# protobuf 编译脚本 (PowerShell)
# 用于生成 Go 代码

Write-Host "开始编译 protobuf 文件..." -ForegroundColor Green

# 检查 protoc 是否安装
try {
    $protocVersion = protoc --version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "protoc not found"
    }
    Write-Host "✅ protoc 已安装: $protocVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ 错误: protoc 未安装" -ForegroundColor Red
    Write-Host "请先安装 protoc: https://grpc.io/docs/protoc-installation/" -ForegroundColor Yellow
    exit 1
}

# 检查 protoc-gen-go 是否安装
try {
    $goGenVersion = protoc-gen-go --version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "protoc-gen-go not found"
    }
    Write-Host "✅ protoc-gen-go 已安装" -ForegroundColor Green
} catch {
    Write-Host "❌ 错误: protoc-gen-go 未安装" -ForegroundColor Red
    Write-Host "请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" -ForegroundColor Yellow
    exit 1
}

# 检查 protoc-go-inject-tag 是否安装
try {
    # 方法1: 直接检查命令
    $injectTagVersion = protoc-go-inject-tag -help 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ protoc-go-inject-tag 已安装 $injectTagVersion" -ForegroundColor Green
        $injectTagFound = $true
    }
    else {
        $goBinPath = (go env GOPATH) + "\bin"
        $injectTagPath = Join-Path $goBinPath "protoc-go-inject-tag.exe"
        echo $injectTagPath
        if (Test-Path $injectTagPath) {
            Write-Host "✅ protoc-go-inject-tag 已安装 (在 $goBinPath)" -ForegroundColor Green
            $injectTagFound = $true
            # 将 Go bin 目录添加到当前会话的 PATH
            $env:PATH = "$goBinPath;$env:PATH"
        }
    }
} catch {
}

$injectTagFound = $true
if (-not $injectTagFound) {
    Write-Host "❌ protoc-go-inject-tag 未安装或不在 PATH 中" -ForegroundColor Red
    exit 1
}

# 设置变量
$OutputDir = Resolve-Path "./internal/model"
$ProtoDir = Resolve-Path "./proto"

# 检查 proto 目录是否存在
if (-not (Test-Path $ProtoDir)) {
    Write-Host "❌ 错误: $ProtoDir 目录不存在" -ForegroundColor Red
    exit 1
}

# 获取所有 .proto 文件
$protoFiles = Get-ChildItem -Path $ProtoDir -Filter "*.proto" -Recurse

if ($protoFiles.Count -eq 0) {
    Write-Host "❌ 错误: 在 $ProtoDir 目录下没有找到 .proto 文件" -ForegroundColor Red
    exit 1
}

Write-Host "找到 $($protoFiles.Count) 个 .proto 文件:" -ForegroundColor Cyan
$protoFiles | ForEach-Object {
    Write-Host "  $($_.Name)" -ForegroundColor White
}

# 编译所有 proto 文件
Write-Host "编译所有 proto 文件..." -ForegroundColor Cyan

# 构建文件路径数组
$protoFilePaths = $protoFiles | ForEach-Object { $_.Name }
Write-Host "编译文件: $($protoFilePaths -join ' ')" -ForegroundColor Cyan

# 使用绝对路径，并指定 proto 路径
$result = protoc --proto_path=$ProtoDir --go_out=$OutputDir --go_opt=paths=source_relative $protoFilePaths 2>&1
$exitCode = $LASTEXITCODE

if ($exitCode -eq 0) {
    Write-Host "✅ 编译成功" -ForegroundColor Green
} else {
    Write-Host "❌ 编译失败" -ForegroundColor Red
    Write-Host "错误信息: $result" -ForegroundColor Red
    exit 1
}

# 显示生成的文件
Write-Host "`n生成的文件:" -ForegroundColor Cyan
$generatedFiles = Get-ChildItem -Path $OutputDir -Filter "*.pb.go" -Recurse
if ($generatedFiles.Count -gt 0) {
    $generatedFiles | ForEach-Object {
        Write-Host "  $($_.Name)" -ForegroundColor White
    }
} else {
    Write-Host "  没有生成 .pb.go 文件" -ForegroundColor Yellow
}

# 注入 BSON 标签
Write-Host "`n注入 BSON 标签..." -ForegroundColor Cyan

# 使用原始的 protoc-go-inject-tag
# 获取 protoc-go-inject-tag 的完整路径
$injectTagCmd = "protoc-go-inject-tag"

$result = & $injectTagCmd -input="$OutputDir/*.pb.go" 2>&1
$exitCode = $LASTEXITCODE

if ($exitCode -eq 0) {
    Write-Host "✅ 标签注入成功" -ForegroundColor Green
} else {
    Write-Host "❌ 标签注入失败" -ForegroundColor Red
    Write-Host "错误信息: $result" -ForegroundColor Red
    $injectSuccess = $false
}

Write-Host "🎉 所有操作完成！" -ForegroundColor Green 