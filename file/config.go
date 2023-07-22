package file

type WriteConfig struct {
	Level     string `mapstructure:"level"`
	FileName  string `mapstructure:"filename"`
	Driver    string `mapstructure:"driver"`
	MaxFiles  int    `mapstructure:"maxFiles"`
	MaxAge    int    `mapstructure:"maxAge"`
	MaxSize   int    `mapstructure:"maxSize"`
	Compress  bool   `mapstructure:"compress"`
	OutFormat string `mapstructure:"outFormat"`
}
