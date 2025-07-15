package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// 配置选项
type Config struct {
	ProgramPath  string
	Args         []string
	RestartDelay time.Duration
}

func main() {
	// 解析命令行参数
	restartDelay := flag.Duration("delay", 5*time.Second, "重启程序的延迟时间")
	flag.Parse()

	// 获取要守护的程序
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("请指定要守护的可执行程序")
	}

	programPath := args[0]
	programArgs := args[1:]

	// 检查程序是否存在
	absPath, err := filepath.Abs(programPath)
	if err != nil {
		log.Fatalf("获取程序路径失败: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("程序不存在: %s", absPath)
	}

	config := &Config{
		ProgramPath:  absPath,
		Args:         programArgs,
		RestartDelay: *restartDelay,
	}

	// 启动守护进程
	daemon := NewDaemon(config)
	daemon.Start()
}

// Daemon 结构体
type Daemon struct {
	config     *Config
	cmd        *exec.Cmd
	quit       chan struct{}
	terminated bool
}

// NewDaemon 创建一个新的守护进程
func NewDaemon(config *Config) *Daemon {
	return &Daemon{
		config: config,
		quit:   make(chan struct{}),
	}
}

// Start 启动守护进程
func (d *Daemon) Start() {
	log.Printf("守护进程启动，监控程序: %s", d.config.ProgramPath)

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动被守护的程序
	go d.run()

	// 等待信号
	<-sigChan
	log.Println("收到终止信号，正在关闭...")
	d.terminate()
}

// run 运行被守护的程序并在需要时重启
func (d *Daemon) run() {
	for {
		select {
		case <-d.quit:
			return
		default:
			if d.terminated {
				return
			}

			// 创建并启动cmd
			cmd := exec.Command(d.config.ProgramPath, d.config.Args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			d.cmd = cmd
			log.Printf("启动程序: %s %v", d.config.ProgramPath, d.config.Args)

			err := cmd.Start()
			if err != nil {
				log.Printf("启动程序失败: %v", err)
				time.Sleep(d.config.RestartDelay)
				continue
			}

			// 等待进程结束
			err = cmd.Wait()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
						exitCode = status.ExitStatus()
					}
				}
				log.Printf("程序退出，状态码: %d, 错误: %v", exitCode, err)
			} else {
				log.Printf("程序正常退出，状态码: 0")
			}

			// 如果是手动终止，则不重启
			if d.terminated {
				return
			}

			log.Printf("%s秒后重启程序...", d.config.RestartDelay)
			time.Sleep(d.config.RestartDelay)
		}
	}
}

// terminate 终止守护进程和被守护的程序
func (d *Daemon) terminate() {
	d.terminated = true
	close(d.quit)

	if d.cmd != nil && d.cmd.Process != nil {
		log.Println("正在停止被守护的程序...")

		// 尝试优雅地终止程序
		d.cmd.Process.Signal(syscall.SIGTERM)

		// 给程序一些时间来清理
		time.Sleep(3 * time.Second)

		// 如果程序还在运行，强制终止
		if err := d.cmd.Process.Kill(); err != nil {
			log.Printf("强制终止程序失败: %v", err)
		}
	}

	log.Println("守护进程已终止")
}
