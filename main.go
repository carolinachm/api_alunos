package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Inicializa conexão com o banco
	InitDB()
	defer db.Close()

	// 2. Configura todas as rotas (chama o arquivo routes.go)
	mux := ConfigurarRotas()

	// 3. Aplica middleware CORS
	servidor := CorsMiddleware(mux)

	// 4. Inicia o servidor 
	porta := ":8080"
	fmt.Println("Servidor rodando em http://localhost", porta)

	err := http.ListenAndServe(porta, servidor)
	if err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}