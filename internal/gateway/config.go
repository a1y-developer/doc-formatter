package gateway

type Config struct {
	HTTPAddr     string
	AuthGRPCAddr string
}

func NewConfig() *Config {
	return &Config{}
}
