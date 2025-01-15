package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lium-product/es-search/pkg/cfg"
	"lium-product/es-search/search/logs"
	"lium-product/es-search/search/routes"
)

func main() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = location
	cfg.InitLoadCfg()
	common := cfg.LoadCommon()
	// 程序退出前处理
	go Finally()

	// 注册路由启动服务
	r := routes.Init(common.Mode)
	logs.GetLogger().Infof("server start at %s:%d", common.Host, common.Port)
	err = r.Run(fmt.Sprintf("%s:%d", common.Host, common.Port))
	if err != nil {
		logs.GetLogger().Fatalf("run err: %v", err)
	}
}

// Finally 程序退出前的处理
func Finally() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	for {
		v := <-signals
		switch v {
		case syscall.SIGTERM:
			logs.GetLogger().Infof("Got signal %s", v.String())
			return
		case syscall.SIGINT:
			logs.GetLogger().Infof("Got signal %s", v.String())
			return
		}
	}
}
