package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Server struct {
		Mode         string `env:"MODE"`
		Port         string `env:"PORT"`
		JWTSecretKey string `env:"JWT_SECRET_KEY"`
	}

	Client struct {
		Url  string `env:"FRONTEND_URL"`
		Port string `env:"CLIENT_PORT"`
	}

	Database struct {
		DBHost     string `env:"DB_HOST"`
		DBPort     string `env:"DB_PORT"`
		DBUser     string `env:"DB_USER"`
		DBPassword string `env:"DB_PASSWORD"`
		DBName     string `env:"DB_NAME"`
	}

	Config struct {
		Server   Server
		Client   Client
		Database Database
	}
)

var Cfg Config

func LoadNative() ([]string, error) {
	var ok bool
	var missing []string

	if _, err := os.Stat("/app/.env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	// ! Load Server configuration ____________________________
	if Cfg.Server.Mode, ok = os.LookupEnv("MODE"); !ok {
		missing = append(missing, "MODE env is not set")
	}
	if Cfg.Server.Port, ok = os.LookupEnv("PORT"); !ok {
		missing = append(missing, "PORT env is not set")
	}
	if Cfg.Server.JWTSecretKey, ok = os.LookupEnv("JWT_SECRET_KEY"); !ok {
		missing = append(missing, "JWT_SECRET_KEY env is not set")
	}
	// ! ______________________________________________________

	// ! Load Client configuration ____________________________
	if Cfg.Client.Url, ok = os.LookupEnv("FRONTEND_URL"); !ok {
		missing = append(missing, "FRONTEND_URL env is not set")
	}
	if Cfg.Client.Port, ok = os.LookupEnv("CLIENT_PORT"); !ok {
		missing = append(missing, "CLIENT_PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Database configuration __________________________
	if Cfg.Database.DBUser, ok = os.LookupEnv("DB_USER"); !ok {
		missing = append(missing, "DB_USER env is not set")
	}
	if Cfg.Database.DBHost, ok = os.LookupEnv("DB_HOST"); !ok {
		missing = append(missing, "DB_HOST env is not set")
	}
	if Cfg.Database.DBPort, ok = os.LookupEnv("DB_PORT"); !ok {
		missing = append(missing, "DB_PORT env is not set")
	}
	if Cfg.Database.DBName, ok = os.LookupEnv("DB_NAME"); !ok {
		missing = append(missing, "DB_NAME env is not set")
	}
	if Cfg.Database.DBPassword, ok = os.LookupEnv("DB_PASSWORD"); !ok {
		missing = append(missing, "DB_PASSWORD env is not set")
	}
	// ! ______________________________________________________

	return missing, nil
}

func LoadByViper() ([]string, error) {
	var missing []string

	config := viper.New()
	if configFile, err := os.Stat("/app/config.json"); err != nil || configFile.IsDir() {
		config.SetConfigFile("config.json")
	} else {
		config.SetConfigFile("/app/config.json")
	}

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	// ! Load Server configuration ____________________________
	if Cfg.Server.Mode = config.GetString("MODE"); Cfg.Server.Mode == "" {
		missing = append(missing, "MODE env is not set")
	}
	if Cfg.Server.Port = config.GetString("PORT"); Cfg.Server.Port == "" {
		missing = append(missing, "PORT env is not set")
	}
	if Cfg.Server.JWTSecretKey = config.GetString("JWT_SECRET_KEY"); Cfg.Server.JWTSecretKey == "" {
		missing = append(missing, "JWT_SECRET_KEY env is not set")
	}
	// ! ______________________________________________________

	// ! Load Client configuration ____________________________
	if Cfg.Client.Url = config.GetString("CLIENT.URL"); Cfg.Client.Url == "" {
		missing = append(missing, "CLIENT.URL env is not set")
	}
	if Cfg.Client.Port = config.GetString("CLIENT.PORT"); Cfg.Client.Port == "" {
		missing = append(missing, "CLIENT.PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Database configuration __________________________
	if Cfg.Database.DBUser = config.GetString("DATABASE.POSTGRESQL.USER"); Cfg.Database.DBUser == "" {
		missing = append(missing, "DATABASE.POSTGRESQL.USER env is not set")
	}
	if Cfg.Database.DBHost = config.GetString("DATABASE.POSTGRESQL.HOST"); Cfg.Database.DBHost == "" {
		missing = append(missing, "DATABASE.POSTGRESQL.HOST env is not set")
	}
	if Cfg.Database.DBPort = config.GetString("DATABASE.POSTGRESQL.PORT"); Cfg.Database.DBPort == "" {
		missing = append(missing, "DATABASE.POSTGRESQL.PORT env is not set")
	}
	if Cfg.Database.DBName = config.GetString("DATABASE.POSTGRESQL.NAME"); Cfg.Database.DBName == "" {
		missing = append(missing, "DATABASE.POSTGRESQL.NAME env is not set")
	}
	if Cfg.Database.DBPassword = config.GetString("DATABASE.POSTGRESQL.PASSWORD"); Cfg.Database.DBPassword == "" {
		missing = append(missing, "DATABASE.POSTGRESQL.PASSWORD env is not set")
	}
	// ! ______________________________________________________

	return missing, nil
}
