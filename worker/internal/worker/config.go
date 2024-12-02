package worker

type Config struct {
	BaseUrl     string
	LogLevel    string
	RabbitMqUrl string
}

func NewConfig(baseUrl, rabbitMQUrl string) *Config {
	return &Config{
		BaseUrl:     baseUrl,
		RabbitMqUrl: rabbitMQUrl,
		LogLevel:    "debug",
	}
}
