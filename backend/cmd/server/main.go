package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Luawig/neoneuro/backend/internal/server"
	"github.com/Luawig/neoneuro/backend/pkg/config"
)

func main() {
	cfg := config.Load()
	engine := server.NewEngine(cfg)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: engine}
	go func() {
		log.Printf("▶ http %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
