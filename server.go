package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// CotacaoResponse representa a estrutura da resposta da API de cotação
type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	// Conecta ao banco de dados SQLite
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria a tabela de cotações se não existir
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, valor REAL, data DATETIME)`)
	if err != nil {
		log.Fatal(err)
	}

	// Define o handler para a rota /cotacao
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		// Cria um contexto com timeout de 200ms para a requisição à API externa
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		// Obtém a cotação
		cotacao, err := getCotacao(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cria um contexto com timeout de 10ms para salvar no banco de dados
		ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelDB()

		// Salva a cotação no banco de dados
		err = saveCotacao(ctxDB, db, cotacao)
		if err != nil {
			log.Printf("Erro ao salvar cotação: %v", err)
		}

		// Retorna a cotação como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"bid": cotacao})
	})

	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getCotacao faz uma requisição à API externa para obter a cotação do dólar
func getCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var cotacaoResp CotacaoResponse
	err = json.Unmarshal(body, &cotacaoResp)
	if err != nil {
		return "", err
	}

	return cotacaoResp.USDBRL.Bid, nil
}

// saveCotacao salva a cotação no banco de dados
func saveCotacao(ctx context.Context, db *sql.DB, valor string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (valor, data) VALUES (?, datetime('now'))", valor)
	return err
}
