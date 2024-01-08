# Currency Exchange App

Este é um exemplo de aplicativo de câmbio de moeda que consiste em um servidor (`server.go`) e um cliente (`client.go`). O servidor fornece a cotação do dólar, enquanto o cliente realiza uma solicitação ao servidor, recebe a cotação e a salva em um arquivo.

## Como Rodar

### Pré-requisitos

- Go instalado: [Instruções de Instalação](https://golang.org/doc/install)

### Rodar o Servidor

1. Abra um terminal e navegue até o diretório onde o arquivo `server.go` está localizado.

2. Execute o seguinte comando para iniciar o servidor:

   ```bash
   go run server.go
   ```

O servidor estará ouvindo em http://localhost:8080/cotacao.

### Rodar o Cliente

1. Abra outro terminal e navegue até o diretório onde o arquivo client.go está localizado.

2. Execute o seguinte comando para iniciar o cliente:
   ```bash
   go run client.go
   ```

O cliente fará uma solicitação ao servidor para obter a cotação do dólar e salvará o valor no arquivo cotacao.txt.
