package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	// 设置目标 URL
	targetURL := "https://example.com/"

	// 设置请求参数
	jsonData := strings.NewReader(`{"body": "body"}`)

	// 设置并发数
	concurrency := 10

	// 设置总请求时长
	duration := 600 * time.Second

	// 创建等待组
	var wg sync.WaitGroup

	// 创建限流器
	throttle := time.Tick(time.Second / time.Duration(concurrency))

	// 计数器
	var counter int64

	// 开始时间
	start := time.Now()

	// 循环发送请求
	for time.Since(start) < duration {
		<-throttle
		wg.Add(1)
		go func() {
			defer wg.Done()

			// resp, err := http.Post(targetURL, "application/json", jsonData)
			req, err := http.NewRequest("POST", targetURL, jsonData)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 添加请求头
			req.Header.Add("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer res.Body.Close()
			counter++
		}()
	}

	// 等待所有请求完成
	wg.Wait()

	// 计算 QPS
	elapsed := time.Since(start)
	qps := float64(counter) / elapsed.Seconds()
	fmt.Printf("QPS: %.2f\n", qps)
}
