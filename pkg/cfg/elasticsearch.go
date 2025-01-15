package cfg

type ElasticSearch struct {
	Address  string `json:"address"`
	Port     uint   `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Sniff    bool   `json:"sniff"`
}

// LoadElastic 加载Elastic配置
func LoadElastic() ElasticSearch {
	return GetInstance().ElasticSearch
}
