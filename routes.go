package main

import "net/http"

// ConfigurarRotas registra todas as rotas do sistema
func ConfigurarRotas() *http.ServeMux {
	mux := http.NewServeMux()

	// Rota de teste
	mux.HandleFunc("GET /", HelloWorld)

	// Rotas de Alunos
	mux.HandleFunc("GET /alunos", ListarAlunos)
	mux.HandleFunc("POST /alunos", CadastrarAluno)
	mux.HandleFunc("PUT /alunos/{codigo}", AlterarAluno)
	mux.HandleFunc("DELETE /alunos/{codigo}", RemoverAluno)

	return mux
}