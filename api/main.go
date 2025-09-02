package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"job_ping/database"
	"job_ping/models"
)

var db *sql.DB

func main() {
	var err error

	// Initialize database
	db, err = database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize database tables
	database.InitDB(db)

	// Set up Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Set up routes
	r.POST("/watch", createWatch)
	r.GET("/watches", listWatches)
	r.GET("/watch/:id", getWatch)
	r.PUT("/watch/:id", updateWatch)
	r.DELETE("/watch/:id", deleteWatch)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	log.Println("Starting API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// POST /watch
func createWatch(c *gin.Context) {
	var w models.Watch
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(
		"INSERT INTO watches (url, interval_ms, expected_status) VALUES ($1, $2, $3) RETURNING id, created_at",
		w.URL, w.IntervalMS, w.ExpectedStatus,
	).Scan(&w.ID, &w.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, w)
}

// GET /watches for all watches in db
func listWatches(c *gin.Context) {
	watches, err := database.GetAllWatches(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, watches)
}

// GET /watch/:id specfcic wathc retrival
func getWatch(c *gin.Context) {
	id := c.Param("id")

	var w models.Watch
	err := db.QueryRow("SELECT id, url, interval_ms, expected_status, created_at FROM watches WHERE id = $1", id).
		Scan(&w.ID, &w.URL, &w.IntervalMS, &w.ExpectedStatus, &w.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Watch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, w)
}

// PUT /watch/:id update a current watch
func updateWatch(c *gin.Context) {
	id := c.Param("id")
	var w models.Watch

	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("UPDATE watches SET url = $1, interval_ms = $2, expected_status = $3 WHERE id = $4",
		w.URL, w.IntervalMS, w.ExpectedStatus, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watch not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch updated successfully"})
}

// DELETE /watch/:id delete a current watched url
func deleteWatch(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM watches WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Watch not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch deleted successfully"})
}
