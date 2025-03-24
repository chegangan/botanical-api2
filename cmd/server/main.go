package main

func main() {
	server := NewServer()
	if err := server.Run(); err != nil {
		panic(err)
	}
}

/*
打包命令
# 4. 编译Windows可执行文件
GOOS=windows GOARCH=amd64 go build -o botanical-api.exe cmd/server/main.go cmd/server/server.go

# 5. 创建分发目录(如果不存在)
mkdir -p dist

# 6. 复制所需文件到分发目录
cp botanical-api.exe dist/
cp -r configs dist/
cp -r docs dist/

# 7. 创建启动脚本
cat > dist/start.bat << EOF
@echo off
echo 正在启动植物管理系统API...
start http://localhost:8000/swagger/index.html
botanical-api.exe
pause
EOF

# 8. 打包分发
cd dist
zip -r ../botanical-api-windows.zip *
*/
