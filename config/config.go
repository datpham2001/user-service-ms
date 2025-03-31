package config

type Config struct {
	Env      string         `yaml:"env" mapstructure:"env"`
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Cache    CacheConfig    `yaml:"cache" mapstructure:"cache"`
	Jwt      JwtConfig      `yaml:"jwt" mapstructure:"jwt"`
}

type ServerConfig struct {
	Grpc struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port string `yaml:"port" mapstructure:"port"`
	} `yaml:"grpc" mapstructure:"grpc"`
	Http struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port string `yaml:"port" mapstructure:"port"`
	} `yaml:"http" mapstructure:"http"`
	TLS struct {
		Enable   bool   `yaml:"enable" mapstructure:"enable"`
		CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
		KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	} `yaml:"tls" mapstructure:"tls"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         string `yaml:"port" mapstructure:"port"`
	User         string `yaml:"user" mapstructure:"user"`
	Password     string `yaml:"password" mapstructure:"password"`
	DBName       string `yaml:"db_name" mapstructure:"db_name"`
	SSLMode      string `yaml:"ssl_mode" mapstructure:"ssl_mode"`
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
}

type CacheConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type JwtConfig struct {
	Secret string `yaml:"secret" mapstructure:"secret"`
}
