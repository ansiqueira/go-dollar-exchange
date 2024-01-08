package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erro ao fazer a requisição HTTP:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Erro na resposta do servidor: %s", resp.Status)
		return
	}

	var exchangeRate ExchangeRate
	err = json.NewDecoder(resp.Body).Decode(&exchangeRate)
	if err != nil {
		log.Fatal("Erro ao decodificar a resposta JSON:", err)
		return
	}

	err = saveToFile("cotacao.txt", exchangeRate.Bid)
	if err != nil {
		log.Fatal("Erro ao salvar a cotação no arquivo:", err)
		return
	}

	fmt.Printf("Cotação do dólar salva com sucesso: %s\n", exchangeRate.Bid)
}

func saveToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(fmt.Sprintf("Dólar: %s\n", content)), 0644)
}
