// Comando para importar o pacote PostgreSQL: go get github.com/lib/pq

// Pacote principal da aplicacao.
// O pacote "main" indica que este arquivo pertence a um programa executavel.
package main

// Importacoes
import (
	// database/sql fornece a interface generica para trabalhar
	// com bancos de dados relacionais.
	"database/sql"

	// encoding/json permite converter estruturas Go para JSON
	// e JSON para estruturas Go.
	"encoding/json"

	// fmt e utilizado para exibir mensagens no terminal
	// e formatar textos.
	"fmt"

	//os é um pacote nativo do Go que fornece uma interface para interagir com o sistema operacional.
	"os"

	// log permite registrar mensagens de erro e encerrar
	// a aplicacao em situacoes criticas.
	"log"

	// net/http fornece recursos para criar o servidor HTTP
	// e trabalhar com requisicoes e respostas.
	"net/http"

	// Biblioteca responsavel pela geracao de UUIDs.
	"github.com/google/uuid"

	// Driver PostgreSQL utilizado pelo database/sql.
	//
	// O "_" significa que nao vamos utilizar diretamente
	// nenhuma funcao desse pacote no codigo.
	// O objetivo do import e registrar o driver "postgres".
	 _"github.com/lib/pq"
)

// ============================================================
// MODELO / STRUCT
// ============================================================

// Aluno representa um aluno dentro da aplicacao.
//
// Essa struct funciona como o modelo dos dados que serao
// recebidos, enviados e armazenados no banco.
type Aluno struct {

	// Codigo unico do aluno.
	Codigo string `json:"codigo"`

	// Nome do aluno.
	Nome string `json:"nome"`

	// Primeira nota.
	Nota1 float64 `json:"nota1"`

	// Segunda nota.
	Nota2 float64 `json:"nota2"`

	// Media calculada pela aplicacao.
	Media float64 `json:"media"`

	// Situacao do aluno:
	// Aprovado(a), Em Recuperacao ou Reprovado(a).
	Situacao string `json:"situacao"`
}

// ============================================================
// CONEXAO COM O BANCO
// ============================================================

// Variavel global que representa a conexao com o banco.
//
// O *sql.DB nao representa necessariamente uma unica conexao.
// Ele gerencia um pool de conexoes com o banco de dados.
var db *sql.DB

// ============================================================
// INICIALIZACAO DO BANCO DE DADOS
// ============================================================

// initDB inicializa a conexao com o PostgreSQL
// e cria a tabela "alunos", caso ela ainda nao exista.
func initDB() {
	// bdURL deve possuir as informacoes necessarias para realizar a conexao com o banco de dados
	dbURL := os.Getenv("DATABASE_URL")
	

	// Declara a variavel que armazenara possiveis erros.
	var err error

	// Abre a conexao com o PostgreSQL.
	//
	// "postgres" informa qual driver sera utilizado.
	db, err = sql.Open("postgres", dbURL)

	// Verifica se houve erro ao abrir/preparar a conexao.
	if err != nil {
		log.Fatalf("Erro ao abrir conexÃ£o com o banco: %v", err)
	}

	// Verifica se realmente conseguimos estabelecer
	// comunicacao com o banco.
	//
	// sql.Open() sozinho nao garante que o banco esteja acessivel.
	// O Ping() faz essa verificacao.
	err = db.Ping()

	if err != nil {
		log.Fatalf("Erro ao conectar no banco do Render: %v", err)
	}

	// Se chegamos aqui, a conexao foi estabelecida.
	fmt.Println("Conectado ao PostgreSQL com sucesso!")

	// ========================================================
	// CRIACAO DA TABELA
	// ========================================================

	// Comando SQL responsavel por criar a tabela.
	//
	// IF NOT EXISTS significa:
	// "crie a tabela somente se ela ainda nao existir".
	query := `
	CREATE TABLE IF NOT EXISTS alunos (
		codigo VARCHAR(36) PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		nota1 NUMERIC(4,2) NOT NULL,
		nota2 NUMERIC(4,2) NOT NULL,
		media NUMERIC(4,2) NOT NULL,
		situacao VARCHAR(20) NOT NULL
	);`

	// Executa o comando SQL no PostgreSQL.
	_, err = db.Exec(query)

	// Verifica se houve algum problema na criacao da tabela.
	if err != nil {
		log.Fatalf("Erro ao criar tabela alunos: %v", err)
	}
}

// ============================================================
// CALCULO DA MEDIA E SITUACAO
// ============================================================

// mediaSituacao calcula a media do aluno
// e determina sua situacao.
//
// O parametro e um ponteiro (*Aluno), pois queremos
// alterar os valores Media e Situacao do proprio aluno.
func mediaSituacao(aluno *Aluno) {

	// Calcula a media das duas notas.
	aluno.Media = (aluno.Nota1 + aluno.Nota2) / 2

	// Verifica a situacao do aluno.
	if aluno.Media >= 7 {

		// Media maior ou igual a 7:
		// aluno aprovado.
		aluno.Situacao = "Aprovado(a)"

	} else if aluno.Media >= 5 {

		// Media entre 5 e 6.99:
		// aluno em recuperacao.
		aluno.Situacao = "Em RecuperaÃ§Ã£o"

	} else {

		// Media abaixo de 5:
		// aluno reprovado.
		aluno.Situacao = "Reprovado(a)"
	}
}

// ============================================================
// HELLO WORLD
// ============================================================

// helloWorld e uma funcao responsavel pela rota:
//
// GET /
//
// Ela serve como uma rota simples para testar
// se a API esta funcionando.
func helloWorld(w http.ResponseWriter, r *http.Request) {

	// Define que a resposta sera JSON.
	w.Header().Set("Content-Type", "application/json")

	// Define o status HTTP da resposta.
	w.WriteHeader(http.StatusCreated)

	// Gera um UUID unico.
	codigoUnico := uuid.New().String()

	// Cria um mapa que sera convertido para JSON.
	mensagem := map[string]string{
		"codigoUnico": codigoUnico,
		"mensagem":    "Hello World!",
	}

	// Converte o mapa para JSON e envia para o cliente.
	json.NewEncoder(w).Encode(mensagem)
}

// ============================================================
// LISTAR ALUNOS
// ============================================================

// listarAlunos e responsavel pela rota:
//
// GET /alunos
//
// Busca todos os alunos cadastrados no PostgreSQL.
func listarAlunos(w http.ResponseWriter, r *http.Request) {

	// Define que a resposta sera JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Executa uma consulta SQL para buscar os alunos.
	rows, err := db.Query(`
		SELECT codigo, nome, nota1, nota2, media, situacao
		FROM alunos
	`)

	// Verifica se houve erro na consulta.
	if err != nil {
		http.Error(
			w,
			"Erro ao buscar alunos no banco",
			http.StatusInternalServerError,
		)
		return
	}

	// Garante que as linhas da consulta serao fechadas
	// quando a funcao terminar.
	defer rows.Close()

	// Criamos uma lista vazia de alunos.
	//
	// Isso garante que, caso nao existam alunos,
	// o JSON retornado seja [] em vez de null.
	listaAlunos := []Aluno{}

	// Percorre cada registro retornado pelo banco.
	for rows.Next() {

		// Cria uma variavel que representa um aluno.
		var a Aluno

		// Transfere os valores das colunas do banco
		// para os campos da struct Aluno.
		err := rows.Scan(
			&a.Codigo,
			&a.Nome,
			&a.Nota1,
			&a.Nota2,
			&a.Media,
			&a.Situacao,
		)

		// Verifica se houve erro durante a leitura.
		if err != nil {
			http.Error(
				w,
				"Erro ao processar dados dos alunos",
				http.StatusInternalServerError,
			)
			return
		}

		// Adiciona o aluno a lista.
		listaAlunos = append(listaAlunos, a)
	}

	// Define o status HTTP 200 (OK).
	w.WriteHeader(http.StatusOK)

	// Converte a lista de alunos para JSON
	// e envia ao cliente.
	json.NewEncoder(w).Encode(listaAlunos)
}

// ============================================================
// CADASTRAR ALUNO
// ============================================================

// cadastrarAluno e responsavel pela rota:
//
// POST /alunos
//
// Recebe os dados de um aluno em JSON
// e salva o aluno no PostgreSQL.
func cadastrarAluno(w http.ResponseWriter, r *http.Request) {

	// Define o formato da resposta.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Cria uma variavel para armazenar
	// os dados recebidos no JSON.
	var aluno Aluno

	// Decodifica o JSON enviado pelo cliente
	// diretamente para a struct Aluno.
	erro := json.NewDecoder(r.Body).Decode(&aluno)

	// Verifica se o JSON e invalido.
	if erro != nil {
		http.Error(
			w,
			"Falha ao decodificar o JSON",
			http.StatusBadRequest,
		)
		return
	}

	// Gera automaticamente um codigo unico para o aluno.
	aluno.Codigo = uuid.New().String()

	// Calcula a media e a situacao.
	mediaSituacao(&aluno)

	// ========================================================
	// INSERT NO BANCO
	// ========================================================

	// Comando SQL para inserir o aluno.
	//
	// $1, $2, $3 etc. sao parametros que serao
	// substituidos pelos valores posteriormente.
	query := `
		INSERT INTO alunos
		(codigo, nome, nota1, nota2, media, situacao)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Executa o INSERT.
	_, err := db.Exec(
		query,
		aluno.Codigo,
		aluno.Nome,
		aluno.Nota1,
		aluno.Nota2,
		aluno.Media,
		aluno.Situacao,
	)

	// Verifica se houve erro ao salvar.
	if err != nil {
		http.Error(
			w,
			"Erro ao salvar aluno no banco",
			http.StatusInternalServerError,
		)
		return
	}

	// Cadastro realizado com sucesso.
	w.WriteHeader(http.StatusCreated)

	// Retorna o aluno cadastrado em JSON.
	json.NewEncoder(w).Encode(aluno)
}

// ============================================================
// ALTERAR ALUNO
// ============================================================

// alterarAluno e responsavel pela rota:
//
// PUT /alunos/{codigo}
//
// Atualiza os dados de um aluno existente.
func alterarAluno(w http.ResponseWriter, r *http.Request) {

	// Define o formato da resposta.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Obtem o codigo do aluno que esta na URL.
	//
	// Exemplo:
	// PUT /alunos/123
	//
	// codigo = "123"
	codigo := r.PathValue("codigo")

	// Cria uma variavel para receber
	// os novos dados enviados pelo cliente.
	var aluno Aluno

	// Decodifica o JSON do corpo da requisicao.
	erro := json.NewDecoder(r.Body).Decode(&aluno)

	// Verifica se o JSON e invalido.
	if erro != nil {
		http.Error(
			w,
			"Falha ao decodificar o JSON",
			http.StatusBadRequest,
		)
		return
	}

	// Mantem o codigo original recebido na URL.
	aluno.Codigo = codigo

	// Recalcula a media e a situacao
	// com as novas notas.
	mediaSituacao(&aluno)

	// ========================================================
	// UPDATE NO BANCO
	// ========================================================

	// Atualiza os dados do aluno cujo codigo
	// corresponde ao codigo recebido na URL.
	query := `
		UPDATE alunos
		SET nome = $1,
			nota1 = $2,
			nota2 = $3,
			media = $4,
			situacao = $5
		WHERE codigo = $6
	`

	// Executa o UPDATE.
	resultado, err := db.Exec(
		query,
		aluno.Nome,
		aluno.Nota1,
		aluno.Nota2,
		aluno.Media,
		aluno.Situacao,
		codigo,
	)

	// Verifica se ocorreu erro no banco.
	if err != nil {
		http.Error(
			w,
			"Erro ao atualizar dados no banco",
			http.StatusInternalServerError,
		)
		return
	}

	// Descobre quantas linhas foram alteradas.
	linhasAfetadas, _ := resultado.RowsAffected()

	// Se nenhuma linha foi alterada,
	// significa que o codigo nao foi encontrado.
	if linhasAfetadas == 0 {
		http.Error(
			w,
			"CÃ³digo nÃ£o encontrado",
			http.StatusNotFound,
		)
		return
	}

	// Retorna sucesso.
	w.WriteHeader(http.StatusOK)

	// Retorna o aluno atualizado.
	json.NewEncoder(w).Encode(aluno)
}

// ============================================================
// REMOVER ALUNO
// ============================================================

// removerAluno e responsavel pela rota:
//
// DELETE /alunos/{codigo}
//
// Remove um aluno do banco de dados.
func removerAluno(w http.ResponseWriter, r *http.Request) {

	// Define o formato da resposta.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Obtem o codigo do aluno atraves da URL.
	codigo := r.PathValue("codigo")

	// ========================================================
	// DELETE NO BANCO
	// ========================================================

	// Remove o aluno que possui o codigo informado.
	query := "DELETE FROM alunos WHERE codigo = $1"

	// Executa o DELETE.
	resultado, err := db.Exec(query, codigo)

	// Verifica se ocorreu algum erro.
	if err != nil {
		http.Error(
			w,
			"Erro ao remover aluno do banco",
			http.StatusInternalServerError,
		)
		return
	}

	// Verifica quantas linhas foram removidas.
	linhasAfetadas, _ := resultado.RowsAffected()

	// Se nenhuma linha foi removida,
	// o codigo informado nao existe.
	if linhasAfetadas == 0 {
		http.Error(
			w,
			"CÃ³digo nÃ£o encontrado",
			http.StatusNotFound,
		)
		return
	}

	// Retorna HTTP 204 (No Content),
	// indicando que a operacao foi realizada
	// sem necessidade de retornar um conteudo.
	w.WriteHeader(http.StatusNoContent)
}

// ============================================================
// CORS
// ============================================================

// corsMiddleware permite que aplicacoes de outros dominios
// facam requisicoes para nossa API.
//
// Isso e especialmente importante quando, por exemplo,
// um frontend esta hospedado em um endereco diferente
// da API.
func corsMiddleware(next http.Handler) http.Handler {

	// Retorna um novo Handler responsavel pelo CORS.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Permite requisicoes de qualquer origem.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Define os metodos HTTP permitidos.
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS",
		)

		// Define quais cabecalhos podem ser enviados.
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type",
		)

		// O navegador pode enviar uma requisicao OPTIONS
		// antes da requisicao real.
		//
		// Isso e chamado de preflight request.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Se nao for OPTIONS, continua normalmente
		// para o proximo Handler.
		next.ServeHTTP(w, r)
	})
}

// ============================================================
// FUNCAO PRINCIPAL
// ============================================================

// main e o ponto de entrada da aplicacao.
func main() {

	// 1. Inicializa a conexao com o banco
	// antes de iniciar o servidor HTTP.
	initDB()

	// Garante que a conexao com o banco
	// sera encerrada quando o programa terminar.
	defer db.Close()

	// ========================================================
	// ROTAS
	// ========================================================

	// GET /
	//
	// Retorna o Hello World.
	http.HandleFunc("GET /", helloWorld)

	// GET /alunos
	//
	// Lista todos os alunos.
	http.HandleFunc("GET /alunos", listarAlunos)

	// POST /alunos
	//
	// Cadastra um novo aluno.
	http.HandleFunc("POST /alunos", cadastrarAluno)

	// PUT /alunos/{codigo}
	//
	// Atualiza um aluno existente.
	http.HandleFunc("PUT /alunos/{codigo}", alterarAluno)

	// DELETE /alunos/{codigo}
	//
	// Remove um aluno.
	http.HandleFunc("DELETE /alunos/{codigo}", removerAluno)

	// ========================================================
	// PORTA
	// ========================================================

	// Definir a porta
	porta := "8080"

	// Exibe no terminal a porta utilizada.
	fmt.Printf("Servidor em execuÃ§Ã£o na porta: %s\n", porta)

	// ========================================================
	// INICIAR SERVIDOR
	// ========================================================

	http.ListenAndServe(
		":"+porta,
		corsMiddleware(http.DefaultServeMux),
	)
}
