package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	UsdBrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	http.HandleFunc("/cotacao", handlerCotacao)
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cotacao, err := buscaCotacao(ctx)
	if err != nil {
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		log.Println("Erro na busca da cotação:", err)
		return
	}

	err = salvaCotacao(ctx, cotacao.UsdBrl.Bid)
	if err != nil {
		log.Println("Erro ao salvar no banco:", err)
	}

	resp := map[string]string{"bid": cotacao.UsdBrl.Bid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func buscaCotacao(ctx context.Context) (Cotacao, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return Cotacao{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Cotacao{}, err
	}
	defer resp.Body.Close()

	var c Cotacao
	err = json.NewDecoder(resp.Body).Decode(&c)
	return c, err
}

func salvaCotacao(ctx context.Context, bid string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid, data) VALUES (?, datetime('now'))", bid)
	return err
}

func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bid" TEXT,
		"data" DATETIME
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Erro criando tabela:", err)
	}
}
