#!/bin/bash

# 定义源文件和输出目录
SOURCE_FILE="src/main/main.go"
OUTPUT_DIR="bin"

# 根据操作系统设置输出文件后缀
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OUTPUT_FILE="$OUTPUT_DIR/keep-online"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OUTPUT_FILE="$OUTPUT_DIR/keep-online"
elif [[ "$OSTYPE" == "cygwin" ]]; then
    OUTPUT_FILE="$OUTPUT_DIR/keep-online.exe"
elif [[ "$OSTYPE" == "msys" ]]; then
    OUTPUT_FILE="$OUTPUT_DIR/keep-online.exe"
elif [[ "$OSTYPE" == "win32" ]]; then
    OUTPUT_FILE="$OUTPUT_DIR/keep-online.exe"
else
    echo "Unsupported OS type: $OSTYPE"
    exit 1
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 编译 Go 文件
go build -o "$OUTPUT_FILE" "$SOURCE_FILE"

# 检查编译是否成功
if [[ $? -eq 0 ]]; then
    echo "编译成功：$OUTPUT_FILE"
else
    echo "编译失败"
    exit 1
fi
