package server

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

const (
	maxConns    = 100
	defaultPort = "7777"
	envFile     = ".env"
)

func NewServer() (*Config, error) {
	godotenv.Load(envFile)

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Println("WARNING: TOKEN_SECRET environment variable not set, using 'secret' as token secret")
		tokenSecret = "secret"

	}

	environment := os.Getenv("ENVIRONMENT")
	db_url := os.Getenv("GOOSE_DBSTRING")
	if environment == "test" {
		db_url = os.Getenv("TEST_DB_URL")
	}
	if db_url == "" {
		return nil, errors.New("error: GOOSE_DBSTRING environment variable not set")
	}

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		return nil, errors.New("error: failed to connect to database: " + err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("WARNING: PORT environment variable not set")
		port = defaultPort
	}

	conns := make(activeConnections, maxConns)
	redisUrl := os.Getenv("REDIS_URL")
	if environment == "test" {
		redisUrl = "localhost:6379"
	}
	if redisUrl == "" {
		return nil, errors.New("error: REDIS_URL environment variable not set")
	}
	redis := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	queries := database.New(db)
	c := Config{
		TokenSecret: tokenSecret,
		Port:        port,
		Users:       conns,
		Rdb:         redis,
		DB:          queries,
	}

	return &c, nil
}
