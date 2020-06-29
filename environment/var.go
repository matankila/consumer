package environment

import "os"

var (
	RedisPort = os.Getenv("REDIS_PORT")
	//RedisUser     = os.Getenv("REDIS_USER")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	RabbitPort     = os.Getenv("RABBIT_PORT")
	RabbitUser     = os.Getenv("RABBIT_USER")
	RabbitPassword = os.Getenv("RABBIT_PASSWORD")

	ServerPort = os.Getenv("SERVER_PORT")
)
