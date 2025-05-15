package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Ranper/iam/internal/pkg/middleware"
	"github.com/Ranper/iam/pkg/core"
	"github.com/Ranper/iam/pkg/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// 通用的api服务器

type GenericAPIServer struct {
	middlewares []string

	InsecureServingInfo *InsecureServingInfo

	*gin.Engine // 包装了一下gin.Engine
	healthz     bool

	insecureServer *http.Server
}

func initGenericAPIServer(s *GenericAPIServer) {
	// do some setup

	s.InstallMiddlewares()
	s.InstallAPIs()
}

// InstallMiddlewares 根据需要安装通用的中间件
func (s *GenericAPIServer) InstallMiddlewares() {
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware %s", m)

			continue
		}

		log.Infof("install middleware %s", m)
		s.Use(mw)
	}
}

func (s *GenericAPIServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}
}

func (s *GenericAPIServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		log.Infof("Start to listen the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	// 测试router是否生效
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *GenericAPIServer) Close() {
	ctx, calcel := context.WithTimeout(context.Background(), 10*time.Second)
	defer calcel()

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}

func (s *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Infof("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep 1 second and retry.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done(): // 已经超时了，直接退出
			log.Fatal("can not ping http server within t he specified time interval.")
		default:
		}
	}
}
