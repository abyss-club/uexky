package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/pelletier/go-toml"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

type Config struct {
	Env            RuntimeEnv `toml:"env"`
	PostgresURI    string     `toml:"postgres_uri"`
	RedisURI       string     `toml:"redis_uri"`
	MigrationFiles string     `toml:"migration_files"`
	Server         struct {
		Proto     string `toml:"proto"`
		Domain    string `toml:"domain"`
		APIDomain string `toml:"api_domain"`
		Port      int    `toml:"port"`
		Host      string `toml:"host"`
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

	filename string `toml:"-"`
}

var c = Config{}

func setDefault() {
	c = Config{}
	c.PostgresURI = "postgres://postgres:postgres@localhost:5432/uexky2?sslmode=disable"
	c.RedisURI = "redis://localhost:6379/0"
	c.MigrationFiles = "./migrations"
	c.Server.Domain = "abyss.club"
	c.Server.Domain = "api.abyss.club"
	c.Server.Proto = "http"
	c.Server.Port = 8000
	c.Server.Host = "localhost"
}

func patchEnv() {
	c.Env = RuntimeEnv(getenv("UEXKY_ENV", string(c.Env)))
	c.PostgresURI = getenv("PG_URI", c.PostgresURI)
	c.RedisURI = getenv("REDIS_URI", c.RedisURI)
	c.MigrationFiles = getenv("MIGRATION_FILES", c.MigrationFiles)
	c.Server.Domain = getenv("DOMAIN", c.Server.Domain)
	c.Server.APIDomain = getenv("API_DOMAIN", c.Server.APIDomain)
	c.Server.Proto = getenv("PROTO", c.Server.Proto)
	c.Server.Port = getenvInt("PORT", c.Server.Port)
	c.Server.Host = getenv("HOST", c.Server.Host)
	c.Mail.PrivateKey = getenv("MAILGUN_PRIVATE_KEY", c.Mail.PrivateKey)
	c.Mail.PublicKey = getenv("MAILGUN_PUBLIC_KEY", c.Mail.PublicKey)
	c.Mail.Domain = getenv("MAILGUN_DOMAIN", c.Mail.Domain)
}

func Load(filename string) error {
	setDefault()
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			return uerr.Wrap(uerr.InternalError, err, "open config file")
		}
		if err := toml.NewDecoder(f).Decode(&c); err != nil {
			return uerr.Wrap(uerr.InternalError, err, "read config file")
		}
	}
	c.filename = filename
	patchEnv()
	mf, err := filepath.Abs(c.MigrationFiles)
	if err != nil {
		return uerr.Wrapf(uerr.InternalError, err, "file migrations file path")
	}
	c.MigrationFiles = mf
	return nil
}

func Get() *Config {
	return &c
}

func GetEnv() RuntimeEnv {
	switch c.Env {
	case TestEnv:
		return TestEnv
	case ProdEnv:
		return ProdEnv
	default:
		return DevEnv
	}
}

type RuntimeEnv string

const (
	DevEnv  RuntimeEnv = "dev"
	TestEnv RuntimeEnv = "test"
	ProdEnv RuntimeEnv = "prod"
)

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
