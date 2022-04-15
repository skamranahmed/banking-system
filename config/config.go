package config

import (
	"os"

	"github.com/spf13/viper"

	"log"
)

// AppEnvironment : string wrapper for environment name
type AppEnvironment string

func (e AppEnvironment) IsLocal() bool {
	return e == AppEnvironmentLocal
}

var (
	// Database Credentials
	DbDriver   string
	DbHost     string
	DbUser     string
	DbName     string
	DbPassword string
	DbPort     string

	// Test Database Credentials
	TestDbDriver   string
	TestDbHost     string
	TestDbUser     string
	TestDbName     string
	TestDbPassword string
	TestDbPort     string

	// Server
	ServerPort string

	// Environment
	Environment AppEnvironment

	// slice of all app environments except the `local`` env
	AppEnvironemnts = []AppEnvironment{
		AppEnvironmentProduction,
	}
)

const (
	// AppEnvironmentLocal : local env
	AppEnvironmentLocal = AppEnvironment("local")

	// AppEnvironmentLocal : production env
	AppEnvironmentProduction = AppEnvironment("production")

	// ConfigFileName : localConfig.yaml
	ConfigFileName string = "localConfig"

	// ConfigFileType : yaml
	ConfigFileType string = "yaml"
)

// func init() {
// 	SetConfigFromViper()
// }

func Load(path string) {
	SetConfigFromViper(path)
}

func SetConfigFromViper(path string) {
	Environment = getCurrentHostEnvironment()
	log.Printf("🚀 Current Host Environment: %s\n", Environment)

	// if env is local, we set the env variables using the config file
	if Environment.IsLocal() {
		setEnvironmentVarsFromConfig(path)
	}

	// fetch the env vars and store in variables
	// Database Credentials
	DbDriver = os.Getenv("DB_DRIVER")
	DbHost = os.Getenv("DB_HOST")
	DbUser = os.Getenv("DB_USER")
	DbName = os.Getenv("DB_NAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbPort = os.Getenv("DB_PORT")

	// Test Database Credentials
	TestDbDriver = os.Getenv("TEST_DB_DRIVER")
	TestDbHost = os.Getenv("TEST_DB_HOST")
	TestDbUser = os.Getenv("TEST_DB_USER")
	TestDbName = os.Getenv("TEST_DB_NAME")
	TestDbPassword = os.Getenv("TEST_DB_PASSWORD")
	TestDbPort = os.Getenv("TEST_DB_PORT")

	// Server
	ServerPort = os.Getenv("SERVER_PORT")
}

func setEnvironmentVarsFromConfig(path string) {
	// add the path of the config file
	viper.AddConfigPath(path)

	// set the config file name
	viper.SetConfigName(ConfigFileName)
	// set the config file type
	viper.SetConfigType(ConfigFileType)

	// read the env vars from the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("unable to read env vars from config file")
	}

	/*
		Step 1. Get the env vars from viper
		Step 2. Set the host OS env vars
	*/

	// Database Credentials
	dbDriver := viper.GetString("DB_DRIVER")
	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbName := viper.GetString("DB_NAME")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbPort := viper.GetString("DB_PORT")
	os.Setenv("DB_DRIVER", dbDriver)
	os.Setenv("DB_HOST", dbHost)
	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("DB_PORT", dbPort)

	// Test Database Credentials
	testDbDriver := viper.GetString("TEST_DB_DRIVER")
	testDbHost := viper.GetString("TEST_DB_HOST")
	testDbUser := viper.GetString("TEST_DB_USER")
	testDbName := viper.GetString("TEST_DB_NAME")
	testDbPassword := viper.GetString("TEST_DB_PASSWORD")
	testDbPort := viper.GetString("TEST_DB_PORT")
	os.Setenv("TEST_DB_DRIVER", testDbDriver)
	os.Setenv("TEST_DB_HOST", testDbHost)
	os.Setenv("TEST_DB_USER", testDbUser)
	os.Setenv("TEST_DB_NAME", testDbName)
	os.Setenv("TEST_DB_PASSWORD", testDbPassword)
	os.Setenv("TEST_DB_PORT", testDbPort)

	// Server
	serverPort := viper.GetString("SERVER_PORT")
	os.Setenv("SERVER_PORT", serverPort)
}

func getCurrentHostEnvironment() AppEnvironment {
	currentHostEnvironment := os.Getenv("ENVIRONMENT")
	for _, env := range AppEnvironemnts {
		if env == AppEnvironment(currentHostEnvironment) {
			return env
		}
	}
	// if env not found return `local`` env
	return AppEnvironmentLocal
}
