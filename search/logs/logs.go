package logs

import (
	"sync"

	"lium-product/es-search/pkg/cfg"
	"lium-product/es-search/pkg/common_logs"
)

var (
	// Logger 日志
	once sync.Once
	Log  *common_logs.Logger

	syncOnce sync.Once
	LogSync  *common_logs.Logger

	onceCrontab sync.Once
	LogCrontab  *common_logs.Logger
)

func GetLogger() *common_logs.Logger {
	if Log == nil {
		once.Do(func() {
			Log = common_logs.InitLogger(cfg.LoadLogger())
		})
	}
	return Log
}

func GetCrontabLogger() *common_logs.Logger {
	if LogCrontab == nil {
		onceCrontab.Do(func() {
			conf := cfg.LoadLogger()
			// 定时任务关闭控制台输出
			conf.OutPutConsole = false
			LogCrontab = common_logs.InitLogger(conf)
		})
	}

	return LogCrontab
}
