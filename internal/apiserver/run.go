package apiserver

import "github.com/Ranper/iam/internal/apiserver/config"

// Run 运行一个指定的APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
