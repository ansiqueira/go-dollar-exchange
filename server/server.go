package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	apiCtx, apiCancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer apiCancel()

	exchangeRate, err := getExchangeRate(apiCtx)
	if err != nil {
		log.Printf("Erro ao obter cotação da API: %v", err)
		http.Error(w, "Erro ao obter cotação da API", http.StatusInternalServerError)
		return
	}

	db, err := getDBConnection()

	if err != nil {
		log.Printf("Erro ao conectar no banco de dados: %v", err)
		http.Error(w, "Erro ao conectar o banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = saveToDatabase(ctx, exchangeRate, db)
	if err != nil {
		log.Printf("Erro ao persistir os dados no banco de dados: %v", err)
		http.Error(w, "Erro ao persistir os dados no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchangeRate)
}

func getExchangeRate(ctx context.Context) (USDBRL, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return USDBRL{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return USDBRL{}, err
	}
	defer resp.Body.Close()

	var USDBRLData map[string]USDBRL
	err = json.NewDecoder(resp.Body).Decode(&USDBRLData)
	if err != nil {
		return USDBRL{}, err
	}

	if len(USDBRLData) == 0 {
		return USDBRL{}, fmt.Errorf("Resposta da API vazia ou inválida")
	}

	return USDBRLData["USDBRL"], nil
}

func saveToDatabase(ctx context.Context, exchangeRate USDBRL, db *sql.DB) error {
	saveCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)

	sql.Open("postgres", "")

	defer cancel()

	_, err := db.ExecContext(saveCtx, `
		INSERT INTO exchange_rates (
			code, codein, name, high, low, var_bid,
			pct_change, bid, ask, timestamp, create_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		exchangeRate.Code, exchangeRate.Codein, exchangeRate.Name, exchangeRate.High, exchangeRate.Low,
		exchangeRate.VarBid, exchangeRate.PctChange, exchangeRate.Bid, exchangeRate.Ask,
		exchangeRate.Timestamp, exchangeRate.CreateDate)

	return err
}

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS exchange_rates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			codein TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid TEXT,
			ask TEXT,
			timestamp TEXT,
			create_date TEXT
		)`)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
