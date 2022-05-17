package config

type Postgres struct {
	User         string `env:"PG_USER"`
	Password     string `env:"PG_PASSWORD"`
	Host         string `env:"PG_HOST"`
	Port         int    `env:"PG_PORT"`
	Database     string `env:"PG_DATABASE"`
	TestDatabase string `env:"PG_TEST_DATABASE"`
}
