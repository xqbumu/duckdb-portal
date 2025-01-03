package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const version = "unknown"

// 定义一个简单的处理函数
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, DuckDB Portal!")
}

func main() {
	// 配置重载方案
	method := flag.String("method", "auto", "config reload method: auto/signal")
	// 定义一个端口参数，默认为8080
	port := flag.String("port", "8080", "HTTP server port")
	// 添加config参数
	configPath := flag.String("config", "config.yaml", "Path to config file")
	// 添加version参数
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	// 检查是否需要显示版本信息
	if *showVersion {
		fmt.Printf("DuckDB Portal version: %s\n", version)
		return
	}

	// 加载配置文件
	if err := LoadConfig(*configPath, *method); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("Failed to load config: %v\n", err)
		return
	}

	// 创建一个新的Chi路由器
	r := chi.NewRouter()

	// 使用中间件
	r.Use(middleware.Logger) // 可选：添加日志中间件
	r.Use(authMiddleware)    // 添加授权中间件

	// 注册路由和处理函数
	r.Get("/hello", helloHandler)

	r.Post("/duckdb", handlerDuckDB)

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("Starting server at %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Println(err)
	}
}
