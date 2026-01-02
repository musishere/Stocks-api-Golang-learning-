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
	_ "github.com/lib/pq"
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

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	query := "INSERT INTO stocks (name, price, company, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(query, stock.Name, stock.Price, stock.Company, stock.CreatedAt, stock.UpdatedAt)
	if err != nil {
		slog.Error("Error inserting stock into database", "error", err)
		return 0
	}

	return stock.ID
}

func getAllStocks() []models.Stock {
	db := CreateConnection()
	defer db.Close()

	query := "SELECT * FROM stocks"
	_, err := db.Query(query)
	if err != nil {
		slog.Error("Error fetching stocks from database", "error", err)
		return []models.Stock{}
	}
	return []models.Stock{}
}

func getStock(id int) models.Stock {
	db := CreateConnection()
	defer db.Close()

	query := "SELECT * FROM stocks WHERE id = $1"
	_, err := db.Query(query, id)
	if err != nil {
		slog.Error("Error fetching stock from database", "error", err)
		return models.Stock{}
	}
	return models.Stock{}
}

func updateStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	query := "UPDATE stocks SET name = $1, price = $2, company = $3, created_at = $4, updated_at = $5 WHERE id = $6"
	_, err := db.Exec(query, stock.Name, stock.Price, stock.Company, stock.CreatedAt, stock.UpdatedAt, stock.ID)
	if err != nil {
		slog.Error("Error updating stock in database", "error", err)
		return 0
	}
	return stock.ID
}

func deleteStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	query := "DELETE FROM stocks WHERE id = $1"
	_, err := db.Exec(query, stock.ID)
	if err != nil {
		slog.Error("Error deleting stock from database", "error", err)
		return 0
	}
	return stock.ID
}
