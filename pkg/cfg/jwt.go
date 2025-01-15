package cfg

type Jwt struct {
	Expire uint   `json:"expire"` // Token过期时间(Minute)
	Key    string `json:"key"`    // Token加密key
}

func LoadJwt() Jwt {
	return GetInstance().Jwt
}
