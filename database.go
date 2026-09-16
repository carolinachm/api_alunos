package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// db Pool de conexões com o banco
var db *sql.DB

// InitDB estabelece a conexão e cria a tabela se não existir
func InitDB() {
	// URL completa com sslmode
	dbURL := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erro ao configurar conexão: %v", err)
	}

	// Verifica se a conexão está ativa
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}
	fmt.Println("Conectado ao PostgreSQL com sucesso!")

	// Cria tabela "alunos"
	criarTabela := `
	CREATE TABLE IF NOT EXISTS alunos (
		codigo VARCHAR(36) PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		nota1 NUMERIC(4,2) NOT NULL,
		nota2 NUMERIC(4,2) NOT NULL,
		media NUMERIC(4,2) NOT NULL,
		situacao VARCHAR(20) NOT NULL
	);`

	_, err = db.Exec(criarTabela)
	if err != nil {
		log.Fatalf("❌ Erro ao criar tabela: %v", err)
	}
	fmt.Println("✅ Tabela 'alunos' pronta!")
}

// ============================================================
// FUNÇÕES COMANDOS SQL
// ============================================================

// SelecionarTodos → SELECT * FROM alunos
func SelecionarTodos() ([]Aluno, error) {
	linhas, err := db.Query(`
		SELECT codigo, nome, nota1, nota2, media, situacao
		FROM alunos
		ORDER BY nome
	`)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var lista []Aluno
	for linhas.Next() {
		var a Aluno
		err := linhas.Scan(&a.Codigo, &a.Nome, &a.Nota1, &a.Nota2, &a.Media, &a.Situacao)
		if err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	return lista, nil
}

// Inserir → INSERT INTO alunos
func Inserir(aluno Aluno) error {
	_, err := db.Exec(`
		INSERT INTO alunos (codigo, nome, nota1, nota2, media, situacao)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, aluno.Codigo, aluno.Nome, aluno.Nota1, aluno.Nota2, aluno.Media, aluno.Situacao)
	return err
}

// Atualizar → UPDATE alunos SET ... WHERE codigo
func Atualizar(aluno Aluno) error {
	resultado, err := db.Exec(`
		UPDATE alunos
		SET nome = $1, nota1 = $2, nota2 = $3, media = $4, situacao = $5
		WHERE codigo = $6
	`, aluno.Nome, aluno.Nota1, aluno.Nota2, aluno.Media, aluno.Situacao, aluno.Codigo)

	if err != nil {
		return err
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return fmt.Errorf("código não encontrado")
	}
	return nil
}

// Excluir → DELETE FROM alunos WHERE codigo
func Excluir(codigo string) error {
	resultado, err := db.Exec(`
		DELETE FROM alunos WHERE codigo = $1
	`, codigo)

	if err != nil {
		return err
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return fmt.Errorf("código não encontrado")
	}
	return nil
}

// GetDB retorna a instância da conexão
func GetDB() *sql.DB {
	return db
}