package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"gharabakloo/search/internal/config"
	"gharabakloo/search/internal/handler/http"
	"gharabakloo/search/internal/initial/db"
	"gharabakloo/search/internal/repository/cache/redis"
	"gharabakloo/search/internal/repository/db/mysql"
	"gharabakloo/search/internal/service/search"
	"gharabakloo/search/pkg/myerr"
)

func main() {
	log.Println("Running...")
	loadEnv()
	cfg := config.New()

	ctx := context.Background()
	mysqlDB, err := db.MySQL(ctx, cfg.DB.MySQL)
	errCheck(myerr.Errorf(err))

	redisClient, err := db.Redis(ctx, cfg.DB.Redis)
	errCheck(myerr.Errorf(err))

	dbRepository := mysql.New(mysqlDB)
	cacheRepository := redis.New(redisClient)
	srv := search.New(dbRepository, cacheRepository)

	gracefulShutdown()

	httpHandler := http.NewHandler(cfg, srv)
	err = http.RunServer(http.SetupHTTPRouter(httpHandler), cfg.HTTP)
	errCheck(myerr.Errorf(err))
}

// loadEnv loads .env file in the current path.
func loadEnv() {
	err := godotenv.Overload()
	errCheck(myerr.Errorf(err))
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		defer close(c)

		log.Println("Stopping...")
		time.Sleep(time.Second)
		log.Println("Stopped")
		os.Exit(0)
	}()
}
