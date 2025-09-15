package api

type Config struct {
	Port int
}

func NewConfig() Config {
	return Config{Port: 8080}
}
