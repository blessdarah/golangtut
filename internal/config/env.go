package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppEnv struct {
	AppPort           string
	AppHost           string
	DB_HOST           string
	DB_PORT           string
	DB_USER           string
	DB_PASSWORD       string
	DB_NAME           string
	DB_SSLMODE        string
	PageSize          int
	Debug             bool
	Env               string
	MigrateURL        string
	OAuthClientID     string
	OAuthClientSecret string
}

func panicIfErr(val string) string {
	v, ok := os.LookupEnv(val)
	if !ok {
		panic(fmt.Sprintf("missing env var %s", val))
	}

	return v
}

func panicIfErrInt(val string) int {
	v, err := strconv.Atoi(panicIfErr(val))
	if err != nil {
		panic(fmt.Sprintf("invalid %s: %s", val, err))
	}
	return v
}

func LoadConfig() *AppEnv {

	_ = godotenv.Load()

	return &AppEnv{
		AppPort: panicIfErr("APP_PORT"),
		AppHost: panicIfErr("APP_HOST"),

		DB_HOST:           panicIfErr("DB_HOST"),
		DB_PORT:           panicIfErr("DB_PORT"),
		DB_USER:           panicIfErr("DB_USER"),
		DB_PASSWORD:       panicIfErr("DB_PASSWORD"),
		DB_NAME:           panicIfErr("DB_NAME"),
		DB_SSLMODE:        panicIfErr("DB_SSLMODE"),
		MigrateURL:        panicIfErr("MIGRATE_DATABASE_URL"),
		OAuthClientID:     panicIfErr("OAUTH_CLIENT_ID"),
		OAuthClientSecret: panicIfErr("OAUTH_CLIENT_SECRET"),

		PageSize: panicIfErrInt("PAGE_SIZE"),
		Debug:    panicIfErr("DEBUG") == "true",
		Env:      panicIfErr("ENV"),
	}
}
