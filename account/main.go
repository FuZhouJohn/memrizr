package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FuZhouJohn/memrizr/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("服务正在启动...")

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R: router,
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务初始化错误：%v\n", err)
		}
	}()

	log.Printf("正在监听端口 %v ", srv.Addr)

	// 等待退出信号
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("正在停止服务...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务被迫关闭：%v\n", err)
	}
}
