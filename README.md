# 📘 API de Gerenciamento de Alunos

API RESTful desenvolvida em **Go (Golang)** com **PostgreSQL** para cadastro, listagem, edição e exclusão de registros de alunos, com cálculo automático de média e situação.

---

## 📁 Estrutura do Projeto

```text
projeto-alunos/
├── main.go          # Ponto de entrada da aplicação
├── models.go        # Definição da estrutura de dados
├── database.go      # Conexão com PostgreSQL e funções SQL
├── handlers.go      # Lógica das rotas / controladores
├── routes.go        # Registro de todas as rotas
├── middleware.go    # Configuração de CORS
├── go.mod           # Gerenciamento de dependências
├── go.sum           # Verificação de integridade
└── README.md        # Documentação do projeto
🚀 Tecnologias Utilizadas
Go 1.22+ — Linguagem de programação
PostgreSQL — Banco de dados relacional
github.com/lib/pq — Driver PostgreSQL para Go
github.com/google/uuid — Geração de identificadores únicos
⚙️ Configuração e Instalação
Pré-requisitos
Go instalado na versão 1.22 ou superior
Banco de dados PostgreSQL ativo
1. Instalar Dependências
go get github.com/lib/pq
go get github.com/google/uuid

Ou:

go mod tidy
2. Configurar Conexão com o Banco

A string de conexão está definida no arquivo database.go.

Exemplo:

dbURL := "postgresql://USUARIO:SENHA@HOST/BANCO?sslmode=require"

⚠️ Para uso local, altere para sslmode=disable no final da URL.

3. Executar o Projeto
go run .

Servidor iniciado em:

http://localhost:8080
📋 Endpoints da API
Método	Rota	Descrição
GET	/	Rota de teste / boas-vindas
GET	/alunos	Listar todos os alunos cadastrados
POST	/alunos	Cadastrar um novo aluno
PUT	/alunos/{codigo}	Atualizar dados de um aluno
DELETE	/alunos/{codigo}	Remover um aluno
📤 Exemplos de Requisição
Cadastrar Aluno — POST /alunos
{
  "nome": "Maria da Silva",
  "nota1": 8.5,
  "nota2": 7.0
}
Resposta
{
  "codigo": "a1b2c3d4-...",
  "nome": "Maria da Silva",
  "nota1": 8.5,
  "nota2": 7.0,
  "media": 7.75,
  "situacao": "Aprovado(a)"
}
Listar Alunos — GET /alunos

Retorna um array com todos os registros ordenados por nome.

Atualizar Aluno — PUT /alunos/{codigo}

Informe o código na URL e envie os dados atualizados no corpo da requisição.

Exemplo:

PUT /alunos/a1b2c3d4-...
{
  "nome": "Maria da Silva Souza",
  "nota1": 9.0,
  "nota2": 8.0
}

A média e a situação são recalculadas automaticamente.

Excluir Aluno — DELETE /alunos/{codigo}

Remove o aluno correspondente ao código informado.

Em caso de sucesso, retorna:

204 No Content
🧮 Regras de Negócio
Cálculo da Média
Média = (nota1 + nota2) / 2
Situação
Média	Situação
>= 7.0	Aprovado(a)
>= 5.0 e < 7.0	Em Recuperação
< 5.0	Reprovado(a)
Regras adicionais
O campo codigo é gerado automaticamente utilizando UUID.
A média é calculada automaticamente.
A situação é definida automaticamente de acordo com a média.
A tabela alunos é criada automaticamente na primeira execução.
🗄️ Estrutura da Tabela alunos
Coluna	Tipo	Descrição
codigo	VARCHAR(36)	Identificador único (PK)
nome	VARCHAR(100)	Nome do aluno
nota1	NUMERIC(4,2)	Primeira nota
nota2	NUMERIC(4,2)	Segunda nota
media	NUMERIC(4,2)	Média calculada
situacao	VARCHAR(20)	Situação do aluno
🔑 Resumo dos Arquivos
Arquivo	Responsabilidade
models.go	Define a struct Aluno com tags JSON
database.go	Conexão, criação de tabela e funções SQL
handlers.go	Lógica de cada rota, validação, cálculo e resposta
routes.go	Associação dos caminhos às funções correspondentes
middleware.go	Liberação de CORS para integração com frontend
main.go	Inicializa banco, carrega rotas e inicia servidor
🔄 Fluxo da Aplicação
Cliente
   │
   ▼
HTTP Request
   │
   ▼
Routes
   │
   ▼
Handlers
   │
   ├── Validação
   ├── Cálculo da média
   └── Definição da situação
   │
   ▼
Database
   │
   ▼
PostgreSQL
   │
   ▼
HTTP Response
🧪 Testando a API

A API pode ser testada utilizando:

Postman
Insomnia
Thunder Client
curl
Frontend HTML, CSS e JavaScript
Exemplo com curl
Cadastrar aluno
curl -X POST http://localhost:8080/alunos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Maria da Silva",
    "nota1": 8.5,
    "nota2": 7.0
  }'
Listar alunos
curl http://localhost:8080/alunos
🌐 CORS

O projeto possui configuração de CORS para permitir requisições de diferentes origens.

Atualmente, a configuração permite requisições de qualquer origem:

*

⚠️ Para ambientes de produção, recomenda-se restringir as origens permitidas aos domínios realmente utilizados pela aplicação.

🔐 Segurança

As informações de conexão com o banco de dados não devem ficar expostas diretamente no código-fonte.

Para ambientes de produção, recomenda-se utilizar variáveis de ambiente.

Exemplo:

DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
DB_SSLMODE

Dessa forma, informações sensíveis não ficam expostas no código ou no repositório.

📌 Observações
O servidor roda na porta 8080.
A tabela alunos é criada automaticamente.
O campo codigo utiliza UUID.
A média é calculada automaticamente.
A situação é determinada automaticamente pela média.
O PostgreSQL é utilizado como banco de dados.
Para conexão local, normalmente utiliza-se sslmode=disable.
Para bancos em nuvem, pode ser necessário utilizar sslmode=require.
👤 Autor

Projeto desenvolvido para fins acadêmicos e prática de desenvolvimento de APIs REST com Go e PostgreSQL.

📄 Licença

Este projeto foi desenvolvido para fins de estudo e aprendizado.