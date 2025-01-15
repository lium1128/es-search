package cfg

type Common struct {
	Mode           string `json:"mode"`
	Port           int64  `json:"port"`
	Host           string `json:"host"`
	CloseCaptcha   string `json:"close_captcha"`
	SysName        string `json:"app_name"`     // 产品名称
	Copyright      string `json:"copyright"`    // 产品名称
	Version        string `json:"version"`      // 版本
	StoragePath    string `json:"root_storage"` // 存储路径, storage目录的路径即可，比如：./storage
	CloseAuthToken string `json:"close_auth_token"`
}

// LoadCommon 加载Common配置
func LoadCommon() Common {
	return GetInstance().Common
}
