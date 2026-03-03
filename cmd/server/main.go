package main

import (
	"context"
	"go-rest-challenge/internal/repository"
	"go-rest-challenge/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-rest-challenge/internal/middleware"
	httpTransport "go-rest-challenge/internal/transport/http"
)

func main() {
	repo := repository.NewMemoryBookRepository(1024)
	uc := usecase.NewBookUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handler.Ping)
	mux.HandleFunc("POST /echo", handler.Echo)
	mux.HandleFunc("POST /auth/token", handler.Token)

	// Wrap the entire resource or individual methods consistently
	mux.Handle("POST /books", middleware.Auth(http.HandlerFunc(handler.CreateBook)))
	mux.Handle("GET /books", middleware.Auth(http.HandlerFunc(handler.GetBooks)))

	// Don't forget the ID-based routes if the challenge requires it!
	mux.Handle("GET /books/{id}", middleware.Auth(http.HandlerFunc(handler.GetBookByID)))
	mux.Handle("PUT /books/{id}", middleware.Auth(http.HandlerFunc(handler.UpdateBook)))
	mux.Handle("DELETE /books/{id}", middleware.Auth(http.HandlerFunc(handler.DeleteBook)))

	server := &http.Server{
		Addr:              ":8099",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		log.Println("Server running on :8099")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
