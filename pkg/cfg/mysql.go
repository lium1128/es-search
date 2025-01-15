package cfg

const (
	LogDebugMode   = "debug"
	LogReleaseMode = "release"
	LogTestMode    = "test"
)

type MySql struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	CharSet  string `json:"charset"`
	LogLevel string `json:"log-level"` // 2345=error,01=warning,-1=debug
	SlowTime int64  `json:"slow-time"`
}

// LoadMysql 加载Mysql配置
func LoadMysql() MySql {
	return GetInstance().MySql
}
