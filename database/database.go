package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
		return nil, err
	}

	/*
		username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
		password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
		uri := "mongodb://" + username + ":" + password + "@localhost:27017"
	*/

	uri := os.Getenv("URL_MONGO")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Подключение к БД успешно")
	return client, nil
}
