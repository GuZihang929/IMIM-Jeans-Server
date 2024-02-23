package confDetail

type Zap struct {
	Filename   string `json:"filename" yaml:"filename"`      //文件名
	MaxSize    int    `json:"max_size" yaml:"maxSize"`       //每个日志文件的最大大小，单位MB
	MaxBackups int    `json:"max_backups" yaml:"maxBackups"` // 最多保留的日志文件数
	MaxAge     int    `json:"max_age" yaml:"maxAge"`         // 保留的最大天数
	Compress   bool   `json:"compress" yaml:"compress"`      //是否进行压缩
}
