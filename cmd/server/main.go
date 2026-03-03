package main

import (
	"context"
	"github.com/gorilla/mux"
	"go-rest-challenge/internal/middleware"
	"go-rest-challenge/internal/repository"
	"go-rest-challenge/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	httpTransport "go-rest-challenge/internal/transport/http"
)

func main() {
	repo := repository.NewMemoryBookRepository(1024)
	uc := usecase.NewBookUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/ping", handler.Ping).Methods("GET")
	r.HandleFunc("/echo", handler.Echo).Methods("POST")
	r.HandleFunc("/auth/token", handler.Token).Methods("POST")

	//r.HandleFunc("/books", handler.CreateBook).Methods("POST")
	r.Handle("/books", middleware.Auth(http.HandlerFunc(handler.GetBooks))).Methods("GET")

	//r.HandleFunc("/books", handler.GetBooks).Methods("GET")

	r.HandleFunc("/books/{id}", handler.GetBookByID).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", handler.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books/{id}", handler.DeleteBook).Methods(http.MethodDelete)

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
