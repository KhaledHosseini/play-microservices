package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	APP_VERSION = "APP_VERSION"

	ENVIRONMENT = "ENVIRONMENT"
	SERVER_PORT = "SERVER_PORT"

	AUTH_SERVICE_URL      = "AUTH_SERVICE_URL"
	SCHEDULER_SERVICE_URL = "SCHEDULER_SERVICE_URL"
	REPORT_SERVICE_URL    = "SCHEDULER_SERVICE_URL"
)

type Config struct {
	AppVersion  string
	Environment string

	ServerPort string

	AuthServiceURL      string
	SchedulerServiceURL string
	ReportServiceURL    string

	Logger_DisableCaller     bool
	Logger_DisableStacktrace bool
	Logger_Encoding          string
	Logger_Level             string
}

func goDotEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error loading enironment variable for: %s.", key)
	}
	return value
}

func goDotEnvVariableInteger(key string) int {

	str := goDotEnvVariable(key)
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Error loading integer enironment variable for: %s.", key)
	}
	return num
}

func fileContent(file string) string {

	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	return text
}

func (cfg *Config) IsDevelopmentEnvironment() bool {
	return cfg.Environment == "development"
}

func InitConfig() (*Config, error) {

	// load .env file
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var c Config
	c.AppVersion = goDotEnvVariable(APP_VERSION)
	c.Environment = goDotEnvVariable(ENVIRONMENT)
	c.ServerPort = goDotEnvVariable(SERVER_PORT)
	c.AuthServiceURL = goDotEnvVariable(AUTH_SERVICE_URL)
	c.SchedulerServiceURL = goDotEnvVariable(SCHEDULER_SERVICE_URL)
	c.ReportServiceURL = goDotEnvVariable(REPORT_SERVICE_URL)

	c.Logger_DisableCaller = false
	c.Logger_DisableStacktrace = false
	c.Logger_Encoding = "json"
	c.Logger_Level = "info"

	return &c, nil
}
