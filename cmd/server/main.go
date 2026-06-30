package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stefanSmk/privaquest/internal/config"
	"github.com/stefanSmk/privaquest/web"
	"github.com/stefanSmk/privaquest/internal/handler"
	"github.com/stefanSmk/privaquest/internal/middleware"
	"github.com/stefanSmk/privaquest/internal/service"
	"github.com/stefanSmk/privaquest/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	st := store.New(db)
	if err := st.Migrate(context.Background()); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	svc := service.New(st)
	api := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", api.Health)
	mux.HandleFunc("GET /api/locales", api.Locales)
	mux.Handle("POST /api/requests", middleware.OptionalPublicToken(cfg.PublicToken)(http.HandlerFunc(api.SubmitRequest)))
	mux.HandleFunc("GET /api/requests/status", api.CheckStatus)

	admin := middleware.AdminKey(cfg.AdminAPIKey)
	mux.Handle("GET /api/admin/dashboard", admin(http.HandlerFunc(api.AdminDashboard)))
	mux.Handle("GET /api/admin/requests", admin(http.HandlerFunc(api.AdminList)))
	mux.Handle("PATCH /api/admin/requests/{id}", admin(http.HandlerFunc(api.AdminUpdate)))
	mux.Handle("GET /api/admin/requests/{id}/audit", admin(http.HandlerFunc(api.AdminAudit)))

	static, _ := fs.Sub(web.Static, "static")
	mux.Handle("GET /", http.FileServer(http.FS(static)))

	srv := &http.Server{Addr: cfg.Addr, Handler: mux}

	go func() {
		log.Printf("privaquest running on %s", cfg.BaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
