package middleware

import (
	"database/sql"
	"encoding/json"
	"go-postgres-stocks/models"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message"`
	// Data    json{}
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging database", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	slog.Info("Successfully connected to database")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		slog.Error("Error decoding stock data", "error", err)
		json.NewEncoder(w).Encode(response{ID: 0, Message: "Error decoding stock data"})
		return
	}
	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
	slog.Info("Stock created successfully")

}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting all stocks")
	stocks := getAllStocks()
	res := response{
		ID:      stocks[0].ID,
		Message: "Stocks fetched successfully",
	}
	json.NewEncoder(w).Encode(res)
	slog.Info("Stocks fetched successfully")
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting stock")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		slog.Error("Error converting id to int", "error", err)
		json.NewEncoder(w).Encode(response{ID: 0, Message: "Error converting id to int"})
		return
	}
	stock := getStock(id)
	res := response{
		ID:      stock.ID,
		Message: "Stock fetched successfully",
	}

	json.NewEncoder(w).Encode(res)
	slog.Info("Stock fetched successfully")
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	slog.Info("Updating stock")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		slog.Error("Error converting id to int", "error", err)
		json.NewEncoder(w).Encode(response{ID: 0, Message: "Error converting id to int"})
		return
	}
	stock := getStock(id)
	json.NewDecoder(r.Body).Decode(&stock)
	insertID := updateStock(stock)
	res := response{
		ID:      insertID,
		Message: "Stock updated successfully",
	}
	json.NewEncoder(w).Encode(res)
	slog.Info("Stock updated successfully")
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	slog.Info("Deleting stock")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		slog.Error("Error converting id to int", "error", err)
		json.NewEncoder(w).Encode(response{ID: 0, Message: "Error converting id to int"})
		return
	}
	stock := getStock(id)
	insertID := deleteStock(stock)
	res := response{
		ID:      insertID,
		Message: "Stock deleted successfully",
	}
	json.NewEncoder(w).Encode(res)
	slog.Info("Stock deleted successfully")
}
