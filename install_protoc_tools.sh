#!/bin/bash

# protoc 工具安装脚本
# 自动安装 protoc 编译和标签注入工具

echo "开始安装 protoc 相关工具..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go 未安装，请先安装 Go"
    echo "安装方法: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go 已安装: $(go version)"

# 安装 protoc-gen-go
echo ""
echo "安装 protoc-gen-go..."
if go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; then
    echo "✅ protoc-gen-go 安装成功"
else
    echo "❌ protoc-gen-go 安装失败"
fi

# 安装 protoc-go-inject-tag
echo ""
echo "安装 protoc-go-inject-tag..."
if go install github.com/favadi/protoc-gen-go-inject-tag@latest; then
    echo "✅ protoc-go-inject-tag 安装成功"
else
    echo "❌ protoc-go-inject-tag 安装失败"
fi

# 检查安装结果
echo ""
echo "检查安装结果..."

tools=("protoc-gen-go" "protoc-go-inject-tag")
all_installed=true

for tool in "${tools[@]}"; do
    if command -v "$tool" &> /dev/null; then
        version=$($tool --version 2>/dev/null || echo "版本未知")
        echo "✅ $tool 已安装: $version"
    else
        echo "❌ $tool 未找到"
        all_installed=false
    fi
done

if [ "$all_installed" = true ]; then
    echo ""
    echo "🎉 所有工具安装成功！"
    echo "现在可以运行 protocol_build.sh 来编译 protobuf 文件"
else
    echo ""
    echo "⚠️ 部分工具安装失败，请检查错误信息"
    exit 1
fi 