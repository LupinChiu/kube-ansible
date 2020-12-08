package conf

type Database struct {
	External      bool   `yaml:"external"`
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	Type          string `yaml:"type"`
	Db            string `yaml:"db"`
	HardcodeDelay int    `yaml:"hardcodedelay"`
	Conn          struct {
		Maxidle int `yaml:"maxidle"`
		Maxopen int `yaml:"maxopen"`
	}
}

type Protocol struct {
	Name   string `yaml:"name"`
	Enable bool   `yaml:"enable"`
	Port   string `yaml:"port"`
}

type ConfigMap struct {
	TCPConf   Protocol `yaml:"tcp"`
	HTTPConf  Protocol `yaml:"http"`
	HTTPSConf Protocol `yaml:"https"`
	GRPCConf  Protocol `yaml:"grpc"`
	Db        Database `yaml:"db"`
}

type Conf struct {
	Proto Protocol
	Db    Database
}
