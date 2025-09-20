package main

import (
	"context"
	"log"
	"spotify-api/app"
)

var ctx = context.Background()

func main() {
	app, err := app.InitApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.Router.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
