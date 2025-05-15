package apiserver

import (
	"github.com/Ranper/iam/internal/apiserver/config"
	genericapiserver "github.com/Ranper/iam/internal/pkg/server"
)

type apiServer struct {
	genericAPIServer *genericapiserver.GenericAPIServer
}

// preparedAPIServer 包含已经准备好的apiserver. 该结构体只负责初始化之后的事项. 划分责任.
type preparedAPIServer struct {
	*apiServer
}

// ExtraConfig 包含除了通用配置之外, apiservre需要的配置
type ExtraConfig struct {
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		genericAPIServer: genericServer,
	}

	return server, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	initRouter(s.genericAPIServer.Engine)

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {

	return s.genericAPIServer.Run() // block, until stop.
}
