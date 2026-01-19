package env

import (
	"fmt"
	"os"
	"strconv"

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

	Redis struct {
		RHost string `env:"REDIS_HOST"`
		RPort string `env:"REDIS_PORT"`
	}

	GoogleOAuth struct {
		GOClientID     string `env:"GOOGLE_CLIENT_ID"`
		GOClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	}

	GithubOAuth struct {
		GHClientID     string `env:"GITHUB_CLIENT_ID"`
		GHClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	}

	MicrosoftOAuth struct {
		MSClientID       string `env:"MICROSOFT_CLIENT_ID"`
		MSClientSecret   string `env:"MICROSOFT_CLIENT_SECRET"`
		MSTenantID       string `env:"MICROSOFT_TENANT_ID"`
		MSClientSecretID string `env:"MICROSOFT_CLIENT_SECRET_ID"`
	}

	OAuth struct {
		Google    GoogleOAuth
		Github    GithubOAuth
		Microsoft MicrosoftOAuth
	}

	GSMTP struct {
		GSHost     string `env:"GOOGLE_SMTP_HOST"`
		GSPort     string `env:"GOOGLE_SMTP_PORT"`
		GSUser     string `env:"GOOGLE_SMTP_USER"`
		GSPassword string `env:"GOOGLE_SMTP_PASSWORD"`
	}

	RabbitMQ struct {
		RMQHost        string `env:"RABBITMQ_HOST"`
		RMQPort        string `env:"RABBITMQ_PORT"`
		RMQUser        string `env:"RABBITMQ_USER"`
		RMQPassword    string `env:"RABBITMQ_PASSWORD"`
		RMQVirtualHost string `env:"RABBITMQ_VIRTUAL_HOST"`
	}

	ZSMTP struct {
		ZSHost     string `env:"ZOHO_SMTP_HOST"`
		ZSPort     string `env:"ZOHO_SMTP_PORT"`
		ZSUser     string `env:"ZOHO_SMTP_USER"`
		ZSPassword string `env:"ZOHO_SMTP_PASSWORD"`
		ZSSecure   string `env:"ZOHO_SMTP_SECURE"`
		ZSAuth     bool   `env:"ZOHO_SMTP_AUTH"`
	}

	Minio struct {
		Host        string `env:"MINIO_HOST"`
		AccessKey   string `env:"MINIO_ROOT_USER"`
		SecretKey   string `env:"MINIO_ROOT_PASSWORD"`
		MaxOpenConn int    `env:"MINIO_MAX_OPEN_CONN"`
		UseSSL      int    `env:"MINIO_USE_SSL"`
	}

	Config struct {
		Server   Server
		Client   Client
		Database Database
		Redis    Redis
		OAuth    OAuth
		GSMTP    GSMTP
		ZSMTP    ZSMTP
		RabbitMQ RabbitMQ
		Minio    Minio
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

	// ! Load Redis configuration _____________________________
	if Cfg.Redis.RHost, ok = os.LookupEnv("REDIS_HOST"); !ok {
		missing = append(missing, "REDIS_HOST env is not set")
	}
	if Cfg.Redis.RPort, ok = os.LookupEnv("REDIS_PORT"); !ok {
		missing = append(missing, "REDIS_PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load Google OAuth configuration ______________________
	if Cfg.OAuth.Google.GOClientID, ok = os.LookupEnv("GOOGLE_CLIENT_ID"); !ok {
		missing = append(missing, "GOOGLE_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Google.GOClientSecret, ok = os.LookupEnv("GOOGLE_CLIENT_SECRET"); !ok {
		missing = append(missing, "GOOGLE_CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Github OAuth configuration ______________________
	if Cfg.OAuth.Github.GHClientID, ok = os.LookupEnv("GITHUB_CLIENT_ID"); !ok {
		missing = append(missing, "GITHUB_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Github.GHClientSecret, ok = os.LookupEnv("GITHUB_CLIENT_SECRET"); !ok {
		missing = append(missing, "GITHUB_CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Microsoft OAuth configuration ___________________
	if Cfg.OAuth.Microsoft.MSClientID, ok = os.LookupEnv("MICROSOFT_CLIENT_ID"); !ok {
		missing = append(missing, "MICROSOFT_CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecret, ok = os.LookupEnv("MICROSOFT_CLIENT_SECRET"); !ok {
		missing = append(missing, "MICROSOFT_CLIENT_SECRET env is not set")
	}
	if Cfg.OAuth.Microsoft.MSTenantID, ok = os.LookupEnv("MICROSOFT_TENANT_ID"); !ok {
		missing = append(missing, "MICROSOFT_TENANT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecretID, ok = os.LookupEnv("MICROSOFT_CLIENT_SECRET_ID"); !ok {
		missing = append(missing, "MICROSOFT_CLIENT_SECRET_ID env is not set")
	}
	// ! ______________________________________________________

	// ! Load Gmail SMTP configuration ______________________________
	if Cfg.GSMTP.GSHost, ok = os.LookupEnv("GOOGLE_SMTP_HOST"); !ok {
		missing = append(missing, "GOOGLE_SMTP_HOST env is not set")
	}
	if Cfg.GSMTP.GSPort, ok = os.LookupEnv("GOOGLE_SMTP_PORT"); !ok {
		missing = append(missing, "GOOGLE_SMTP_PORT env is not set")
	}
	if Cfg.GSMTP.GSUser, ok = os.LookupEnv("GOOGLE_SMTP_USER"); !ok {
		missing = append(missing, "GOOGLE_SMTP_USER env is not set")
	}
	if Cfg.GSMTP.GSPassword, ok = os.LookupEnv("GOOGLE_SMTP_PASSWORD"); !ok {
		missing = append(missing, "GOOGLE_SMTP_PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load RabbitMQ configuration __________________________
	if Cfg.RabbitMQ.RMQUser, ok = os.LookupEnv("RABBITMQ_USER"); !ok {
		missing = append(missing, "RABBITMQ_USER env is not set")
	}
	if Cfg.RabbitMQ.RMQPassword, ok = os.LookupEnv("RABBITMQ_PASSWORD"); !ok {
		missing = append(missing, "RABBITMQ_PASSWORD env is not set")
	}
	if Cfg.RabbitMQ.RMQHost, ok = os.LookupEnv("RABBITMQ_HOST"); !ok {
		missing = append(missing, "RABBITMQ_HOST env is not set")
	}
	if Cfg.RabbitMQ.RMQPort, ok = os.LookupEnv("RABBITMQ_PORT"); !ok {
		missing = append(missing, "RABBITMQ_PORT env is not set")
	}
	if Cfg.RabbitMQ.RMQVirtualHost, ok = os.LookupEnv("RABBITMQ_VIRTUAL_HOST"); !ok {
		missing = append(missing, "RABBITMQ_VIRTUAL_HOST env is not set")
	}
	// ! ______________________________________________________

	// ! Load Zoho SMTP configuration __________________________
	if Cfg.ZSMTP.ZSHost, ok = os.LookupEnv("ZOHO_SMTP_HOST"); !ok {
		missing = append(missing, "ZOHO_SMTP_HOST env is not set")
	}
	if Cfg.ZSMTP.ZSPort, ok = os.LookupEnv("ZOHO_SMTP_PORT"); !ok {
		missing = append(missing, "ZOHO_SMTP_PORT env is not set")
	}
	if Cfg.ZSMTP.ZSUser, ok = os.LookupEnv("ZOHO_SMTP_USER"); !ok {
		missing = append(missing, "ZOHO_SMTP_USER env is not set")
	}
	if Cfg.ZSMTP.ZSPassword, ok = os.LookupEnv("ZOHO_SMTP_PASSWORD"); !ok {
		missing = append(missing, "ZOHO_SMTP_PASSWORD env is not set")
	}
	if Cfg.ZSMTP.ZSSecure, ok = os.LookupEnv("ZOHO_SMTP_SECURE"); !ok {
		missing = append(missing, "ZOHO_SMTP_SECURE env is not set")
	}
	if zohoAuth, ok := os.LookupEnv("ZOHO_SMTP_AUTH"); !ok {
		missing = append(missing, "ZOHO_SMTP_AUTH env is not set")
	} else {
		Cfg.ZSMTP.ZSAuth = zohoAuth == "true"
	}
	// ! ______________________________________________________

	// ! Load MinIO configuration _____________________________
	if Cfg.Minio.Host, ok = os.LookupEnv("MINIO_HOST"); !ok {
		missing = append(missing, "MINIO_HOST env is not set")
	}
	if Cfg.Minio.AccessKey, ok = os.LookupEnv("MINIO_ROOT_USER"); !ok {
		missing = append(missing, "MINIO_ROOT_USER env is not set")
	}
	if Cfg.Minio.SecretKey, ok = os.LookupEnv("MINIO_ROOT_PASSWORD"); !ok {
		missing = append(missing, "MINIO_ROOT_PASSWORD env is not set")
	}
	if val, ok := os.LookupEnv("MINIO_MAX_OPEN_CONN"); !ok {
		missing = append(missing, "MINIO_MAX_OPEN_CONN env is not set")
	} else {
		var err error
		if Cfg.Minio.MaxOpenConn, err = strconv.Atoi(val); err != nil {
			missing = append(missing, fmt.Sprintf("MINIO_MAX_OPEN_CONN must be int, got %s", val))
		}
	}
	if val, ok := os.LookupEnv("MINIO_USE_SSL"); !ok {
		missing = append(missing, "MINIO_USE_SSL env is not set")
	} else {
		var err error
		if Cfg.Minio.UseSSL, err = strconv.Atoi(val); err != nil {
			missing = append(missing, fmt.Sprintf("MINIO_USE_SSL must be int, got %s", val))
		}
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

	// ! Load Redis configuration _____________________________
	if Cfg.Redis.RHost = config.GetString("REDIS.HOST"); Cfg.Redis.RHost == "" {
		missing = append(missing, "REDIS.HOST env is not set")
	}
	if Cfg.Redis.RPort = config.GetString("REDIS.PORT"); Cfg.Redis.RPort == "" {
		missing = append(missing, "REDIS.PORT env is not set")
	}
	// ! ______________________________________________________

	// ! Load RabbitMQ configuration __________________________
	if Cfg.RabbitMQ.RMQUser = config.GetString("MESSAGE-BROKER.RABBITMQ.USER"); Cfg.RabbitMQ.RMQUser == "" {
		missing = append(missing, "MESSAGE-BROKER.RABBITMQ.USER env is not set")
	}
	if Cfg.RabbitMQ.RMQPassword = config.GetString("MESSAGE-BROKER.RABBITMQ.PASSWORD"); Cfg.RabbitMQ.RMQPassword == "" {
		missing = append(missing, "MESSAGE-BROKER.RABBITMQ.PASSWORD env is not set")
	}
	if Cfg.RabbitMQ.RMQHost = config.GetString("MESSAGE-BROKER.RABBITMQ.HOST"); Cfg.RabbitMQ.RMQHost == "" {
		missing = append(missing, "MESSAGE-BROKER.RABBITMQ.HOST env is not set")
	}
	if Cfg.RabbitMQ.RMQPort = config.GetString("MESSAGE-BROKER.RABBITMQ.PORT"); Cfg.RabbitMQ.RMQPort == "" {
		missing = append(missing, "MESSAGE-BROKER.RABBITMQ.PORT env is not set")
	}
	if Cfg.RabbitMQ.RMQVirtualHost = config.GetString("MESSAGE-BROKER.RABBITMQ.VIRTUAL_HOST"); Cfg.RabbitMQ.RMQVirtualHost == "" {
		missing = append(missing, "MESSAGE-BROKER.RABBITMQ.VIRTUAL_HOST env is not set")
	}
	// ! ______________________________________________________

	// ! Load Google OAuth configuration ______________________
	if Cfg.OAuth.Google.GOClientID = config.GetString("OAUTH.GOOGLE.CLIENT_ID"); Cfg.OAuth.Google.GOClientID == "" {
		missing = append(missing, "OAUTH.GOOGLE.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Google.GOClientSecret = config.GetString("OAUTH.GOOGLE.CLIENT_SECRET"); Cfg.OAuth.Google.GOClientSecret == "" {
		missing = append(missing, "OAUTH.GOOGLE.CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Github OAuth configuration ______________________
	if Cfg.OAuth.Github.GHClientID = config.GetString("OAUTH.GITHUB.CLIENT_ID"); Cfg.OAuth.Github.GHClientID == "" {
		missing = append(missing, "OAUTH.GITHUB.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Github.GHClientSecret = config.GetString("OAUTH.GITHUB.CLIENT_SECRET"); Cfg.OAuth.Github.GHClientSecret == "" {
		missing = append(missing, "OAUTH.GITHUB.CLIENT_SECRET env is not set")
	}
	// ! ______________________________________________________

	// ! Load Microsoft OAuth configuration ___________________
	if Cfg.OAuth.Microsoft.MSClientID = config.GetString("OAUTH.MICROSOFT.CLIENT_ID"); Cfg.OAuth.Microsoft.MSClientID == "" {
		missing = append(missing, "OAUTH.MICROSOFT.CLIENT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecret = config.GetString("OAUTH.MICROSOFT.CLIENT_SECRET"); Cfg.OAuth.Microsoft.MSClientSecret == "" {
		missing = append(missing, "OAUTH.MICROSOFT.CLIENT_SECRET env is not set")
	}
	if Cfg.OAuth.Microsoft.MSTenantID = config.GetString("OAUTH.MICROSOFT.TENANT_ID"); Cfg.OAuth.Microsoft.MSTenantID == "" {
		missing = append(missing, "OAUTH.MICROSOFT.TENANT_ID env is not set")
	}
	if Cfg.OAuth.Microsoft.MSClientSecretID = config.GetString("OAUTH.MICROSOFT.CLIENT_SECRET_ID"); Cfg.OAuth.Microsoft.MSClientSecretID == "" {
		missing = append(missing, "OAUTH.MICROSOFT.CLIENT_SECRET_ID env is not set")
	}
	// ! ______________________________________________________

	// ! Load Gmail SMTP configuration ______________________________
	if Cfg.GSMTP.GSHost = config.GetString("SMTP.GOOGLE.HOST"); Cfg.GSMTP.GSHost == "" {
		missing = append(missing, "SMTP.GOOGLE.HOST env is not set")
	}
	if Cfg.GSMTP.GSPort = config.GetString("SMTP.GOOGLE.PORT"); Cfg.GSMTP.GSPort == "" {
		missing = append(missing, "SMTP.GOOGLE.PORT env is not set")
	}
	if Cfg.GSMTP.GSUser = config.GetString("SMTP.GOOGLE.USER"); Cfg.GSMTP.GSUser == "" {
		missing = append(missing, "SMTP.GOOGLE.USER env is not set")
	}
	if Cfg.GSMTP.GSPassword = config.GetString("SMTP.GOOGLE.PASSWORD"); Cfg.GSMTP.GSPassword == "" {
		missing = append(missing, "SMTP.GOOGLE.PASSWORD env is not set")
	}
	// ! ______________________________________________________

	// ! Load Zoho SMTP configuration __________________________
	if Cfg.ZSMTP.ZSHost = config.GetString("SMTP.ZOHO.HOST"); Cfg.ZSMTP.ZSHost == "" {
		missing = append(missing, "SMTP.ZOHO.HOST env is not set")
	}
	if Cfg.ZSMTP.ZSPort = config.GetString("SMTP.ZOHO.PORT"); Cfg.ZSMTP.ZSPort == "" {
		missing = append(missing, "SMTP.ZOHO.PORT env is not set")
	}
	if Cfg.ZSMTP.ZSUser = config.GetString("SMTP.ZOHO.USER"); Cfg.ZSMTP.ZSUser == "" {
		missing = append(missing, "SMTP.ZOHO.USER env is not set")
	}
	if Cfg.ZSMTP.ZSPassword = config.GetString("SMTP.ZOHO.PASSWORD"); Cfg.ZSMTP.ZSPassword == "" {
		missing = append(missing, "SMTP.ZOHO.PASSWORD env is not set")
	}
	if Cfg.ZSMTP.ZSSecure = config.GetString("SMTP.ZOHO.SECURE"); Cfg.ZSMTP.ZSSecure == "" {
		missing = append(missing, "SMTP.ZOHO.SECURE env is not set")
	}
	if Cfg.ZSMTP.ZSAuth = config.GetBool("SMTP.ZOHO.AUTH"); !Cfg.ZSMTP.ZSAuth {
		missing = append(missing, "SMTP.ZOHO.AUTH env is not set")
	}
	// ! ______________________________________________________

	// ! Load Minio configuration __________________________
	if Cfg.Minio.Host = config.GetString("OBJECT-STORAGE.MINIO.HOST"); Cfg.Minio.Host == "" {
		missing = append(missing, "OBJECT-STORAGE.MINIO.HOST env is not set")
	}
	if Cfg.Minio.AccessKey = config.GetString("OBJECT-STORAGE.MINIO.USER"); Cfg.Minio.AccessKey == "" {
		missing = append(missing, "OBJECT-STORAGE.MINIO.USER env is not set")
	}
	if Cfg.Minio.SecretKey = config.GetString("OBJECT-STORAGE.MINIO.PASSWORD"); Cfg.Minio.SecretKey == "" {
		missing = append(missing, "OBJECT-STORAGE.MINIO.PASSWORD env is not set")
	}
	if Cfg.Minio.MaxOpenConn = config.GetInt("OBJECT-STORAGE.MINIO.MAX_OPEN_CONN_POOL"); Cfg.Minio.MaxOpenConn == 0 {
		missing = append(missing, "OBJECT-STORAGE.MINIO.MAX_OPEN_CONN_POOL env is not set")
	}
	if Cfg.Minio.UseSSL = config.GetInt("OBJECT-STORAGE.MINIO.USE_SSL"); Cfg.Minio.UseSSL < 0 || Cfg.Minio.UseSSL > 1 {
		missing = append(missing, "OBJECT-STORAGE.MINIO.USE_SSL env is not valid")
	}
	
	// ! ______________________________________________________

	return missing, nil
}
