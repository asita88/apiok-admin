package cores

import "apiok-admin/routers"

func InitRouter(conf *ConfigGlobal) error {
	routers.RouterRegister(conf.Runtime.Gin)
	return nil
}
