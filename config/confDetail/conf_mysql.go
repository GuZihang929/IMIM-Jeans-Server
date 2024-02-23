package confDetail

import (
	"strconv"
	"time"
)

type Mysql struct {
	Username string        `json:"username" yaml:"username"`
	Password string        `json:"password" yaml:"password"`
	Host     string        `json:"host" yaml:"host"`
	Port     int           `json:"port" yaml:"port"`
	Dbname   string        `json:"dbname" yaml:"dbname"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
	Config   string        `json:"config" yaml:"config"`
}

func (m Mysql) Dsn() string {
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.Dbname + "?" + m.Config
	return dsn
}
