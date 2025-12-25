package routers

import (
	"apiok-admin/app/controllers/admin"
	"apiok-admin/app/middlewares"

	"github.com/gin-gonic/gin"
)

func RouterRegister(routerEngine *gin.Engine) {

	// ACME HTTP-01 challenge verification (must be at root path, no login required)
	routerEngine.GET("/.well-known/acme-challenge/:token", admin.LetsEncryptChallenge)

	noLoginRouter := routerEngine.Group("admin")
	{
		user := noLoginRouter.Group("user")
		{
			user.POST("/login", admin.UserLogin)
		}
	}

	adminRouter := routerEngine.Group("admin", middlewares.CheckUserLogin, middlewares.AuditLog())
	{
		// user
		user := adminRouter.Group("user")
		{
			user.POST("/logout", admin.UserLogout)
			user.PUT("/change-password", admin.UserChangePassword)
			user.GET("/list", admin.UserList)
			user.GET("/info/:res_id", admin.UserInfo)
			user.POST("/add", admin.UserAdd)
			user.PUT("/update/:res_id", admin.UserUpdate)
			user.DELETE("/delete/:res_id", admin.UserDelete)
		}

		// service
		service := adminRouter.Group("service")
		{
			service.POST("/add", admin.ServiceAdd)
			service.GET("/list", admin.ServiceList)
			service.GET("/name/list", admin.ServiceNameList)
			service.GET("/info/:res_id", admin.ServiceInfo)
			service.PUT("/update/:res_id", admin.ServiceUpdate)
			service.DELETE("/delete/:res_id", admin.ServiceDelete)
			service.PUT("/update/name/:res_id", admin.ServiceUpdateName)
			service.PUT("/switch/enable/:res_id", admin.ServiceSwitchEnable)
			service.PUT("/switch/release/:res_id", admin.ServiceSwitchRelease)
		}

		servicePlugin := adminRouter.Group("service/plugin/config")
		{
			servicePlugin.POST("/add", admin.ServicePluginConfigAdd)
			servicePlugin.GET("/list/:service_res_id", admin.ServicePluginConfigList)
			servicePlugin.GET("/info/:res_id", admin.ServicePluginConfigInfo)
			servicePlugin.PUT("/update/:res_id", admin.ServicePluginConfigUpdate)
			servicePlugin.DELETE("/delete/:res_id", admin.ServicePluginConfigDelete)
			servicePlugin.PUT("/switch/enable/:res_id", admin.ServicePluginConfigSwitchEnable)
		}

		// router
		router := adminRouter.Group("router")
		{
			// router
			router.POST("/add", admin.RouterAdd)
			router.GET("/list", admin.RouterList)
			router.GET("/info/:router_res_id", admin.RouterInfo)
			router.PUT("/update/:router_res_id", admin.RouterUpdate)
			router.DELETE("/delete/:router_res_id", admin.RouterDelete)
			router.PUT("/update/name/:router_res_id", admin.RouterUpdateName)
			router.PUT("/switch/enable/:router_res_id", admin.RouterSwitchEnable)
			router.PUT("/switch/release/:router_res_id", admin.RouterSwitchRelease)
			router.POST("/copy/:router_res_id", admin.RouterCopy)
		}

		// router plugin
		routerPlugin := adminRouter.Group("router/plugin/config")
		{
			// router plugin
			routerPlugin.POST("/add", admin.RouterPluginConfigAdd)
			routerPlugin.GET("/list/:router_res_id", admin.RouterPluginConfigList)
			routerPlugin.GET("/info/:res_id", admin.RouterPluginConfigInfo)
			routerPlugin.PUT("/update/:res_id", admin.RouterPluginConfigUpdate)
			routerPlugin.DELETE("/delete/:res_id", admin.RouterPluginConfigDelete)
			routerPlugin.PUT("/switch/enable/:res_id", admin.RouterPluginConfigSwitchEnable)
		}

		// upstream
		upstream := adminRouter.Group("upstream")
		{
			upstream.POST("/add", admin.UpstreamAdd)
			upstream.GET("/list", admin.UpstreamList)
			upstream.GET("/info/:res_id", admin.UpstreamInfo)
			upstream.GET("/name/list", admin.UpstreamNameList)
			upstream.PUT("/update/:res_id", admin.UpstreamUpdate)
			upstream.DELETE("/delete/:res_id", admin.UpstreamDelete)
			upstream.PUT("/update/name/:res_id", admin.UpstreamUpdateName)
			upstream.PUT("/switch/enable/:res_id", admin.UpstreamSwitchEnable)
			upstream.PUT("/switch/release/:res_id", admin.UpstreamSwitchRelease)
		}

		// plugin
		plugin := adminRouter.Group("plugin")
		{
			plugin.GET("/type-list", admin.PluginTypeList)
			plugin.GET("/add-list", admin.PluginAddList)
			plugin.GET("/info/:plugin_res_id", admin.PluginInfo)
			// plugin.GET("/list", admin.PluginList)
			// plugin.PUT("/update/:id", admin.PluginUpdate)
			// plugin.DELETE("/delete/:id", admin.PluginDelete)
		}

		// global plugin
		globalPlugin := adminRouter.Group("global/plugin/config")
		{
			globalPlugin.POST("/add", admin.GlobalPluginConfigAdd)
			globalPlugin.GET("/list", admin.GlobalPluginConfigList)
			globalPlugin.GET("/info/:res_id", admin.GlobalPluginConfigInfo)
			globalPlugin.PUT("/update/:res_id", admin.GlobalPluginConfigUpdate)
			globalPlugin.DELETE("/delete/:res_id", admin.GlobalPluginConfigDelete)
			globalPlugin.PUT("/switch/enable/:res_id", admin.GlobalPluginConfigSwitchEnable)
		}

		// certificate
		certificate := adminRouter.Group("certificate")
		{
			certificate.GET("/list", admin.CertificateList)
			certificate.POST("/add", admin.CertificateAdd)
			certificate.GET("/info/:id", admin.CertificateInfo)
			certificate.PUT("/update/:id", admin.CertificateUpdate)
			certificate.DELETE("/delete/:id", admin.CertificateDelete)
			certificate.PUT("/switch/enable/:id", admin.CertificateSwitchEnable)
		}

		// letsencrypt
		letsencrypt := adminRouter.Group("letsencrypt")
		{
			letsencrypt.POST("/request", admin.LetsEncryptRequest)
		}

		// cluster node
		clusterNode := adminRouter.Group("cluster-node")
		{
			clusterNode.POST("/add", admin.ClusterNodeAdd)
			clusterNode.GET("/list", admin.ClusterNodeList)
			clusterNode.DELETE("/delete/:id", admin.ClusterNodeDelete)
		}

		// log
		log := adminRouter.Group("log")
		{
			log.GET("/list", admin.LogList)
			log.GET("/access/list", admin.AccessLogList)
			log.GET("/access/aggregation", admin.AccessLogAggregation)
			log.GET("/access/field/aggregation", admin.FieldAggregation)
		}
	}
}
