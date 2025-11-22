package main

import (
	"LinksChecker/internal/delivery"
	"LinksChecker/internal/repository"
	"LinksChecker/internal/repository/inmemory"
	"LinksChecker/internal/service/checker"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	graceful = flag.Bool("graceful", false, "enable graceful restart")
)

func main() {
	flag.Parse()

	repo := inmemory.New("state.json")
	checker := checker.New(repo)
	handler := delivery.NewHandler(checker)

	port := os.Getenv("AUTH_PORT")

	server := &http.Server{
		Addr:    port,
		Handler: nil,
	}

	http.HandleFunc("POST /check", handler.CheckLinks)
	http.HandleFunc("GET /report", handler.GenerateReport)

	// Channels for graceful shutdown/restart
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	restart := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(restart, syscall.SIGHUP)

	go func() {
		log.Println("Server is running on port", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	go func() {
		for {
			select {
			case <-quit:
				log.Println("Received shutdown signal...")
				handler.SetReady(false)
				gracefulShutdown(server, handler, repo, done)
				return

			case <-restart:
				log.Println("Received restart signal...")
				handler.SetReady(false)
				gracefulRestart(server, handler, repo)
				return
			}
		}
	}()

	// Awaiting completion
	<-done
	log.Println("Server stopped")
}

func gracefulShutdown(server *http.Server, handler *delivery.Handler, repo repository.Repo, done chan bool) {
	// Waiting for active tasks to complete (maximum 30 seconds)
	if handler.WaitForActiveTasks(30) {
		log.Println("All active tasks completed successfully")
	} else {
		log.Println("Some tasks were not completed within timeout")
	}

	// Saving the state before stopping
	log.Println("Saving state before shutdown...")
	repo.SaveState()

	// Stopping the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the server: %v\n", err)
	}

	close(done)
}

func gracefulRestart(server *http.Server, handler *delivery.Handler, repo repository.Repo) {
	// Waiting for active tasks to be completed
	if handler.WaitForActiveTasks(30) {
		log.Println("All active tasks completed, proceeding with restart")
	} else {
		log.Println("Proceeding with restart despite active tasks")
	}

	// Saving the state before restarting
	log.Println("Saving state before restart...")
	repo.SaveState()

	// Stopping the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the server for restart: %v\n", err)
	}

	// Launching a new instance
	execName, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get executable name: %v\n", err)
	}

	args := []string{"-graceful"}
	cmd := exec.Command(execName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not restart server: %v\n", err)
	}

	log.Println("New server process started, PID:", cmd.Process.Pid)
	os.Exit(0)
}
