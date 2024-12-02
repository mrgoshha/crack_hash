package manager

type Config struct {
	BaseUrl     string `toml:"bind_addr"`
	MongoUrl    string `toml:"mongo_url"`
	RabbitMQUrl string `toml:"rabbit_url"`
	LogLevel    string `toml:"log_level"`
}

func NewConfig(mongoUrl, rabbitMQUrl string) *Config {
	return &Config{
		BaseUrl:     ":8080",
		MongoUrl:    mongoUrl,
		RabbitMQUrl: rabbitMQUrl,
		LogLevel:    "debug",
	}
}
