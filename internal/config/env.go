package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppEnv struct {
	AppPort     string
	AppHost     string
	DB_NAME     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	PageSize    int
	Debug       bool
	Env         string
	MigrateURL  string
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

		DB_NAME:     panicIfErr("DB_NAME"),
		DB_PORT:     panicIfErr("DB_PORT"),
		DB_USER:     panicIfErr("DB_USER"),
		DB_PASSWORD: panicIfErr("DB_PASSWORD"),
		MigrateURL:  panicIfErr("MIGRATE_DATABASE_URL"),

		PageSize: panicIfErrInt("PAGE_SIZE"),
		Debug:    panicIfErr("DEBUG") == "true",
		Env:      panicIfErr("ENV"),
	}
}
