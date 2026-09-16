package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Calcula média e situação
func mediaSituacao(aluno *Aluno) {
	aluno.Media = (aluno.Nota1 + aluno.Nota2) / 2
	switch {
	case aluno.Media >= 7:
		aluno.Situacao = "Aprovado(a)"
	case aluno.Media >= 5:
		aluno.Situacao = "Em Recuperação"
	default:
		aluno.Situacao = "Reprovado(a)"
	}
}

// Listar → chama função do database.go
func ListarAlunos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	lista, err := SelecionarTodos()
	if err != nil {
		http.Error(w, "Erro ao consultar", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lista)
}

// Cadastrar → chama função do database.go
func CadastrarAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var aluno Aluno
	if err := json.NewDecoder(r.Body).Decode(&aluno); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	aluno.Codigo = uuid.New().String()
	mediaSituacao(&aluno)
	if err := Inserir(aluno); err != nil {
		http.Error(w, "Erro ao salvar", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(aluno)
}

// Editar → chama função do database.go
func AlterarAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	codigo := r.PathValue("codigo")
	var aluno Aluno
	if err := json.NewDecoder(r.Body).Decode(&aluno); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	aluno.Codigo = codigo
	mediaSituacao(&aluno)
	if err := Atualizar(aluno); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(aluno)
}

// Remover → chama função do database.go
func RemoverAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	codigo := r.PathValue("codigo")
	if err := Excluir(codigo); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Rota de teste
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"mensagem": "API funcionando! ",
	})
}