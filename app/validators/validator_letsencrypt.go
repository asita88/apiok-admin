package validators

type LetsEncryptRequest struct {
	Domain string `json:"domain" zh:"域名" en:"Domain name" binding:"required"`
	Enable int    `json:"enable" zh:"是否启用" en:"Enable certificate" binding:"omitempty,oneof=1 2"`
}
