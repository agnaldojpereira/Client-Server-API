# Projeto de Cotação do Dólar

Este projeto consiste em dois programas em Go que trabalham juntos para obter e armazenar a cotação atual do dólar em relação ao real brasileiro.

## Componentes

1. `server.go`: Um servidor HTTP que obtém a cotação do dólar de uma API externa e a armazena em um banco de dados SQLite.
2. `client.go`: Um cliente que faz uma requisição ao servidor para obter a cotação atual e a salva em um arquivo de texto.

## Funcionalidades

### Server (`server.go`)

- Cria um servidor HTTP na porta 8080.
- Expõe um endpoint `/cotacao` que retorna a cotação atual do dólar.
- Utiliza a API https://economia.awesomeapi.com.br/json/last/USD-BRL para obter a cotação.
- Armazena cada cotação recebida em um banco de dados SQLite.
- Implementa timeouts:
  - 200ms para chamar a API de cotação.
  - 10ms para persistir os dados no banco.

### Client (`client.go`)

- Faz uma requisição HTTP para o servidor na rota `/cotacao`.
- Recebe a cotação atual do dólar.
- Salva a cotação em um arquivo `cotacao.txt`.
- Implementa um timeout de 300ms para receber o resultado do servidor.

## Requisitos

- Go 1.15 ou superior
- SQLite3

## Como executar

1. Certifique-se de ter o Go e o SQLite3 instalados em seu sistema.
2. Clone este repositório:
   ```
   git clone https://seu-repositorio.git
   cd seu-repositorio
   ```
3. Instale as dependências:
   ```
   go mod tidy
   ```
4. Em um terminal, execute o servidor:
   ```
   go run server.go
   ```
5. Em outro terminal, execute o cliente:
   ```
   go run client.go
   ```

## Observações

- O servidor utiliza o pacote `context` para gerenciar timeouts.
- O cliente também utiliza `context` para limitar o tempo de espera da resposta.
- Erros de timeout são registrados nos logs.
- A cotação é salva no formato "Dólar: {valor}" no arquivo `cotacao.txt`.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

