package main

import (
	"context"
	httpapi "github.com/makhmudovs1/go-tasks-api/internal/http"
	"github.com/makhmudovs1/go-tasks-api/internal/logging"
	"github.com/makhmudovs1/go-tasks-api/internal/service"
	"github.com/makhmudovs1/go-tasks-api/internal/storage/memory"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	logChan := make(chan logging.LogEvent, 10)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	logging.StartLogger(ctx, &wg, logChan)

	repo := memory.NewTaskRepo()
	svc := service.NewTaskService(repo)
	h := httpapi.NewHandler(svc, logChan)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	logChan <- logging.LogEvent{
		Time:    time.Now(),
		Action:  "TEST EVENT",
		Details: "First log",
	}
	logChan <- logging.LogEvent{
		Time:    time.Now(),
		Action:  "ANOTHER_EVENT",
		Details: "Another log message",
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	cancel()
	close(logChan)
	wg.Wait()

	log.Println("Server stopped")
}
