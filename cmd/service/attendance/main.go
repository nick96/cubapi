package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/nick96/cubapi/attendance"
	"github.com/nick96/cubapi/db"
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
		log.Fatalf("Failed to connect to databse on host %s: %v", os.Getenv("DB_HOST"), err)
	}

	attendanceStore := attendance.NewAttendanceStore(dbHandle)
	cubStore := attendance.NewCubStore(dbHandle)

	handler := attendance.NewHandler(cubStore, attendanceStore)

	log.Printf("Successfully start attendance service")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
