package config

// Config struct
type Config struct {
	App     AppConfig     `yaml:"app"`
	Elastic ElasticConfig `yaml:"elastic"`
}

// AppConfig struct
type AppConfig struct {
	Port string `yaml:"port"`
}

// DatabasesConfig struct
type DatabasesConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	MySQL    MySQLConfig    `yaml:"mysql"`
}

// PostgresConfig struct
type PostgresConfig struct {
	Master string `yaml:"master"`
	Slave  string `yaml:"slave"`
}

// MySQLConfig struct
type MySQLConfig struct {
	Master string `yaml:"master"`
	Slave  string `yaml:"slave"`
}

// ElasticConfig struct
type ElasticConfig struct {
	Host string `yaml:"host"`
}
