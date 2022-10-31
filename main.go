package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang_blog_backend/dao/mysql"
	"golang_blog_backend/dao/redis"
	"golang_blog_backend/logger"
	"golang_blog_backend/routes"
	"golang_blog_backend/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Println("init settings failed, err:", err)
		return
	}
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("init logger failed, err:", err)
		return
	}
	// 把缓冲区的日志追加到我们的日志文件中
	defer zap.L().Sync()
	// 初始化Mysql
	if err := mysql.Init(); err != nil {
		fmt.Println("init mysql failed, err:", err)
		return
	}
	// 关闭mysql服务
	defer mysql.Close()
	// 初始化redis
	if err := redis.Init(); err != nil {
		fmt.Println("init redis failed, err:", err)
		return
	}
	// 关闭redis服务
	redis.Close()
	// 注册路由
	r := routes.SetupRouter()
	// 启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务器
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送syscall.SIGTERM信号
	// kill -2 发送syscall.SIGINT信号，我们常用的ctrl + c就是触发了SIGTERM信号
	// kill -9 发送 syscall.SIGKILL信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的syscall.SIGTERM或syscall.SIGINT信号抓发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号后才往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
