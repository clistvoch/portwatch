package config

import "fmt"

type RedisConfig struct {
	Enabled   bool   `toml:"enabled"`
	Address   string `toml:"address"`
	Password  string `toml:"password"`
	DB        int    `toml:"db"`
	Topic     string `toml:"topic"`
	TLSEnable bool   `toml:"tls_enable"`
}

func defaultRedisConfig() RedisConfig {
	return RedisConfig{
		Enabled:   false,
		Address:   "localhost:6379",
		Password:  "",
		DB:        0,
		Topic:     "portwatch:changes",
		TLSEnable: false,
	}
}

func validateRedis(cfg RedisConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Address == "" {
		return ValidationError{Field: "redis.address", Msg: "address is required when redis is enabled"}
	}
	if cfg.Topic == "" {
		return fmt.Errorf("redis.topic: must not be empty")
	}
	if cfg.DB < 0 || cfg.DB > 15 {
		return ValidationError{Field: "redis.db", Msg: "db must be between 0 and 15"}
	}
	return nil
}
