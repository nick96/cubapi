package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/nick96/cubapi/attendance"
	"github.com/nick96/cubapi/db"
	"github.com/nick96/cubapi/middleware"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	logger = logger.Named("attendance-service")

	dbHandle, err := db.NewConn(
		logger,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// attendanceStore := attendance.NewAttendanceStore(dbHandle)
	// cubStore := attendance.NewCubStore(dbHandle)

	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(middleware.Logger(logger))
	router.Use(middleware.DefaultContentType(logger, "application/json"))

	// handler := attendance.NewHandler(cubStore, attendanceStore)
	// router.Route("/attendance", handler)

	logger.Info("Successfully start attendance service")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Service exited with an error", zap.Error(err))
	}
}
