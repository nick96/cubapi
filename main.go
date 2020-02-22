package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/nick96/cubapi/repo"
)

func main() {
	db, err := DBConn(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database")

	err = InitDB(db)
	if err != nil {
		log.Fatalf("Failed to initialise the database: %v", err)
	}
	log.Printf("Successfully initialised database")

	cubStore := repo.CubStore{db}
	attendanceStore := repo.AttendanceStore{db}

	router := gin.Default()
	router.POST("/attendance", AttendanceHandler(cubStore, attendanceStore))

	port := os.Getenv("APP_PORT")
	if strings.TrimSpace(port) {
		port = "8080"
	}
	log.Printf("Starting app on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
