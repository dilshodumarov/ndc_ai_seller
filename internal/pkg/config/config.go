package config

import (
	"os"

	"github.com/spf13/cast"
)

type (
	// Config -.
	Config struct {
		App         App
		Environment string
		Server      Server
		Context     Context
		HTTP        HTTP
		GRPC        GRPC
		Log         Log
		PG          PG
		DB          DB
		RMQ         RMQ
		Redis       Redis
		Email       EmailConfig
		//OTLPCollector WebAddress
		JWT       JWT
		AppConfig AppConfig
	}

	Server struct {
		Host         string `json:"HOST"`
		Port         string `json:"PORT"`
		ReadTimeout  string `json:"readTimeout"`
		WriteTimeout string `json:"writeTimeout"`
		IdleTimeout  string `json:"idleTimeout"`
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Context struct {
		Timeout string
	}

	// GRPC -.
	GRPC struct {
		Port string `env:"GRPC_PORT" envDefault:"50051"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// EmailConfig -.
	EmailConfig struct {
		From     string `env:"EMAIL_FROM,required"`
		Password string `env:"EMAIL_PASSWORD,required"`
		Host     string `env:"EMAIL_HOST,required"`
		Port     string `env:"EMAIL_PORT,required"`
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		Name     string
	}

	WebAddress struct {
		Host string
		Port string
	}

	JWT struct {
		Secret string `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
	}
	AppConfig struct {
		ClientID     string
		ClientSecret string
		GrantType    string
		RedirectURI  string
	}
)

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.App.Name = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.Log.Level = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	config.Server.Port = getEnv("SERVER_PORT", ":8080")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "ai-seller")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "seller_db")
	config.DB.User = getEnv("POSTGRES_USER", "ai-seller")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "ai-seller_secret")
	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	config.PG.PoolMax = cast.ToInt(getEnv("POSTGRES_POOL_MAX", "1"))

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "ai-seller-redis")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	config.Email.From = getEnv("EMAIL_FROM", "the.aura.fashionn@gmail.com")
	config.Email.Password = getEnv("EMAIL_PASSWORD", "uqik mvwn qtfe wbte")
	config.Email.Port = getEnv("EMAIL_PORT", "587")
	config.Email.Host = getEnv("EMAIL_HOST", "smtp.gmail.com")

	// config.ContentService.Host = getEnv("CONTENT_SERVICE_GRPC_HOST", "localhost")
	// config.ContentService.Port = getEnv("CONTENT_SERVICE_GRPC_PORT", ":50025")

	// otlp collector configuration
	//config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "localhost")
	//config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")
	config.JWT.Secret = getEnv("JWT_SECRET", "jlakdjfadkjfl")
	// kafka configuration
	// config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:29092"), ",")
	// config.Kafka.Topic.InvestmentPaymentTransaction = getEnv("KAFKA_TOPIC_INVESTMENT_PAYMENT_TRANSACTION", "investment.payment.transaction")

	config.AppConfig.ClientID = getEnv("CLIENT_ID", "")
	config.AppConfig.ClientSecret = getEnv("CLIENT_SECRET", "")
	config.AppConfig.GrantType = getEnv("GRANT_TYPE", "authorization_code")
	config.AppConfig.RedirectURI = getEnv("REDIRECT_URI", "")
	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// package config

// import (
// 	"os"

// 	"github.com/spf13/cast"
// )

// type (
// 	// Config -.
// 	Config struct {
// 		App         App
// 		Environment string
// 		Server      Server
// 		Context     Context
// 		HTTP        HTTP
// 		GRPC        GRPC
// 		Log         Log
// 		PG          PG
// 		DB          DB
// 		RMQ         RMQ
// 		Redis       Redis
// 		Email       EmailConfig
// 		//OTLPCollector WebAddress
// 		JWT JWT
// 	}

// 	Server struct {
// 		Host         string `json:"HOST"`
// 		Port         string `json:"PORT"`
// 		ReadTimeout  string `json:"readTimeout"`
// 		WriteTimeout string `json:"writeTimeout"`
// 		IdleTimeout  string `json:"idleTimeout"`
// 	}

// 	// App -.
// 	App struct {
// 		Name    string `env:"APP_NAME,required"`
// 		Version string `env:"APP_VERSION,required"`
// 	}

// 	// HTTP -.
// 	HTTP struct {
// 		Port string `env:"HTTP_PORT,required"`
// 	}

// 	Context struct {
// 		Timeout string
// 	}

// 	// GRPC -.
// 	GRPC struct {
// 		Port string `env:"GRPC_PORT" envDefault:"50051"`
// 	}

// 	// Log -.
// 	Log struct {
// 		Level string `env:"LOG_LEVEL,required"`
// 	}

// 	// PG -.
// 	PG struct {
// 		PoolMax int    `env:"PG_POOL_MAX,required"`
// 		URL     string `env:"PG_URL,required"`
// 	}

// 	DB struct {
// 		Host     string
// 		Port     string
// 		Name     string
// 		User     string
// 		Password string
// 		SSLMode  string
// 	}

// 	// RMQ -.
// 	RMQ struct {
// 		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
// 		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
// 		URL            string `env:"RMQ_URL,required"`
// 	}

// 	// EmailConfig -.
// 	EmailConfig struct {
// 		From     string `env:"EMAIL_FROM,required"`
// 		Password string `env:"EMAIL_PASSWORD,required"`
// 		Host     string `env:"EMAIL_HOST,required"`
// 		Port     string `env:"EMAIL_PORT,required"`
// 	}

// 	Redis struct {
// 		Host     string
// 		Port     string
// 		Password string
// 		Name     string
// 	}

// 	WebAddress struct {
// 		Host string
// 		Port string
// 	}

// 	JWT struct {
// 		Secret string `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
// 	}
// )

// func NewConfig() (*Config, error) {
// 	var config Config

// 	// general configuration
// 	config.App.Name = getEnv("APP", "app")
// 	config.Environment = getEnv("ENVIRONMENT", "develop")
// 	config.Log.Level = getEnv("LOG_LEVEL", "debug")
// 	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

// 	// server configuration
// 	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
// 	config.Server.Port = getEnv("SERVER_PORT", ":8080")
// 	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
// 	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
// 	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

// 	// db configuration

// 	config.DB.Host = getEnv("POSTGRES_HOST", "108.181.201.147")
// 	config.DB.Port = getEnv("POSTGRES_PORT", "5435")
// 	config.DB.Name = getEnv("POSTGRES_DATABASE", "seller_db")
// 	config.DB.User = getEnv("POSTGRES_USER", "ai-seller")
// 	config.DB.Password = getEnv("POSTGRES_PASSWORD", "ai-seller_secret")
// 	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

// 	config.PG.PoolMax = cast.ToInt(getEnv("POSTGRES_POOL_MAX", "1"))

// 	// redis configuration
// 	config.Redis.Host = getEnv("REDIS_HOST", "108.181.201.147")
// 	config.Redis.Port = getEnv("REDIS_PORT", "6378")
// 	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
// 	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

// 	config.Email.From = getEnv("EMAIL_FROM", "the.aura.fashionn@gmail.com")
// 	config.Email.Password = getEnv("EMAIL_PASSWORD", "uqik mvwn qtfe wbte")
// 	config.Email.Port = getEnv("EMAIL_PORT", "587")
// 	config.Email.Host = getEnv("EMAIL_HOST", "smtp.gmail.com")

// 	// config.ContentService.Host = getEnv("CONTENT_SERVICE_GRPC_HOST", "localhost")
// 	// config.ContentService.Port = getEnv("CONTENT_SERVICE_GRPC_PORT", ":50025")

// 	// otlp collector configuration
// 	//config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "localhost")
// 	//config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")
// 	config.JWT.Secret = getEnv("JWT_SECRET", "jlakdjfadkjfl")
// 	// kafka configuration
// 	// config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "localhost:29092"), ",")
// 	// config.Kafka.Topic.InvestmentPaymentTransaction = getEnv("KAFKA_TOPIC_INVESTMENT_PAYMENT_TRANSACTION", "investment.payment.transaction")

// 	return &config, nil
// }

// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}

// 	return defaultValue
// }
