package conf

type Conf struct {
	Name     string `yaml:"name"`
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Interval int    `yaml:"interval"`
}

type Protocol struct {
	TCPConf   Conf `yaml:"tcp"`
	HTTPConf  Conf `yaml:"http"`
	HTTPSConf Conf `yaml:"https"`
	GRPCConf  Conf `yaml:"grpc"`
}
