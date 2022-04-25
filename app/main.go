package main

import (
	"fmt"
	"net/http"
	"os"

	redis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var rdb *redis.Client

const (
	defaultIPStackURL  = "http://api.ipstack.com/"
	defaultServicePort = "8080"
)

func main() {
	port := os.Getenv("SVC_PORT")
	if port == "" {
		port = defaultServicePort
	}

	setupRoutes()
	connectToCache()
	logrus.Infof("svc-ip-location started on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
