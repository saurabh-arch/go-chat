package main

type Config struct {
	APIPort int
}

func LoadConfig() Config {
	var cfg Config

	cfg.APIPort = 6000

	return cfg
}
