package packages

var globalConfig interface{}

func SetConfig(conf interface{}) {
	globalConfig = conf
}

func GetConfig() interface{} {
	return globalConfig
}

type configApiOk struct {
	Protocol string
	Ip       string
	Port     int
	Domain   string
	Secret   string
}

var ConfigApiOk configApiOk

func SetConfigApiOk(protocol string, ip string, port int, domain string, secret string) {
	apiOk := configApiOk{
		Protocol: protocol,
		Ip:       ip,
		Port:     port,
		Domain:   domain,
		Secret:   secret,
	}

	ConfigApiOk = apiOk
}
