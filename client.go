package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro criando request:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro na requisição:", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal("Erro decodificando resposta:", err)
	}

	bid := result["bid"]
	fmt.Println("Cotação do dólar:", bid)

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("Erro criando arquivo:", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
	if err != nil {
		log.Fatal("Erro escrevendo no arquivo:", err)
	}

	log.Println("Cotação salva em cotacao.txt")
}
