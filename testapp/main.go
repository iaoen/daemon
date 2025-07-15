package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// 解析命令行参数
	exitAfter := flag.Duration("exit-after", 30*time.Second, "程序运行多久后退出")
	randomExit := flag.Bool("random-exit", true, "是否随机退出")
	flag.Parse()

	fmt.Printf("测试程序启动，PID: %d\n", os.Getpid())
	fmt.Printf("参数: exit-after=%s, random-exit=%v\n", *exitAfter, *randomExit)

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 计算退出时间
	startTime := time.Now()
	exitTime := startTime.Add(*exitAfter)

	// 每秒输出一次信息
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	counter := 0
	for range ticker.C {
		counter++
		fmt.Printf("测试程序运行中... 已运行%d秒\n", counter)

		// 检查是否应该退出
		now := time.Now()
		if now.After(exitTime) {
			fmt.Println("达到指定运行时间，正常退出")
			os.Exit(0)
		}

		// 有10%的概率随机退出
		if *randomExit && rand.Float32() < 0.1 {
			exitCode := rand.Intn(2) + 1 // 1或2的退出码
			fmt.Printf("随机退出，退出码: %d\n", exitCode)
			os.Exit(exitCode)
		}
	}
}
