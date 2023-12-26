package app

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr  string
	ServerPort uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddr:  "localhost:6379",
		ServerPort: 8000,
	}

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.RedisAddr = redisAddr
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}
	return cfg
}
