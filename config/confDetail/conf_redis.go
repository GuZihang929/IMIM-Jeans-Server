package confDetail

type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Password string `json:"password" yaml:"password"`
}

func (r Redis) Addr() string {
	return r.Host
}
