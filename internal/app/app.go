package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httprouter "github.com/pliffdax/sparrow-api/internal/http"
	"github.com/pliffdax/sparrow-api/internal/util"
)

type App struct {
	httpServer *http.Server
}

func New() *App {
	version := util.Getenv("APP_VERSION", "0.1.0")
	_ = version

	addr := ":" + util.Getenv("PORT", "8080")

	r := httprouter.NewRouter()

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{httpServer: srv}
}

func (a *App) Run() error {
	go func() {
		log.Printf("sparrow-api listening on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err == http.ErrServerClosed {
			log.Printf("listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped gracefully, bye 👋")
	return nil
}
