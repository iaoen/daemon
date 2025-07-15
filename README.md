# Go守护进程程序

这是一个用Go语言编写的守护进程程序，可以监控并在需要时自动重启指定的可执行程序。

## 功能特点

- 监控指定的可执行程序
- 在程序崩溃或退出后自动重启
- 支持自定义重启延迟时间
- 优雅地处理终止信号（Ctrl+C）

## 使用方法

### 编译

```bash
cd daemon
go build -o daemon
```

### 运行

基本用法：

```bash
./daemon 你的程序路径 [程序参数...]
```

例如：

```bash
./daemon ./myapp --port=8080
```

### 选项

- `-restart-delay`：程序崩溃后重启的延迟时间（默认为5秒）

例如，设置10秒的重启延迟：

```bash
./daemon -restart-delay=10s ./myapp
```

## 示例

假设有一个名为`web-server`的Web服务器程序：

```bash
# 在80端口启动web-server，并在崩溃后3秒重启
./daemon -restart-delay=3s ./web-server -port=80
```

## 退出

按下`Ctrl+C`可以安全地终止守护进程和被守护的程序。 