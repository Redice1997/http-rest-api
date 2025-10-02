package api

// Config holds the configuration settings for the API server.
type Config struct {
	ServerAddress      string `yaml:"server_address"`
	LogLevel           string `yaml:"log_level"`
	DbConnectionString string `yaml:"db_connection_string"`
	SessionKey         string `yaml:"session_key"`
}

// NewConfig provides default configuration values for the API server.
func NewConfig() *Config {
	return &Config{
		ServerAddress:      ":8080",
		DbConnectionString: "host=localhost port=5432 user=api password=password dbname=test_api_db sslmode=disable",
		LogLevel:           "debug",
		SessionKey:         "secret_key",
	}
}
