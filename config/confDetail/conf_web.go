package confDetail

type Web struct {
	Host string `json:"host" yaml:"Host"`
	Port string `json:"port" yaml:"Port"`
	Env  string `json:"env" yaml:"Env"`
}
