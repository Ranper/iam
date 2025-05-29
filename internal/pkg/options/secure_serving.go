package options

import (
	"fmt"
	"net"
	"path"

	"github.com/Ranper/iam/internal/pkg/server"
	"github.com/spf13/pflag"
)

/*
	抽象出SecureServingOptions, 从命令行、配置文件中获取配置
	再用配置项初始化SecureServingInfo
*/

type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port"`
	// 如果为true,表示BindPort不能为0
	Required bool
	// ServerCert
	ServerCert GeneratableKeyCert `json:"tls" mapstructure:"tls"`
}

// CertKey 包含证书相关的配置项
type CertKey struct {
	// 证书文件 - 证明服务器身份
	// CertFile 包含PEM-编码证书的文件, 可能还包含完整的证书链.
	CertFile string `json:"cert-file" mapstructure:"cert-file"`
	// 私钥文件 - 用于建立加密连接
	// Keyfile 包含PEM-编码的证书私钥的文件, 由CertFile指定.
	KeyFile string `json:"private-key-file" mapstructure:"private-key-file"`
}

type GeneratableKeyCert struct {
	// CertKey 允许设置一个显式的cert/key来使用
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	// 如果没有明确设置CertFile/KeyFile， CertDirectory指定一个目录来写入生成的证书。
	// PairName用于确定CertDirectory中的文件名。
	// 如果不配置“CertDirectory”和“PairName”，则生成内存证书。
	CertDirectory string `json:"cert-dir" mapstructure:"cert-dir"`
	// PairName是与CertDirectory一起用来生成证书和密钥文件名的名称。
	//变成CertDirectory/PairName。crt和CertDirectory/PairName.key
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// NewSecureServingOptions 生成一个默认的SecureServingOptions
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		Required:    true,
		ServerCert: GeneratableKeyCert{
			PairName: "iam",
			// /var/run/ 目录的用途为存储运行时数据（如PID文件、套接字文件、临时锁文件(防冲突)等）
			CertDirectory: "/var/run/iam",
		},
	}
}

// ApplyTo 根据SecureServingOptions的配置初始化服务器SecureServingInfo
// options: 用我从命令行或者配置文件读取的配置来初始化传入的Config
func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		CertKey: server.CertKey{
			CertFile: s.ServerCert.CertKey.CertFile,
			KeyFile:  s.ServerCert.CertKey.KeyFile,
		},
	}

	return nil
}

// Validate 校验SecureServingOptions的配置
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}
	if s.Required && s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				s.BindPort,
			),
		)
	} else if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf("--secure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port", s.BindPort),
		)
	}

	return errors
}

func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	desc := "The port on which to serve HTTPS with authentication and authorization."
	if s.Required {
		desc += " It cannot be switched off with 0."
	} else {
		desc += " if 0, don't serve HTTPS at all."
	}
	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}

func (s *SecureServingOptions) Complete() error {
	if s == nil || s.BindPort == 0 {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 { // 配置了证书和私钥路径
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}

// CreateListener create net listener by given address and returns it and port.
func CreateListener(addr string) (net.Listener, int, error) {
	network := "tcp"

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to listen on %v: %w", addr, err)
	}

	// get port
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()

		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}

	return ln, tcpAddr.Port, nil
}
