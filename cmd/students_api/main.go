package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/surajNirala/students_api/internal/config"
	"github.com/surajNirala/students_api/internal/http/handlers/student"
	"github.com/surajNirala/students_api/internal/storage/sqlite"
)

func main() {
	// fmt.Println("Welcome to the golang rest api")
	//Load Config
	cfg := config.MustLoad()
	//Database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("Version", "1.0.0"))
	//Setup Route
	router := http.NewServeMux()
	router.HandleFunc("GET /api/students", student.GetStudentList(storage))
	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudentById(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudentById(storage))
	//Setup Server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started", slog.String("address", cfg.Addr))
	// fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server")
		}
	}()
	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown successfully.")

}
