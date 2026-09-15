// Pacote principal
package main

// Importa os pacotes necessários
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Aluno representa o modelo de dados de um aluno
type Aluno struct {
	Codigo   string  `json:"codigo"`
	Nome     string  `json:"nome"`
	Nota1    float64 `json:"nota1"`
	Nota2    float64 `json:"nota2"`
	Media    float64 `json:"media"`
	Situacao string  `json:"situacao"`
}

// Lista de alunos armazenada em memória
var alunos = []Aluno{}

// calcularMedia calcula a média e define a situação do aluno
func calcularMedia(aluno *Aluno) {
	// Calcula a média aritmética
	aluno.Media = (aluno.Nota1 + aluno.Nota2) / 2

	// Define a situação conforme a média
	switch {
	case aluno.Media >= 7:
		aluno.Situacao = "Aprovado"
	case aluno.Media >= 5:
		aluno.Situacao = "Recuperação"
	default:
		aluno.Situacao = "Reprovado"
	}
}

// Função Middleware para gerenciar e liberar as permissões de CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Permitir requisições de qualquer origem (Frontend)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 2. Definir quais métodos HTTP externos são aceitos pela API
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 3. Definir os cabeçalhos permitidos nas requisições (obrigatório para JSON)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 4. Tratar a requisição preflight (pedido de verificação prévia do navegador)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Se tudo estiver certo, passa a requisição adiante para as nossas rotas
		next.ServeHTTP(w, r)
	})
}

// listarAlunos retorna todos os alunos cadastrados em formato JSON
func listarAlunos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alunos)
}

// cadastrarAluno recebe os dados, gera um código, calcula a média e cadastra o aluno
func cadastrarAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var aluno Aluno
	err := json.NewDecoder(r.Body).Decode(&aluno)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Gera identificador único
	aluno.Codigo = uuid.New().String()

	// Calcula média e situação
	calcularMedia(&aluno)

	// Adiciona à lista
	alunos = append(alunos, aluno)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(aluno)
}

// alterarAluno busca o aluno pelo código e atualiza seus dados
func alterarAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Obtém o código da URL
	codigo := r.PathValue("codigo")

	// Percorre a lista procurando o aluno
	for indice := range alunos {
		if alunos[indice].Codigo == codigo {
			//Objeto do tipo Aluno para receber os dados atualizados
			var aluno Aluno
			// Decodifica o JSON recebido
			err := json.NewDecoder(r.Body).Decode(&aluno)
			if err != nil {
				http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
				return
			}

			// Preserva o código original
			aluno.Codigo = codigo

			// Recalcula média e situação
			calcularMedia(&aluno)

			// Atualiza os dados na lista
			alunos[indice] = aluno

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(aluno)
			return
		}
	}

	// Aluno não encontrado
	http.Error(w, "Aluno não encontrado", http.StatusNotFound)
}

func removerAluno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Obtém o código da URL
	codigo := r.PathValue("codigo")

	// Percorre a lista procurando o aluno
	for indice := range alunos {
		if alunos[indice].Codigo == codigo {
			// Remove o aluno da lista
			alunos = append(alunos[:indice], alunos[indice+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Aluno não encontrado
	http.Error(w, "Aluno não encontrado", http.StatusNotFound)
}

// Função principal — inicializa o servidor e define as rotas
func main() {
	// Rotas
	http.HandleFunc("GET /alunos", listarAlunos)
	http.HandleFunc("POST /alunos", cadastrarAluno)
	http.HandleFunc("PUT /alunos/{codigo}", alterarAluno)
	http.HandleFunc("DELETE /alunos/{codigo}", removerAluno)

	// Inicia o servidor
	fmt.Println("Servidor em execução em: http://localhost:8080")
	// Configurar servidor aplicando o nosso middleware de CORS sobre o roteador padrão (http.DefaultServeMux)
	http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux))
}
