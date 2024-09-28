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

// CotacaoResponse representa a estrutura da resposta do servidor
type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	// Cria um contexto com timeout de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Prepara a requisição para o servidor
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Realiza a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Decodifica a resposta JSON
	var cotacaoResp CotacaoResponse
	err = json.Unmarshal(body, &cotacaoResp)
	if err != nil {
		log.Fatal(err)
	}

	// Salva a cotação em um arquivo
	err = saveCotacao(cotacaoResp.Bid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cotação do dólar salva com sucesso: %s\n", cotacaoResp.Bid)
}

// saveCotacao salva a cotação em um arquivo texto
func saveCotacao(valor string) error {
	content := fmt.Sprintf("Dólar: %s", valor)
	return ioutil.WriteFile("cotacao.txt", []byte(content), 0644)
}
