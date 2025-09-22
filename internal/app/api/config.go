package api

// Config holds the configuration settings for the API server.
type Config struct {
	ServerAddress      string `yaml:"server_address"`
	LogLevel           string `yaml:"log_level"`
	DbConnectionString string `yaml:"db_connection_string"`
}

// DefaultConfig provides default configuration values for the API server.
func NewConfig() *Config {
	return &Config{
		ServerAddress: ":8080",
	}
}

func (c *Config) GetConnectionString() string {
	return c.DbConnectionString
}
