package utils

import (
	"apiok-admin/app/enums"
	"apiok-admin/app/packages"
	"bytes"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
)

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func RandomStrGenerate(len int) string {
	randomId := createRandomString(len)
	return strings.ToLower(randomId)
}

func IdGenerate(idType string) (string, error) {
	randomId := createRandomString(IdLength)

	var id string
	switch strings.ToLower(idType) {
	case IdTypeUser:
		id = IdTypeUser + "-" + randomId
	case IdTypeUserToken:
		id = IdTypeUserToken + "-" + randomId
	case IdTypeService:
		id = IdTypeService + "-" + randomId
	case IdTypeServiceDomain:
		id = IdTypeServiceDomain + "-" + randomId
	case IdTypeServiceNode:
		id = IdTypeServiceNode + "-" + randomId
	case IdTypeRouter:
		id = IdTypeRouter + "-" + randomId
	case IdTypePlugin:
		id = IdTypePlugin + "-" + randomId
	case IdTypePluginConfig:
		id = IdTypePluginConfig + "-" + randomId
	case IdTypeCertificate:
		id = IdTypeCertificate + "-" + randomId
	case IdTypeClusterNode:
		id = IdTypeClusterNode + "-" + randomId
	case IdTypeUpstream:
		id = IdTypeUpstream + "-" + randomId
	case IdTypeUpstreamNode:
		id = IdTypeUpstreamNode + "-" + randomId
	case IdTypeAcmeChallenge:
		id = IdTypeAcmeChallenge + "-" + randomId
	case IdTypeLog:
		id = IdTypeLog + "-" + randomId
	default:
		return "", fmt.Errorf("id type error")
	}

	return id, nil
}

func DiscernIP(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("(%s) is illegal ip", s)
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return IPV4, nil
		case ':':
			return IPV6, nil
		}
	}
	return "", nil
}

type CertificateInfo struct {
	CommonName    string
	NotBefore     time.Time
	NotAfter      time.Time
	KeyAlgorithm  string
	Issuer        string
}

func DiscernCertificate(certificate *string) (CertificateInfo, error) {
	certificateInfo := CertificateInfo{}
	pemBlock, _ := pem.Decode([]byte(*certificate))
	if pemBlock == nil {
		return certificateInfo, errors.New(enums.CodeMessages(enums.CertificateFormatError))
	}

	parseCert, parseCertErr := x509.ParseCertificate(pemBlock.Bytes)
	if parseCertErr != nil {
		return certificateInfo, errors.New(enums.CodeMessages(enums.CertificateParseError))
	}

	certificateInfo.CommonName = parseCert.Subject.CommonName
	certificateInfo.NotBefore = parseCert.NotBefore
	certificateInfo.NotAfter = parseCert.NotAfter
	certificateInfo.KeyAlgorithm = discernKeyAlgorithm(parseCert)
	certificateInfo.Issuer = discernIssuer(parseCert)

	return certificateInfo, nil
}

func discernIssuer(cert *x509.Certificate) string {
	issuer := &cert.Issuer
	if issuer.CommonName != "" {
		return issuer.CommonName
	}
	if len(issuer.Organization) > 0 {
		return issuer.Organization[0]
	}
	return ""
}

func discernKeyAlgorithm(cert *x509.Certificate) string {
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		bits := pub.N.BitLen()
		if bits >= 4096 {
			return "rsa4096"
		}
		if bits >= 3072 {
			return "rsa3072"
		}
		return "rsa2048"
	case *ecdsa.PublicKey:
		switch pub.Curve.Params().BitSize {
		case 256:
			return "ecdsa_p256"
		case 384:
			return "ecdsa_p384"
		case 521:
			return "ecdsa_p521"
		default:
			return "ecdsa"
		}
	default:
		return "unknown"
	}
}

type enumInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func LoadBalanceList() []enumInfo {
	loadBalanceList := []enumInfo{
		{Id: LoadBalanceRoundRobin, Name: LoadBalanceNameRoundRobin},
		{Id: LoadBalanceCHash, Name: LoadBalanceNameCHash},
	}

	return loadBalanceList
}

func ConfigBalanceList() []enumInfo {
	configBalanceList := []enumInfo{
		{Id: LoadBalanceRoundRobin, Name: ConfigBalanceNameRoundRobin},
		{Id: LoadBalanceCHash, Name: ConfigBalanceNameCHash},
	}

	return configBalanceList
}

func ConfigUpstreamNodeHealthList() []enumInfo {
	configHealthList := []enumInfo{
		{Id: HealthY, Name: ConfigHealthY},
		{Id: HealthN, Name: ConfigHealthN},
	}

	return configHealthList
}

func PluginAllTypes() []enumInfo {
	pluginTypeList := []enumInfo{
		{Id: PluginTypeIdAuth, Name: PluginTypeNameAuth},
		{Id: PluginTypeIdLimit, Name: PluginTypeNameLimit},
		{Id: PluginTypeIdSafety, Name: PluginTypeNameSafety},
		{Id: PluginTypeIdFlowControl, Name: PluginTypeNameFlowControl},
		{Id: PluginTypeIdOther, Name: PluginTypeNameOther},
	}

	return pluginTypeList
}

func PluginAllKeys() []string {
	pluginKeysList := []string{
		PluginKeyCors,
		PluginKeyMock,
		PluginKeyKeyAuth,
		PluginKeyJwtAuth,
		PluginKeyLimitReq,
		PluginKeyLimitConn,
		PluginKeyLimitCount,
		PluginKeyWaf,
		PluginKeyLogKafka,
		PluginKeyLogMysql,
		PluginKeyTrafficTag,
		PluginKeyRequestRewrite,
		PluginKeyResponseRewrite,
	}

	return pluginKeysList
}

type ConfigPluginData struct {
	ResID       string
	PluginKey   string
	Icon        string
	Type        int
	Description string
}

func AllConfigPluginData() []ConfigPluginData {

	allConfigPluginData := []ConfigPluginData{
		{
			ResID:       PluginIdCors,
			PluginKey:   PluginKeyCors,
			Icon:        PluginIconCors,
			Type:        PluginTypeIdSafety,
			Description: PluginDescCors,
		},
		{
			ResID:       PluginIdMock,
			PluginKey:   PluginKeyMock,
			Icon:        PluginIconMock,
			Type:        PluginTypeIdOther,
			Description: PluginDescMock,
		},
		{
			ResID:       PluginIdKeyAuth,
			PluginKey:   PluginKeyKeyAuth,
			Icon:        PluginIconKeyAuth,
			Type:        PluginTypeIdAuth,
			Description: PluginDescKeyAuth,
		},
		{
			ResID:       PluginIdJwtAuth,
			PluginKey:   PluginKeyJwtAuth,
			Icon:        PluginIconJwtAuth,
			Type:        PluginTypeIdAuth,
			Description: PluginDescJwtAuth,
		},
		{
			ResID:       PluginIdLimitReq,
			PluginKey:   PluginKeyLimitReq,
			Icon:        PluginIconLimitReq,
			Type:        PluginTypeIdLimit,
			Description: PluginDescLimitReq,
		},
		{
			ResID:       PluginIdLimitConn,
			PluginKey:   PluginKeyLimitConn,
			Icon:        PluginIconLimitConn,
			Type:        PluginTypeIdLimit,
			Description: PluginDescLimitConn,
		},
		{
			ResID:       PluginIdLimitCount,
			PluginKey:   PluginKeyLimitCount,
			Icon:        PluginIconLimitCount,
			Type:        PluginTypeIdLimit,
			Description: PluginDescLimitCount,
		},
		{
			ResID:       PluginIdWaf,
			PluginKey:   PluginKeyWaf,
			Icon:        PluginIconWaf,
			Type:        PluginTypeIdSafety,
			Description: PluginDescWaf,
		},
		{
			ResID:       PluginIdLogKafka,
			PluginKey:   PluginKeyLogKafka,
			Icon:        PluginIconLogKafka,
			Type:        PluginTypeIdOther,
			Description: PluginDescLogKafka,
		},
		{
			ResID:       PluginIdLogMysql,
			PluginKey:   PluginKeyLogMysql,
			Icon:        PluginIconLogMysql,
			Type:        PluginTypeIdOther,
			Description: PluginDescLogMysql,
		},
		{
			ResID:       PluginIdTrafficTag,
			PluginKey:   PluginKeyTrafficTag,
			Icon:        PluginIconTrafficTag,
			Type:        PluginTypeIdFlowControl,
			Description: PluginDescTrafficTag,
		},
		{
			ResID:       PluginIdRequestRewrite,
			PluginKey:   PluginKeyRequestRewrite,
			Icon:        PluginIconRequestRewrite,
			Type:        PluginTypeIdFlowControl,
			Description: PluginDescRequestRewrite,
		},
		{
			ResID:       PluginIdResponseRewrite,
			PluginKey:   PluginKeyResponseRewrite,
			Icon:        PluginIconResponseRewrite,
			Type:        PluginTypeIdFlowControl,
			Description: PluginDescResponseRewrite,
		},
	}

	return allConfigPluginData
}

func AllRequestMethod() []string {
	return []string{
		RequestMethodALL,
		RequestMethodGET,
		RequestMethodPOST,
		RequestMethodPUT,
		RequestMethodPATH,
		RequestMethodDELETE,
		RequestMethodOPTIONS,
	}
}

func ConfigAllRequestMethod() []string {
	return []string{
		RequestMethodGET,
		RequestMethodPOST,
		RequestMethodPUT,
		RequestMethodPATH,
		RequestMethodDELETE,
		RequestMethodOPTIONS,
	}
}

func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	srcMd5 := hex.EncodeToString(m.Sum(nil))

	return srcMd5
}

type ExpireToken struct {
	Expire int64
	Token  string
}

type TokenClaims struct {
	Encryption string `json:"encryption"`
	Timestamp  string `json:"timestamp"`
	Secret     string `json:"secret"`
	Issuer     string `json:"issuer"`
}

func GenToken(encryption string) (string, error) {
	tokenClaims := TokenClaims{
		Encryption: encryption,
		Timestamp:  Md5(strconv.FormatInt(time.Now().UnixNano(), 10)),
		Secret:     packages.Token.TokenSecret,
		Issuer:     packages.Token.TokenIssuer,
	}

	var (
		jsonValue []byte
		err       error
	)
	if jsonValue, err = json.Marshal(tokenClaims); err != nil {
		return "", err
	}

	token := strings.TrimRight(base64.URLEncoding.EncodeToString(jsonValue), "=")

	return token, nil
}

func ParseToken(tokenString string) (string, error) {
	if l := len(tokenString) % 4; l > 0 {
		tokenString += strings.Repeat("=", 4-l)
	}

	tokenStructStr, tokenStructStrErr := base64.URLEncoding.DecodeString(tokenString)
	if tokenStructStrErr != nil {
		return "", tokenStructStrErr
	}

	tokenClaims := TokenClaims{}
	unmarshalErr := json.Unmarshal(tokenStructStr, &tokenClaims)
	if unmarshalErr != nil {
		return "", unmarshalErr
	}

	if tokenClaims.Issuer != packages.Token.TokenIssuer || tokenClaims.Secret != packages.Token.TokenSecret {
		return "", errors.New("token parsing failed")
	}

	return tokenClaims.Encryption, nil
}

func IPNameToType(ipName string) (int, error) {
	iPNameToTypeMap := map[string]int{
		IPV4: IPTypeV4,
		IPV6: IPTypeV6,
	}

	ipType, ipTypeExist := iPNameToTypeMap[ipName]
	if ipTypeExist == false {
		return -1, errors.New("IP type does not exist")
	}

	return ipType, nil
}

func IpTypeNameList() (list []enumInfo) {
	list = []enumInfo{
		{Id: IPTypeV4, Name: IPV4},
		{Id: IPTypeV6, Name: IPV6},
	}

	return
}

func IpIdNameMap() (ipIdNameMap map[int]string) {
	ipTypeNameList := IpTypeNameList()

	ipIdNameMap = make(map[int]string)
	for _, IpTypeNameDetail := range ipTypeNameList {
		ipIdNameMap[IpTypeNameDetail.Id] = IpTypeNameDetail.Name
	}

	return
}

func IpNameIdMap() (nameIdMap map[string]int) {
	ipTypeNameList := IpTypeNameList()

	nameIdMap = make(map[string]int)
	for _, IpTypeNameDetail := range ipTypeNameList {
		nameIdMap[IpTypeNameDetail.Name] = IpTypeNameDetail.Id
	}

	return
}

func HealthTypeNameList() (list []enumInfo) {
	list = []enumInfo{
		{Id: HealthY, Name: HealthNameY},
		{Id: HealthN, Name: HealthNameN},
	}

	return
}

func HealthTypeNameMap() (healthNameMap map[int]string) {
	healthNameList := HealthTypeNameList()

	healthNameMap = make(map[int]string)
	for _, healthNameDetail := range healthNameList {
		healthNameMap[healthNameDetail.Id] = healthNameDetail.Name
	}

	return
}

func InterceptSni(domains []string) ([]string, error) {
	domainSniInfos := make([]string, 0)
	if len(domains) == 0 {
		return domainSniInfos, nil
	}

	tmpDomainSniMap := make(map[string]byte, 0)
	for _, domain := range domains {
		disassembleDomains := strings.Split(domain, ".")
		if len(disassembleDomains) < 2 {
			return domainSniInfos, errors.New(enums.CodeMessages(enums.ServiceDomainFormatError))
		}

		disassembleDomains[0] = "*"
		domainSniInfo := strings.Join(disassembleDomains, ".")

		_, exit := tmpDomainSniMap[domainSniInfo]
		if exit {
			continue
		}

		tmpDomainSniMap[domainSniInfo] = 0
		domainSniInfos = append(domainSniInfos, domainSniInfo)
	}

	return domainSniInfos, nil
}

func ParseSizeToBytes(sizeStr *string) (*int64, error) {
	if sizeStr == nil || *sizeStr == "" {
		return nil, nil
	}

	str := strings.TrimSpace(*sizeStr)
	str = strings.ToLower(str)

	if str == "0" || str == "" {
		return nil, nil
	}

	var multiplier int64 = 1
	var numStr string

	if strings.HasSuffix(str, "k") {
		multiplier = 1024
		numStr = strings.TrimSuffix(str, "k")
	} else if strings.HasSuffix(str, "m") {
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(str, "m")
	} else if strings.HasSuffix(str, "g") {
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(str, "g")
	} else {
		numStr = str
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size format: %s", str)
	}

	if num == 0 {
		return nil, nil
	}

	result := int64(num * float64(multiplier))
	return &result, nil
}
