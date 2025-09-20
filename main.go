package main

import (
	"context"
	"log"
	"os"
	"spotify-api/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	redis_addr := os.Getenv("REDIS_ADDR")
	redis_pass := os.Getenv("REDIS_PASSWORD")
	redis_db := 0 // os.Getenv("REDIS_DB")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       redis_db,
	})

	err := rdb.Set(ctx, "ping", "pong", 0).Err()
	if err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/musics", handlers.GetMusics)
	router.GET("/musics/:playlistID", handlers.GetPlaylist)

	if err := router.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
