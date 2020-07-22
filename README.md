# uexky

Uexky is the backend of Abyss. (repo: https://gitlab.com/abyss.club, website: https://abyss.club).

Abyss is an anonymous-able forum with tag-managed topics.
We aim to provide an efficient and convenient discussion environment for sub-cultural topics.


## Building
`uexky` is written in `go`. Follow instructions [here](https://golang.org/) if you haven't installed `go` yet.

From repository root:

```shell
make build
```

This command will use `go mod` to download dependencies.
Binaries will be compiled after dependencies are satisfied, then placed at `./dist/uexky`.

## Configuration

### Config file

`uexky` can accept a `toml` config file. Structure of the file is:

```go
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
}
```

Default values: 

```go
PostgresURI    = "postgres://postgres:postgres@localhost:5432/uexky2?sslmode=disable"
RedisURI       = "redis://localhost:6379/0"
MigrationFiles = "./migrations"
Server.Domain  = "abyss.club"
Server.Domain  = "api.abyss.club"
Server.Proto   = "http"
Server.Port    = 8000
Server.Host    = "localhost"
```

### Environments

Configurations are overridden when following environment variables are set.
**Environment variables will take precedence** over the `toml` config file when there're conflicts.

```shell
UEXKY_ENV               // running env: prod, test, dev
PG_URI                  // postgres URI
REDIS_URI               // redis URI
MIGRATION_FILES         // migration files dir
DOMAIN                  // frontend(`tt392`) domain
API_DOMAIN              // api domain
PROTO                   // network proto: http, https
PORT                    // listening port
HOST                    // listening host
MAILGUN_PRIVATE_KEY     // mailgun private key
MAILGUN_PUBLIC_KEY      // mailgun public key
MAILGUN_DOMAIN          // mail domain
```

## Usage

### Prerequisites

- redis >= 6.0.6
- PostgresQL >= 12.3

### Database preparations

Before running `uexky` for the first time, db migrations and setting "main tags" are required: 

```shell
./dist/uexky -c config.toml migrate up
./dist/uexky -c config.toml admin settags "anime,games"
```

### Running uexky

Starting up `uexky`:

```shell
./dist/uexky -c config.toml
```

More cli usage:

```shell
./dist/uexky --help
```

## `docker`

The docker image registry is `registry.gitlab.com/abyss.club/uexky:latest`.  
Basic usage:

```shell
docker run --rm registry.gitlab.com/abyss.club/uexky:latest uexky --help
```

## Contributing

Feel free to file an issue on [our GitLab repo](https://gitlab.com/abyss.club/uexky) if you have any question or feedback.
As this project is in early development, if you're interested in contribution to the project, please contact us directly.
