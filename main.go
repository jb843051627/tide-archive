package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jb843051627/tide-archive/internal/api"
	"github.com/jb843051627/tide-archive/internal/service"
	"github.com/jb843051627/tide-archive/internal/store"
)

func main() {
	address := flag.String("addr", ":8098", "HTTP listen address")
	dbPath := flag.String("db", "tide-archive.db", "SQLite database path")
	flag.Parse()
	repo, err := store.Open(filepath.Clean(*dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	app := service.New(repo)
	defer app.Close()
	server := &http.Server{Addr: *address, Handler: api.New(app).Handler()}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server: %v", err)
		}
	}()
	<-ctx.Done()
	_ = server.Shutdown(context.Background())
}
