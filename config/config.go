package config

type Config struct {
	PostgresURI string `json:"postgres_uri"`
	RedisURL    string `json:"redis_url"`
	Server      struct {
		Proto     string `json:"proto"`
		Domain    string `json:"domain"`
		APIDomain string `json:"api_domain"`
	} `json:"server"`
	Mail struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	} `json:"mail"`
	RateLimit struct {
		HTTPHeader     string `json:"http_header"`
		QueryLimit     int    `json:"query_limit"`
		QueryResetTime int    `json:"query_reset_time"`
		MutLimit       int    `json:"mut_limit"`
		MutResetTime   int    `json:"mut_reset_time"`
		Cost           struct {
			CreateUser int `json:"create_user"`
			PubThread  int `json:"pub_thread"`
			PubPost    int `json:"pub_post"`
		} `json:"cost"`
	} `json:"rate_limit"`
}

var c = Config{}

func Get() *Config {
	return &c
}
