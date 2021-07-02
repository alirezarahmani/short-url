package main

import (
	"fmt"
	"github.com/alirezarahmani/short-url/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alirezarahmani/short-url/api"
	redisRepository "github.com/alirezarahmani/short-url/repository/redis"
)

func main() {
	repo := Repo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), router)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s uiu", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func Repo() shortener.RedirectRepository {
	redisURL := os.Getenv("REDIS_URL")
	repo, err := redisRepository.NewRedisRepository(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
