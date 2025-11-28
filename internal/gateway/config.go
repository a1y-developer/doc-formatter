package gateway

type Config struct {
	Address     string
	AuthService string
}

func NewConfig() *Config {
	return &Config{}
}
