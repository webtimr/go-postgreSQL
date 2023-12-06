package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/webtimr/go-postgreSQL.git/models"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URl"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("successfuly connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("unable to decode the requst body. %v", err)
	}

	insertID := insertStock(stock)
	res := response{
		ID:      insertID,
		Message: "stock created successfuly",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("unable to convert the string into int. %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStock()

	if err != nil {
		log.Fatalf("unable to get all the stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("unable to convert the string into int. %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("unable to decode the requst body. %v", err)
	}

	updateRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("stock updated successfuly. total rows/records affected %v", updateRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string to int %v", err)
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("stock deleted successfuly. Total rows/records %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	fmt.Printf("inserted a single record %v", err)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("now rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row %v", err)
	}

	return stock, err
}

func getAllStock() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("error while checking the affected rows %v", err)
	}

	fmt.Printf("total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("error while checking the affected rows %v", err)
	}

	fmt.Printf("total rows/record affected %v", rowsAffected)

	return rowsAffected
}
