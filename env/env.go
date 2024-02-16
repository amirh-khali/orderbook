package env

import (
	"github.com/joho/godotenv"
	"os"
)

type env struct {
	BootstrapServers string
	KafkaTopic       string
}

var ENV env

func LoadEnv() {
	_ = godotenv.Load("env/.env")

	ENV = env{
		BootstrapServers: os.Getenv("BOOTSTRAP_SERVERS"),
		KafkaTopic:       os.Getenv("KAFKA_TOPIC"),
	}
}
