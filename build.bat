@echo off
chcp 65001

echo 切换到linux环境...
go env -w GOOS=linux

echo 编译守护进程linux版本...
go build -ldflags "-s -w" .
echo 压缩二进制文件...
upx -9 daemon

cd testapp
echo 编译测试程序...
go build -ldflags "-s -w" .
echo 压缩二进制文件...
upx -9 testapp
cd ..

echo 切换回windows环境...
go env -w GOOS=windows
echo 完成！
echo.
