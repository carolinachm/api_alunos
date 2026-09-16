package main

// Aluno representa os dados de um aluno no sistema
type Aluno struct {
	Codigo   string  `json:"codigo"`   // Identificador único
	Nome     string  `json:"nome"`     // Nome completo
	Nota1    float64 `json:"nota1"`    // Primeira nota
	Nota2    float64 `json:"nota2"`    // Segunda nota
	Media    float64 `json:"media"`    // Média calculada
	Situacao string  `json:"situacao"` // Situação: Aprovado, Recuperação, Reprovado
}