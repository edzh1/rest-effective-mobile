package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/edzh1/rest-effective-mobile/internal"
	"github.com/edzh1/rest-effective-mobile/internal/models"
	"github.com/joho/godotenv"
)

type application struct {
	logger        *slog.Logger
	subscriptions *models.SubscriptionModel
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatalf("Wrong port %s", err)
	}

	dbCfg := internal.DSN{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBname:   os.Getenv("POSTGRES_DB"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := internal.InitDB(dbCfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer db.Close()

	app := application{
		logger:        logger,
		subscriptions: &models.SubscriptionModel{DB: db},
	}

	err = http.ListenAndServe(os.Getenv("ADDR"), app.routes())
	if err != nil {
		log.Fatal("can't start server " + err.Error())
	}
}
