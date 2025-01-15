package cfg

type Logger struct {
	Level         string `json:"level"`
	OutPutConsole bool   `json:"output_console"`
	OutPutFile    bool   `json:"output_file"`
	FileName      string `json:"file_name"`
	MaxSize       int    `json:"max_size"`
	MaxAge        int    `json:"max_age"` //util:day
	MaxBackups    int    `json:"max_backups"`
	LocalTime     bool   `json:"local_time"`
	Compress      bool   `json:"compress"`
}

// LoadLogger 加载日志配置
func LoadLogger() Logger {
	GetInstance()
	l := GetInstance().Logger
	l.OutPutConsole = true
	return l
}
