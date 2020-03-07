package main

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/nick96/cubapi/db"
	"github.com/nick96/cubapi/middleware"
	"github.com/nick96/cubapi/user"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	logger = logger.Named("user-service")

	dbHandle, err := db.NewConn(
		logger,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err != nil {
		logger.Fatal("Failed to connect to database")
	}

	store := user.NewStore(dbHandle)

	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(middleware.Logger(logger))
	router.Use(middleware.DefaultContentType(logger, "application/json"))

	router.Route("/user", user.NewUserRouter(logger, store))
	router.Route("/auth", user.NewAuthRouter(logger, store))

	logger.Info("Successfully started user service")
	logger.Fatal("Service exited with error", zap.Error(http.ListenAndServe(":8080", router)))
}
