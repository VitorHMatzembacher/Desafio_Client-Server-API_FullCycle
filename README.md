# Desafio Client-Server API - FullCycle

## Descrição

O projeto consiste em dois sistemas:

- `server.go`  Servidor HTTP responsável por consumir a cotação do dólar de uma API externa, salvar no banco e retornar para o cliente.
- `client.go`  Cliente que faz uma requisição HTTP para o servidor, obtém a cotação e salva em um arquivo `.txt`.

### Server (`server.go`)

- Endpoint HTTP `/cotacao` na porta `8080`.
- Consome a API pública:  
  ➝ `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- Timeout de **200ms** para obter a cotação da API.
- Persiste o valor no banco SQLite (`database.db`) com timeout de **10ms**.
- Retorna para o cliente o valor do campo `"bid"` no formato JSON.

### Client (`client.go`)

- Faz uma requisição para `http://localhost:8080/cotacao`.
- Timeout de **300ms** para resposta do servidor.
- Salva o valor recebido no arquivo `cotacao.txt` no formato:  
  ➝ `Dólar: {valor}`

