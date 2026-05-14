// Package app provides the application setup.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nabilfikrisp/sv-be/config"
	"github.com/nabilfikrisp/sv-be/internal/controller/restapi"
	"github.com/nabilfikrisp/sv-be/internal/repo/inmem"
	// "github.com/nabilfikrisp/sv-be/internal/repo/persistent"
	"github.com/nabilfikrisp/sv-be/internal/usecase/post"
	"github.com/nabilfikrisp/sv-be/pkg/httpserver"
	"github.com/nabilfikrisp/sv-be/pkg/logger"
	// "github.com/nabilfikrisp/sv-be/pkg/mysql"
)

type useCases struct {
	post *post.UseCase
}

type servers struct {
	http *httpserver.Server
}

func initUseCases() useCases {
	postRepo := inmem.NewPostInMemRepo()

	return useCases{
		post: post.New(postRepo),
	}
}

func initServer(cfg *config.Config, uc useCases, l logger.Interface) servers {
	httpserver := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpserver.Engine, cfg, uc.post, l)

	return servers{
		http: httpserver,
	}
}

func (s *servers) startServers() {
	s.http.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - shutdownServers - httpServer.Shutdown: %w", err))
	}
}

// Run initializes application dependencies, starts the server, and waits for shutdown.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	// mysql, err := mysql.New(cfg.MySQL.URL, mysql.MaxPoolSize(cfg.MySQL.PoolMax))
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - mysql.New: %w", err))
	// }
	// defer mysql.Close()

	// uc := initUseCases(mysql)
	uc := initUseCases()
	s := initServer(cfg, uc, l)
	s.startServers()
	s.waitForShutdown(l)
}
