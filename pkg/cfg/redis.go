package cfg

type Redis struct {
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

// LoadRedis 加载consul配置
func LoadRedis() Redis {
	return GetInstance().Redis
}
