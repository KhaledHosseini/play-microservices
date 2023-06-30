package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	APP_VERSION = "APP_VERSION"

	ENVIRONMENT = "ENVIRONMENT"
	TOPICS_FILE = "TOPICS_FILE"

	SERVER_PORT = "SERVER_PORT"

	DATABASE_USER_FILE    = "DATABASE_USER_FILE"
	DATABASE_PASS_FILE    = "DATABASE_PASS_FILE"
	DATABASE_DB_NAME_FILE = "DATABASE_DB_NAME_FILE"
	DATABASE_SCHEMA       = "DATABASE_SCHEMA"
	DATABASE_HOST_NAME    = "DATABASE_HOST_NAME"
	DATABASE_PORT         = "DATABASE_PORT"

	KAFKA_BROKERS               = "KAFKA_BROKERS"
	TOPIC_JOB_CREATE            = "TOPIC_JOB_CREATE"
	TOPIC_JOB_CREATE_PARTITIONS = "TOPIC_JOB_CREATE_PARTITIONS"
	TOPIC_JOB_CREATE_REPLICAS   = "TOPIC_JOB_CREATE_REPLICAS"
	TOPIC_JOB_UPDATE            = "TOPIC_JOB_UPDATE"
	TOPIC_JOB_UPDATE_PARTITIONS = "TOPIC_JOB_UPDATE_PARTITIONS"
	TOPIC_JOB_UPDATE_REPLICAS   = "TOPIC_JOB_UPDATE_REPLICAS"
	TOPIC_JOB_RUN               = "TOPIC_JOB_RUN"
	TOPIC_JOB_RUN_PARTITIONS    = "TOPIC_JOB_RUN_PARTITIONS"
	TOPIC_JOB_RUN_REPLICAS      = "TOPIC_JOB_RUN_REPLICAS"
)

type Config struct {
	AppVersion  string
	Environment string

	ServerPort string

	DatabaseUser     string
	DatabasePass     string
	DatabaseDBName   string
	DatabaseSchema   string
	DatabaseHostName string
	DatabasePort     string
	DatabaseURI      string

	KafkaBrokers             []string
	TopicJobCreate           string
	TopicJobCreatePartitions int
	TopicJobCreateReplicas   int
	TopicJobUpdate           string
	TopicJobUpdatePartitions int
	TopicJobUpdateReplicas   int
	TopicJobRun              string
	TopicJobRunPartitions    int
	TopicJobRunReplicas      int

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

	topicsFile := os.Getenv(TOPICS_FILE)
	if topicsFile == "" {
		topicsFile = "./config/.env.topics"
	}
	err2 := godotenv.Load(topicsFile)
	if err2 != nil {
		log.Fatalf("Error loading topics file")
	}

	var c Config
	c.AppVersion = goDotEnvVariable(APP_VERSION)
	c.Environment = goDotEnvVariable(ENVIRONMENT)
	c.ServerPort = goDotEnvVariable(SERVER_PORT)

	db_user_file := goDotEnvVariable(DATABASE_USER_FILE)
	c.DatabaseUser = fileContent(db_user_file)
	db_pass_file := goDotEnvVariable(DATABASE_PASS_FILE)
	c.DatabasePass = fileContent(db_pass_file)
	db_name_file := goDotEnvVariable(DATABASE_DB_NAME_FILE)
	c.DatabaseDBName = fileContent(db_name_file)
	c.DatabaseSchema = goDotEnvVariable(DATABASE_SCHEMA)
	c.DatabaseHostName = goDotEnvVariable(DATABASE_HOST_NAME)
	c.DatabasePort = goDotEnvVariable(DATABASE_PORT)
	c.DatabaseURI = c.DatabaseSchema + "://" + c.DatabaseHostName + ":" + c.DatabasePort

	c.KafkaBrokers = strings.Split(goDotEnvVariable(KAFKA_BROKERS), ",")
	c.TopicJobCreate = goDotEnvVariable(TOPIC_JOB_CREATE)
	c.TopicJobCreatePartitions = goDotEnvVariableInteger(TOPIC_JOB_CREATE_PARTITIONS)
	c.TopicJobCreateReplicas = goDotEnvVariableInteger(TOPIC_JOB_CREATE_REPLICAS)
	c.TopicJobUpdate = goDotEnvVariable(TOPIC_JOB_UPDATE)
	c.TopicJobUpdatePartitions = goDotEnvVariableInteger(TOPIC_JOB_UPDATE_PARTITIONS)
	c.TopicJobUpdateReplicas = goDotEnvVariableInteger(TOPIC_JOB_UPDATE_REPLICAS)
	c.TopicJobRun = goDotEnvVariable(TOPIC_JOB_RUN)
	c.TopicJobRunPartitions = goDotEnvVariableInteger(TOPIC_JOB_RUN_PARTITIONS)
	c.TopicJobRunReplicas = goDotEnvVariableInteger(TOPIC_JOB_RUN_REPLICAS)

	c.Logger_DisableCaller = false
	c.Logger_DisableStacktrace = false
	c.Logger_Encoding = "json"
	c.Logger_Level = "info"

	return &c, nil
}
