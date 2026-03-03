package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go-rest-challenge/internal/middleware"
	"go-rest-challenge/internal/repository"
	httpTransport "go-rest-challenge/internal/transport/http"
	"go-rest-challenge/internal/usecase"
)

func main() {
	repo := repository.NewMemoryBookRepository(1024)
	uc := usecase.NewBookUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/ping", handler.Ping).Methods("GET")
	r.HandleFunc("/echo", handler.Echo).Methods("POST")
	r.HandleFunc("/auth/token", handler.Token).Methods("POST")

	// Public POST to create book
	r.HandleFunc("/books", handler.CreateBook).Methods("POST")

	// Protected GET routes
	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.Auth) // Apply auth middleware
	//protected.HandleFunc("/books", handler.GetBooks).Methods("GET")
	protected.HandleFunc("/books/{id}", handler.GetBookByID).Methods("GET")
	protected.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	protected.HandleFunc("/books/{id}", handler.DeleteBook).Methods("DELETE")

	server := &http.Server{
		Addr:              ":8099",
		Handler:           r,
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
