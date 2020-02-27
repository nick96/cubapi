package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nick96/cubapi/db"
)

func main() {
	handle, err := db.DBConn(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer handle.Close()
	log.Printf("Successfully connected to database")

	err = db.InitDB(handle)
	if err != nil {
		log.Fatalf("Failed to initialise the database: %v", err)
	}
	log.Printf("Successfully initialised database")

	cubStore := CubStore{handle}
	attendanceStore := AttendanceStore{handle}

	attendanceHandler := NewAttendanceHandler(cubStore, attendanceStore)
	cubsHandler := NewCubsHandler(cubStore)

	router := gin.Default()
	router.POST("/attendance", attendanceHandler)
	router.GET("/cub", cubsHandler)
	router.GET("/cub/:id", cubHandlerji)

	port := os.Getenv("APP_PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}
	log.Printf("Starting app on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
