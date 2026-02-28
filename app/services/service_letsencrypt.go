package services

import (
	"apiok-admin/app/models"
	"apiok-admin/app/packages"
	"apiok-admin/app/utils"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/gorm"
)

type LetsEncryptConfig struct {
	Enabled         bool
	Email           string
	UseStaging      bool
	CertDir         string
	RenewBeforeDays int
}

type LetsEncryptService struct {
	client       *acme.Client
	manager      *autocert.Manager
	challenges   map[string]string // token -> keyAuthorization
	challengesMu sync.RWMutex
	config       *LetsEncryptConfig
}

var (
	letsEncryptService *LetsEncryptService
	letsEncryptOnce    sync.Once
)

func NewLetsEncryptService() *LetsEncryptService {
	letsEncryptOnce.Do(func() {
		letsEncryptService = &LetsEncryptService{
			challenges: make(map[string]string),
		}
	})
	return letsEncryptService
}

// InitLetsEncrypt 初始化Let's Encrypt客户端
func (s *LetsEncryptService) InitLetsEncrypt(config *LetsEncryptConfig) error {
	if !config.Enabled {
		return nil
	}

	if config.Email == "" {
		return errors.New("Let's Encrypt email is required")
	}

	s.config = config

	// 创建证书存储目录
	if err := os.MkdirAll(config.CertDir, 0755); err != nil {
		return fmt.Errorf("failed to create cert directory: %v", err)
	}

	// 选择ACME服务器
	var acmeURL string
	if config.UseStaging {
		acmeURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else {
		acmeURL = "https://acme-v02.api.letsencrypt.org/directory"
	}

	// 创建ACME客户端
	client := &acme.Client{
		DirectoryURL: acmeURL,
		Key:          s.getOrCreateAccountKey(),
		UserAgent:    "apiok-admin/1.0",
	}

	// 注册账户
	ctx := context.Background()
	account := &acme.Account{
		Contact: []string{"mailto:" + config.Email},
	}
	_, err := client.Register(ctx, account, acme.AcceptTOS)
	if err != nil {
		// 如果账户已存在，这是正常情况，继续使用现有账户
		errMsg := strings.ToLower(err.Error())
		if !strings.Contains(errMsg, "already registered") &&
			!strings.Contains(errMsg, "account already exists") {
			return fmt.Errorf("failed to register account: %v", err)
		}
		packages.Log.Infof("ACME account already exists, using existing account")
	}

	s.client = client

	// 创建autocert管理器（用于HTTP-01验证）
	s.manager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(), // 允许所有域名
		Cache:      autocert.DirCache(config.CertDir),
		Email:      config.Email,
	}

	return nil
}

// getOrCreateAccountKey 获取或创建账户私钥
func (s *LetsEncryptService) getOrCreateAccountKey() *rsa.PrivateKey {
	keyPath := filepath.Join(s.config.CertDir, "account.key")

	// 尝试读取现有密钥
	if data, err := os.ReadFile(keyPath); err == nil {
		block, _ := pem.Decode(data)
		if block != nil {
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err == nil {
				return key
			}
		}
	}

	// 创建新密钥
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("failed to generate account key: %v", err))
	}

	// 保存密钥
	keyData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err := os.WriteFile(keyPath, keyData, 0600); err != nil {
		packages.Log.Warnf("failed to save account key: %v", err)
	}

	return key
}

// RequestCertificate 申请Let's Encrypt证书
func (s *LetsEncryptService) RequestCertificate(domain string, enable bool) (string, error) {
	if s.client == nil {
		conf := packages.GetConfig()
		if conf == nil {
			return "", errors.New("配置未加载，无法使用 Let's Encrypt")
		}

		confValue := reflect.ValueOf(conf).Elem()
		letsEncryptField := confValue.FieldByName("LetsEncrypt")
		if !letsEncryptField.IsValid() {
			return "", errors.New("Let's Encrypt 配置不存在")
		}

		enabled := letsEncryptField.FieldByName("Enabled").Bool()
		if !enabled {
			return "", errors.New("Let's Encrypt 未启用，请在 config/app.yaml 中设置 letsencrypt.enabled: true 并配置邮箱后重启服务")
		}

		email := letsEncryptField.FieldByName("Email").String()
		if email == "" {
			return "", errors.New("Let's Encrypt 邮箱未配置，请在 config/app.yaml 中设置 letsencrypt.email")
		}

		useStaging := letsEncryptField.FieldByName("UseStaging").Bool()
		certDir := letsEncryptField.FieldByName("CertDir").String()
		if certDir == "" {
			certDir = "./certs"
		}
		renewBeforeDays := int(letsEncryptField.FieldByName("RenewBeforeDays").Int())
		if renewBeforeDays == 0 {
			renewBeforeDays = 30
		}

		config := &LetsEncryptConfig{
			Enabled:         enabled,
			Email:           email,
			UseStaging:      useStaging,
			CertDir:         certDir,
			RenewBeforeDays: renewBeforeDays,
		}

		if err := s.InitLetsEncrypt(config); err != nil {
			return "", fmt.Errorf("Let's Encrypt 初始化失败: %v", err)
		}
	}

	ctx := context.Background()

	// 创建证书请求
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// 创建证书签名请求
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domain,
		},
		DNSNames: []string{domain},
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		return "", fmt.Errorf("failed to create CSR: %v", err)
	}

	// 创建订单
	order, err := s.client.AuthorizeOrder(ctx, acme.DomainIDs(domain))
	if err != nil {
		return "", fmt.Errorf("failed to authorize order: %v", err)
	}

	// 处理授权挑战
	var challengeTokens []string
	for _, authzURL := range order.AuthzURLs {
		authz, err := s.client.GetAuthorization(ctx, authzURL)
		if err != nil {
			return "", fmt.Errorf("failed to get authorization: %v", err)
		}

		// 查找HTTP-01挑战
		var challenge *acme.Challenge
		for _, ch := range authz.Challenges {
			if ch.Type == "http-01" {
				challenge = ch
				break
			}
		}

		if challenge == nil {
			return "", errors.New("no http-01 challenge found")
		}

		// 生成验证token
		token, err := s.client.HTTP01ChallengeResponse(challenge.Token)
		if err != nil {
			return "", fmt.Errorf("failed to generate challenge response: %v", err)
		}

		// 存储挑战信息（用于HTTP验证路由）
		expiredAt := time.Now().Add(1 * time.Hour)
		challengeModel := &models.AcmeChallenges{}
		if err := challengeModel.ChallengeAdd(challenge.Token, token, expiredAt); err != nil {
			return "", fmt.Errorf("failed to save challenge to database: %v", err)
		}

		s.challengesMu.Lock()
		s.challenges[challenge.Token] = token
		s.challengesMu.Unlock()
		challengeTokens = append(challengeTokens, challenge.Token)

		// 接受挑战
		_, err = s.client.Accept(ctx, challenge)
		if err != nil {
			return "", fmt.Errorf("failed to accept challenge: %v", err)
		}

		// 等待挑战完成
		_, err = s.client.WaitAuthorization(ctx, authz.URI)
		if err != nil {
			return "", fmt.Errorf("failed to wait authorization: %v", err)
		}
	}

	// 完成订单并获取证书
	der, _, err := s.client.CreateOrderCert(ctx, order.FinalizeURL, csr, true)
	if err != nil {
		return "", fmt.Errorf("failed to create order cert: %v", err)
	}

	// 解析证书
	cert, err := x509.ParseCertificate(der[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %v", err)
	}

	// 转换为PEM格式
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der[0],
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	// 保存到数据库
	issuer := ""
	if len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	}
	if issuer == "" && cert.Issuer.CommonName != "" {
		issuer = cert.Issuer.CommonName
	}

	certificateData := &models.Certificates{
		Certificate:   string(certPEM),
		PrivateKey:    string(keyPEM),
		ExpiredAt:     cert.NotAfter,
		Enable:        utils.EnableOff, // 默认不启用，等待用户确认
		Sni:           domain,
		CaProvider:    "letsencrypt",
		KeyAlgorithm:  "rsa2048",
		Issuer:        issuer,
	}

	err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
		resID, err := (&models.Certificates{}).CertificatesAdd(tx, certificateData)
		if err != nil {
			return err
		}

		certificateData.ResID = resID

		// 如果请求时指定启用，则启用证书
		if enable {
			if err = SyncDataSideCertificate(tx, certificateData, ""); err != nil {
				return err
			}
			certificateData.Enable = utils.EnableOn
			if err = (&models.Certificates{}).CertificatesUpdate(tx, resID, certificateData); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// 清理挑战信息
	challengeModel := &models.AcmeChallenges{}
	s.challengesMu.Lock()
	for _, token := range challengeTokens {
		delete(s.challenges, token)
		challengeModel.ChallengeDelete(token)
	}
	s.challengesMu.Unlock()

	return certificateData.ResID, nil
}

// GetChallengeToken 获取HTTP-01挑战token
func (s *LetsEncryptService) GetChallengeToken(token string) (string, bool) {
	s.challengesMu.RLock()
	keyAuth, ok := s.challenges[token]
	s.challengesMu.RUnlock()

	if ok {
		return keyAuth, true
	}

	challengeModel := &models.AcmeChallenges{}
	keyAuth, err := challengeModel.ChallengeGet(token)
	if err != nil {
		return "", false
	}

	s.challengesMu.Lock()
	s.challenges[token] = keyAuth
	s.challengesMu.Unlock()

	return keyAuth, true
}

// RenewExpiringCertificates 续期即将过期的证书
func (s *LetsEncryptService) RenewExpiringCertificates() error {
	if s.client == nil {
		return nil
	}

	// 查询即将过期的证书
	renewBefore := time.Now().AddDate(0, 0, s.config.RenewBeforeDays)

	var certificates []models.Certificates
	err := packages.GetDb().Table((&models.Certificates{}).TableName()).
		Where("expired_at <= ? AND expired_at > ?", renewBefore, time.Now()).
		Where("enable = ?", utils.EnableOn).
		Where("ca_provider = ?", "letsencrypt").
		Find(&certificates).Error

	if err != nil {
		return fmt.Errorf("failed to query expiring certificates: %v", err)
	}

	for _, cert := range certificates {
		packages.Log.Infof("Renewing certificate for domain: %s", cert.Sni)

		// 申请新证书
		resID, err := s.RequestCertificate(cert.Sni, false)
		if err != nil {
			packages.Log.Errorf("Failed to renew certificate for %s: %v", cert.Sni, err)
			continue
		}

		// 获取新证书信息
		newCert, err := (&models.Certificates{}).CertificateInfoById(resID)
		if err != nil {
			packages.Log.Errorf("Failed to get new certificate info: %v", err)
			continue
		}

		// 更新旧证书为新证书并启用
		err = packages.GetDb().Transaction(func(tx *gorm.DB) error {
			// 禁用旧证书
			if err := (&models.Certificates{}).CertificateSwitchEnable(tx, cert.ResID, utils.EnableOff); err != nil {
				return err
			}

			// 删除数据面旧证书
			apiokDataModel := models.ApiokData{}
			if err := apiokDataModel.Delete("certificates", cert.ResID); err != nil {
				packages.Log.Warnf("Failed to delete old certificate from data side: %v", err)
			}

			// 启用新证书
			if err := SyncDataSideCertificate(tx, &newCert, ""); err != nil {
				return err
			}

			newCert.Enable = utils.EnableOn
			if err := (&models.Certificates{}).CertificatesUpdate(tx, resID, &newCert); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			packages.Log.Errorf("Failed to switch to new certificate: %v", err)
		} else {
			packages.Log.Infof("Successfully renewed certificate for domain: %s", cert.Sni)
		}
	}

	return nil
}

