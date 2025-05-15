package server

import "github.com/gin-gonic/gin"

// Config 通用API服务器的配置
type Config struct {
	InsecureServing *InsecureServingInfo

	// 基础类型, 直接初始化,不需要单独封装为文件了
	Mode        string
	Middlewares []string
	Healthz     bool
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
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
