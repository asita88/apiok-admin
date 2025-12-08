package cores

import (
	"apiok-admin/app/services"
	"time"
)

func InitGoroutineFunc() {
	go dynamicValidationPluginData()
	go letsEncryptRenewal()
}

func dynamicValidationPluginData() {

	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

	for {
		services.PluginBasicInfoMaintain()

		<-timer.C
	}
}

func letsEncryptRenewal() {
	// 每天检查一次证书续期
	timer := time.NewTicker(24 * time.Hour)
	defer timer.Stop()

	// 启动时立即检查一次
	time.Sleep(1 * time.Minute) // 等待服务完全启动
	services.NewLetsEncryptService().RenewExpiringCertificates()

	for {
		<-timer.C
		services.NewLetsEncryptService().RenewExpiringCertificates()
	}
}
