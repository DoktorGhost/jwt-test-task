package auth

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func CreateToken(guid string, client *mongo.Client, collection *mongo.Collection) (*TokenPair, error) {

	jwtSecretKey := []byte(os.Getenv("JWT_KEY"))

	// Создание Access токена
	payload := jwt.MapClaims{
		"sub": guid,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	accessTokenString, err := accessToken.SignedString(jwtSecretKey)

	if err != nil {
		log.Printf("Ошибка при подписи JWT токена: %v", err)
		return nil, err
	}

	// Создание Refresh токена
	refreshUUID := uuid.New().String()
	refreshToken := base64.StdEncoding.EncodeToString([]byte(refreshUUID))

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Ошибка хэширования refresh токена: %v", err)
		return nil, err
	}

	// Сохраняем хэшированный Refresh токен в базе данных
	err = updateRefresToken(guid, hashedToken, collection)
	if err != nil {
		log.Printf("Ошибка при записи refresh токена в бд: %v", err)
		return nil, err
	}

	log.Println("Acces и Refrsh токены созданы")
	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshToken(guid, refreshToken string, collection *mongo.Collection, client *mongo.Client) (*TokenPair, error) {

	var result struct {
		// хэш
		RefreshToken string `bson:"refresh_token"`
	}

	err := collection.FindOne(context.Background(), bson.M{"guid": guid}).Decode(&result)

	if err != nil {
		log.Printf("Ошибка при поиске пользователя: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.RefreshToken), []byte(refreshToken))

	if err != nil {
		log.Printf("Refresh токен не совпадает: %v", err)
		return nil, err
	}

	tokenPair, err := CreateToken(guid, client, collection)

	if err != nil {
		log.Printf("Ошибка создания токенов: %v", err)
		return nil, err
	}

	log.Println("Acces и Refrsh токены успешно обновлены")
	return tokenPair, nil
}

func updateRefresToken(guid string, refreshToken []byte, collection *mongo.Collection) error {

	err := collection.FindOne(context.Background(), bson.M{"guid": guid}).Err()

	if err != nil {
		log.Printf("Ошибка при поиске: %v", err)
		return err
	}

	filter := bson.M{"guid": guid}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Ошибка обновления refresh токена: %v", err)
		return err
	}

	log.Println("Refrsh токены успешно обновлен в базе данных")
	return nil
}
