package main

import (
	"context"
	"encoding/json"
	"jwt-test/auth"
	"jwt-test/database"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

	client, err := database.ConnectDB()

	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return
	}
	defer client.Disconnect(context.Background())

	db_name := os.Getenv("DB_NAME")
	collection_name := os.Getenv("COLLECTION_NAME")

	collection := client.Database(db_name).Collection(collection_name)

	router := mux.NewRouter()

	// Маршрут для создания пары токенов.
	router.HandleFunc("/auth/{guid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		guid := vars["guid"]

		tokenPair, err := auth.CreateToken(guid, client, collection)

		if err != nil {
			errorResponse := ErrorResponse{Error: "Ошибка при создании токенов"}
			jsonResponse, _ := json.Marshal(errorResponse)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonResponse)
			return
		}

		json.NewEncoder(w).Encode(tokenPair)

	}).Methods("POST")

	// Маршрут для обновления пары токенов по Refresh токену.
	router.HandleFunc("/auth/refresh/{guid}/{refreshToken}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		guid := vars["guid"]
		refreshToken := vars["refreshToken"]

		tokenPair, err := auth.RefreshToken(guid, refreshToken, collection, client)
		if err != nil {
			errorResponse := ErrorResponse{Error: "Ошибка при обновлении токенов"}
			jsonResponse, _ := json.Marshal(errorResponse)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonResponse)
			return
		}

		json.NewEncoder(w).Encode(tokenPair)

	}).Methods("POST")

	http.ListenAndServe(":8080", router)

}
