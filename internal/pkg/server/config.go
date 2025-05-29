package server

import (
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Config 通用API服务器的配置
type Config struct {
	InsecureServing *InsecureServingInfo
	SecureServing   *SecureServingInfo

	// 基础类型, 直接初始化,不需要单独封装为文件了
	Mode        string
	Middlewares []string
	Healthz     bool
}

// 自签名证书（用于开发测试）
// 由 CA 颁发的证书（用于生产环境）

// CertKey 包含证书相关的配置项
type CertKey struct {
	// 证书文件 - 证明服务器身份
	// CertFile 包含PEM-编码证书的文件, 可能还包含完整的证书链.
	CertFile string
	// 私钥文件 - 用于建立加密连接
	// Keyfile 包含PEM-编码的证书私钥的文件, 由CertFile指定.
	KeyFile string
}

type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 非安全服务器的配置
type InsecureServingInfo struct {
	Address string
}

func NewConfig() *Config {
	// 默认值
	return &Config{
		Healthz:     true,
		Mode:        gin.ReleaseMode,
		Middlewares: []string{},
	}
}

// CompletedConfig GenericAPIServer完整的配置
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// Complete 根据指定的配置,生成一个GenericAPIServer实例
func (c CompletedConfig) New() (*GenericAPIServer, error) {
	// setMode before gin.New()
	gin.SetMode(c.Mode)

	//! 用config的参数来初始化GenericAPIServer.
	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
