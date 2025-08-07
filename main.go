package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

// Initialize Redis client
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default DB
	})
}

// SearchResult represents a bus search result
type SearchResult struct {
	BusID     string `json:"bus_id"`
	BusName   string `json:"bus_name"`
	Price     string `json:"price"`
	Departure string `json:"departure"`
}

// Function to simulate a database query for search results
func queryDatabase(origin, destination, date string) ([]SearchResult, error) {
	// Simulate a database query with hardcoded values for demonstration
	results := []SearchResult{
		{BusID: "1", BusName: "AC Express", Price: "$25", Departure: "2025-04-17 10:00:00"},
		{BusID: "2", BusName: "Super Deluxe", Price: "$30", Departure: "2025-04-17 12:00:00"},
	}
	return results, nil
}

// Function to search bus tickets
func searchTickets(c *gin.Context) {
	origin := c.DefaultQuery("origin", "New York")
	destination := c.DefaultQuery("destination", "Boston")
	date := c.DefaultQuery("date", "2025-04-17")

	cacheKey := fmt.Sprintf("%s-%s-%s", origin, destination, date)

	// Check if the results are already cached in Redis
	cachedResults, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Cache miss, query the database
		results, err := queryDatabase(origin, destination, date)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to query the database"})
			return
		}

		// Cache the search results for 1 hour
		rdb.Set(ctx, cacheKey, fmt.Sprintf("%v", results), 1*time.Hour)

		// Return the database query results
		c.JSON(200, results)
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Failed to connect to Redis"})
		return
	} else {
		// Cache hit, return cached data (assuming it's in a string format)
		c.JSON(200, cachedResults)
	}
}

func main() {
	// Initialize Redis
	initRedis()

	// Create a Gin router
	r := gin.Default()

	// Define search route
	r.GET("/searchcache", searchTickets)

	// Start the server
	if err := r.Run(":8060"); err != nil {
		log.Fatalf("Error starting Gin server: %v", err)
	}
}
