package env

import (
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	BootstrapServers     string
	KafkaTopic           string
	KafkaGroupID         string
	KafkaAutoOffsetReset string
	MongoURI             string
	MongoDatabaseName    string
}

var ENV env

func LoadEnv() {
	_ = godotenv.Load("env/.env")

	ENV = env{
		BootstrapServers:     os.Getenv("BOOTSTRAP_SERVERS"),
		KafkaTopic:           os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:         os.Getenv("KAFKA_GROUP_ID"),
		KafkaAutoOffsetReset: os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
		MongoURI:             os.Getenv("MONGO_URI"),
		MongoDatabaseName:    os.Getenv("MONGO_DATABASE_NAME"),
	}
}
