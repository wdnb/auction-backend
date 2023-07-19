package main

import (
	"auction-website/api"
	"auction-website/conf"
	"auction-website/message/im"
	"auction-website/message/nsq/consumer"
	"auction-website/scheduler"
	"auction-website/utils/logger"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	//*******************系统依赖初始化 start******************
	//1：初始化viper
	if err := conf.Viper(); err != nil {
		fmt.Printf("init viper failed,err:%v\n", err)
		return
	}
	//读取系统配置
	c := conf.Init()
	//2：初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init zap failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync() //写入硬盘
	//defer com.Mysql.Close()
	//defer com.MongoDB.Disconnect(context.Background())
	//defer com.Redis.Close()
	//defer com.Producer.Stop()
	//************************系统依赖初始化 end******************************
	//************************项目应用初始化 start******************************
	//注册消费者
	consumer.InitConsumer(c)
	//初始化im
	go im.NewIM(c).Start()
	//defer tickerScheduler.Stop()
	//************************项目应用初始化 end******************************
	//************************启动http service start******************************
	//初始化系统依赖
	//初始化redis
	//注册拍卖定时器
	scheduler.NewTickerScheduler("auction-lock", c).Start()
	//6：注册路由
	r := api.Setup(c)
	//7：启动http服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("http.port")),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
	//************************启动http service end******************************
}
