package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MSEarn/go-neo4j/config"
	"github.com/MSEarn/go-neo4j/pkg/auth"
	"github.com/MSEarn/go-neo4j/pkg/neo4j_driver"
	"github.com/MSEarn/go-neo4j/routes"
	"github.com/MSEarn/go-neo4j/version"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	httpServer *http.Server
)

func main() {
	cfg, err := InitConfig()
	if err != nil {
		panic(err)
	}

	jwt, err := auth.NewJWT("HS256", []byte("cfg.Auth.JWTSecret"), "15m")
	if err != nil {
		zap.L().Error(fmt.Sprintf("unable to newJWT, err: %v", err))
		panic(err)
	}

	neo4jCfg := cfg.Neo4j
	neo4jDriver := neo4j_driver.NewDriver()

	router := routes.Setup(neo4jCfg, neo4jDriver, jwt)
	httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Server.Port),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(cfg.Server.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.Server.ReadTimeout),

		Handler: router, // Pass our instance of gorilla/mux in.
	}

	fmt.Printf(
		"Starting the service...\ncommit: %s, build time: %s, release: %s\n",
		version.GitCommit, version.Buildtime, version.Version,
	)

	fmt.Printf("server serve on :%d\n", cfg.Server.Port)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().Error("HTTP Server listen failed", zap.Error(err))
			return err
		}
		return nil
	})
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}

func InitConfig() (*config.Config, error) {
	viper.SetDefault("NEO4J.URI", "bolt://localhost:7687")
	viper.SetDefault("NEO4J.USERNAME", "neo4j_admin1")
	viper.SetDefault("NEO4J.PASSWORD", "Passw0rd")

	viper.SetDefault("SERVER.PORT", 8091)
	viper.SetDefault("SERVER.WRITETIMEOUT", 100)
	viper.SetDefault("SERVER.READTIMEOUT", 100)
	viper.SetDefault("SERVER.IDLETIMEOUT", 100)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
