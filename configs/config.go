package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Service struct {
	URL string
}

type Config struct {
	Services []Service
}

func Load() Config {
	_ = godotenv.Load()

	raw := os.Getenv("SERVICES")

	var services []Service
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		services = append(services, Service{URL: p})
	}

	return Config{
		Services: services,
	}
}
