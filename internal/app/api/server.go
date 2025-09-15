package api

type server struct {
	config *Config
}

// New creates a new API server instance
func newServer(cfg *Config) *server {
	return &server{config: cfg}
}
