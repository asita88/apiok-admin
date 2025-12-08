package cores

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/services"
)

// InitLetsEncrypt 初始化Let's Encrypt服务
func InitLetsEncrypt(conf *ConfigGlobal) error {
	if !conf.LetsEncrypt.Enabled {
		return nil
	}

	config := &services.LetsEncryptConfig{
		Enabled:         conf.LetsEncrypt.Enabled,
		Email:           conf.LetsEncrypt.Email,
		UseStaging:      conf.LetsEncrypt.UseStaging,
		CertDir:         conf.LetsEncrypt.CertDir,
		RenewBeforeDays: conf.LetsEncrypt.RenewBeforeDays,
	}

	err := services.NewLetsEncryptService().InitLetsEncrypt(config)
	if err != nil {
		packages.Log.Errorf("Failed to initialize Let's Encrypt: %v", err)
		return err
	}

	packages.Log.Info("Let's Encrypt initialized successfully")
	return nil
}
