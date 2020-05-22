package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	PostgresURI string `toml:"postgres_uri"`
	RedisURL    string `toml:"redis_url"`
	Server      struct {
		Proto     string `toml:"proto"`
		Domain    string `toml:"domain"`
		APIDomain string `toml:"api_domain"`
		Port      int    `toml:"port"`
	} `toml:"server"`
	Mail struct {
		PrivateKey string `toml:"private_key"`
		PublicKey  string `toml:"public_key"`
		Domain     string `toml:"domain"`
	} `toml:"mail"`
	RateLimit struct {
		HTTPHeader     string `toml:"http_header"`
		QueryLimit     int    `toml:"query_limit"`
		QueryResetTime int    `toml:"query_reset_time"`
		MutLimit       int    `toml:"mut_limit"`
		MutResetTime   int    `toml:"mut_reset_time"`
		Cost           struct {
			CreateUser int `toml:"create_user"`
			PubThread  int `toml:"pub_thread"`
			PubPost    int `toml:"pub_post"`
		} `toml:"cost"`
	} `toml:"rate_limit"`
}

var c = Config{}

func setDefault() {
	c = Config{}
	c.Server.Domain = "abyss.club"
	c.Server.Domain = "api.abyss.club"
	c.Server.Proto = "http"
	c.Server.Port = 8000
}

func patchEnv() {
	c.PostgresURI = getenv("PG_URI", c.PostgresURI)
	c.RedisURL = getenv("REDIS_URI", c.PostgresURI)
	c.Server.Domain = getenv("DOMAIN", c.Server.Domain)
	c.Server.APIDomain = getenv("API_DOMAIN", c.Server.APIDomain)
	c.Server.Proto = getenv("PROTO", c.Server.Proto)
	c.Server.Port = getenvInt("PORT", c.Server.Port)
	c.Mail.PrivateKey = getenv("MAILGUN_PRIVATE_KEY", c.Mail.PrivateKey)
	c.Mail.PublicKey = getenv("MAILGUN_PUBLIC_KEY", c.Mail.PublicKey)
	c.Mail.Domain = getenv("MAILGUN_DOMAIN", c.Mail.Domain)
}

func Load(filename string) error {
	setDefault()
	if filename == "" {
		return nil
	}
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	patchEnv()
	return nil
}

func Get() *Config {
	return &c
}

func getenv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func getenvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}
