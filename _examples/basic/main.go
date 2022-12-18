package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devfacet/ginman"
	"github.com/gin-gonic/gin"
)

type GetRes struct {
	ginman.Response
	Message string `json:"message,omitempty"`
}

func main() {
	r, err := ginman.NewWithOptions(ginman.Options{
		ContextMetadata:   map[string]any{"env": "dev"},
		EnableCompression: true,
		EnableLocation:    true,
		EnableRecovery:    true,
		EnableRequestID:   true,
		Mode:              "dev",
		Validations:       []string{"base64Any", "duration", "json"},
	})
	if err != nil {
		log.Fatalf("couldn't initialize the server due to %s", err)
	}
	r.Handle("GET", "/", func(c *gin.Context) {
		r := GetRes{
			Response: ginman.Response{Code: http.StatusOK},
			Message:  "Hello there",
		}
		r.Reply(c, r)
	})

	serverAddress, ok := os.LookupEnv("APP_SERVER_ADDRESS")
	if !ok {
		serverAddress = "localhost:8080"
	}
	log.Default().Printf("listening for request on %s", serverAddress)
	if err := r.Run(serverAddress); err != nil {
		log.Fatalf("couldn't start the server due to %s", err)
	}
}
