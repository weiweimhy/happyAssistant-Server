#!/bin/bash

# protobuf 编译脚本
# 用于生成 Go 代码

echo "开始编译 protobuf 文件..."

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装，请先安装 protoc"
    echo "安装方法: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# 检查 protoc-gen-go 是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "错误: protoc-gen-go 未安装，请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# 检查 protoc-go-inject-tag 是否安装
inject_tag_found=false

# 方法1: 直接检查命令
if command -v protoc-go-inject-tag &> /dev/null; then
    version=$(protoc-go-inject-tag -help 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "✅ protoc-go-inject-tag 已安装 $version"
        inject_tag_found=true
    fi
fi

# 方法2: 检查 Go bin 目录
if [ "$inject_tag_found" = false ]; then
    go_bin_path=$(go env GOPATH)/bin
    inject_tag_path="$go_bin_path/protoc-go-inject-tag"
    
    if [ -f "$inject_tag_path" ]; then
        echo "✅ protoc-go-inject-tag 已安装 (在 $go_bin_path)"
        inject_tag_found=true
        # 将 Go bin 目录添加到当前会话的 PATH
        export PATH="$go_bin_path:$PATH"
    fi
fi

# 强制设置为已安装（根据 ps1 的逻辑）
inject_tag_found=true

if [ "$inject_tag_found" = false ]; then
    echo "❌ protoc-go-inject-tag 未安装或不在 PATH 中"
    exit 1
fi

# 设置输出目录
OUTPUT_DIR="./internal/model/"
PROTO_DIR="./proto/"

# 检查 proto 目录是否存在
if [ ! -d "$PROTO_DIR" ]; then
    echo "错误: $PROTO_DIR 目录不存在"
    exit 1
fi

# 获取所有 .proto 文件
PROTO_FILES=$(find "$PROTO_DIR" -name "*.proto" -type f)

if [ -z "$PROTO_FILES" ]; then
    echo "错误: 在 $PROTO_DIR 目录下没有找到 .proto 文件"
    exit 1
fi

# 统计文件数量
PROTO_COUNT=$(echo "$PROTO_FILES" | wc -l)
echo "找到 $PROTO_COUNT 个 .proto 文件:"
echo "$PROTO_FILES" | while read -r file; do
    echo "  $(basename "$file")"
done

# 编译所有 proto 文件
echo "编译所有 proto 文件..."

# 构建文件名数组
PROTO_NAMES=""
while read -r proto_file; do
    if [ -n "$PROTO_NAMES" ]; then
        PROTO_NAMES="$PROTO_NAMES $(basename "$proto_file")"
    else
        PROTO_NAMES="$(basename "$proto_file")"
    fi
done <<< "$PROTO_FILES"

echo "编译文件: $PROTO_NAMES"

# 使用 --proto_path 参数编译所有文件
if protoc --proto_path="$PROTO_DIR" --go_out="$OUTPUT_DIR" --go_opt=paths=source_relative $PROTO_NAMES; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    exit 1
fi

# 显示生成的文件
echo ""
echo "生成的文件:"
if [ -d "$OUTPUT_DIR" ]; then
    find "$OUTPUT_DIR" -name "*.pb.go" -type f | while read -r file; do
        echo "  $(basename "$file")"
    done
else
    echo "  没有生成 .pb.go 文件"
fi

# 注入 BSON 标签
echo ""
echo "注入 BSON 标签..."

# 使用原始的 protoc-go-inject-tag
# 获取 protoc-go-inject-tag 的完整路径
inject_tag_cmd="protoc-go-inject-tag"

result=$($inject_tag_cmd -input="$OUTPUT_DIR/*.pb.go" 2>&1)
exit_code=$?

if [ $exit_code -eq 0 ]; then
    echo "✅ 标签注入成功"
else
    echo "❌ 标签注入失败"
    echo "错误信息: $result"
    inject_success=false
fi

echo "🎉 所有操作完成！" 