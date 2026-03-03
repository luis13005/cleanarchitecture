# 🏆 Clean Architecture - Order System

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![GraphQL](https://img.shields.io/badge/-GraphQL-E10098?style=for-the-badge&logo=graphql&logoColor=white)
![gRPC](https://img.shields.io/badge/-gRPC-4285F4?style=for-the-badge&logo=google&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white)

Sistema de gerenciamento de ordens em **Go**, seguindo **Clean Architecture**. Expõe a mesma lógica de negócio através de três interfaces: **REST**, **gRPC** e **GraphQL**.

---

## 📁 Estrutura do Projeto

```
cleanarchitecture/
├── cmd/
│   └── ordersystem/
│       ├── main.go          # Entrypoint da aplicação
│       ├── wire.go          # Definição das dependências (Wire)
│       ├── wire_gen.go      # Código gerado automaticamente pelo Wire
│       └── .env             # Variáveis de ambiente (fica aqui!)
├── configs/
│   └── config.go            # Leitura do .env com Viper
├── internal/
│   ├── entity/              # Regras de negócio puras (Domain)
│   ├── usecase/             # Casos de uso (orquestração)
│   ├── infra/
│   │   ├── database/        # Repositórios (PostgreSQL)
│   │   ├── grpc/
│   │   │   ├── protofiles/  # Arquivos .proto
│   │   │   ├── pb/          # Código gerado pelo protoc
│   │   │   └── service/     # Implementação dos serviços gRPC
│   │   ├── graph/
│   │   │   ├── schema.graphqls   # Schema GraphQL
│   │   │   ├── resolver.go       # Resolvers
│   │   │   └── model/            # Tipos gerados pelo gqlgen
│   │   └── web/             # Handlers REST (HTTP)
│   └── event/               # Eventos de domínio (RabbitMQ)
├── pkg/
│   └── events/              # Utilitários de eventos
├── sql/
│   └── migrations/          # Arquivos de migração SQL
├── docker-compose.yaml
├── Makefile
└── go.mod
```

---

## 🔧 Pré-requisitos

| Ferramenta | Versão | Uso |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.21+ | Linguagem principal |
| [Docker + Compose](https://docs.docker.com/get-docker/) | Latest | PostgreSQL + RabbitMQ |
| [migrate](https://github.com/golang-migrate/migrate) | Latest | Migrations do banco |
| [protoc](https://grpc.io/docs/protoc-installation/) | Latest | Gerar código gRPC |
| [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf) | Latest | Plugin Go para protoc |
| [wire](https://github.com/google/wire) | Latest | Injeção de dependências |
| [Evans](https://github.com/ktr0731/evans) | Latest | Cliente gRPC interativo |
| WSL | - | Rodar migrate no Windows |

### Instalando as ferramentas Go

```bash
# Wire - Injeção de dependências
go install github.com/google/wire/cmd/wire@latest

# migrate - Migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# gqlgen - GraphQL
go get github.com/99designs/gqlgen
```

---

> ⚠️ **Windows:** O Go busca o `.env` relativo ao diretório onde o comando é executado. Sempre rode de dentro de `cmd/ordersystem/`, ou use o Makefile.

### 1. Subir infraestrutura com Docker

```bash
docker-compose up -d
```

Isso sobe:
- **PostgreSQL** na porta `5432`
- **RabbitMQ** na porta `5672` (management UI: `http://localhost:15672`)

### 3. Verificar containers

```bash
docker-compose ps
```

---

## 🗄️ Migrations (Banco de Dados)

As migrations ficam em `sql/migrations/` e seguem o padrão:

```
sql/migrations/
├── 000001_create_orders_table.up.sql
└── 000001_create_orders_table.down.sql
```

### Rodar migrations (WSL no Windows)

```bash
# Subir todas as migrations
wsl migrate -path=sql/migrations -database "postgresql://admin:admin@localhost:5432/faculdade?sslmode=disable" -verbose up
```

### Criar nova migration

```bash
wsl migrate create -ext sql -dir sql/migrations -seq nome_da_migration
```

---

## 💉 Injeção de Dependências (Wire)

O [Wire](https://github.com/google/wire) gera automaticamente o código de inicialização das dependências.

### Como funciona

- `wire.go` — você define os `Provider` e o `Injector`
- `wire_gen.go` — gerado automaticamente, **nunca edite manualmente**

### Quando regenerar

Sempre que adicionar ou alterar uma dependência (novo usecase, novo repositório, etc.):

```bash
# Entrar na pasta e rodar o wire
cd cmd/ordersystem
wire
```

---

## 🔵 gRPC

### Estrutura dos arquivos

```
internal/infra/grpc/
├── protofiles/
│   └── order.proto          # Definição do contrato
├── pb/
│   ├── order.pb.go          # Gerado pelo protoc
│   └── order_grpc.pb.go     # Gerado pelo protoc
└── service/
    └── order_service.go     # Implementação do servidor gRPC
```

### Gerar o código Go a partir do .proto

```bash
# Pelo Makefile
make gen

# Ou manualmente
protoc --proto_path=internal/infra/grpc/protofiles \
       internal/infra/grpc/protofiles/*.proto \
       --go_out=internal/infra/grpc/pb \
       --go-grpc_out=internal/infra/grpc/pb \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative
```

> ⚠️ Rode **sempre que alterar o `.proto`**.

Os arquivos gerados em `pb/` não devem ser editados manualmente.

### Testando com Evans (cliente interativo)

```bash
# Pelo Makefile (via Docker)
make evans
```

---

## 🟣 GraphQL

### Estrutura dos arquivos

```
internal/infra/graph/
├── schema.graphqls          # Definição do schema (você escreve)
├── gqlgen.yml               # Configuração do gqlgen
├── generated.go             # Código gerado (não edite)
├── resolver.go              # Struct Resolver + dependências
├── schema.resolvers.go      # Assinaturas geradas; você implementa o corpo
└── model/
    └── models_gen.go        # Tipos Go gerados automaticamente
```

### Gerar o código GraphQL

```bash
# Pelo Makefile
make graph

# Ou manualmente
go run github.com/99designs/gqlgen generate
```

> ⚠️ Rode **sempre que alterar o `schema.graphqls`**.

### Implementando os resolvers

Após o generate, `schema.resolvers.go` terá as assinaturas. Você implementa apenas o corpo:

```go
// internal/infra/graph/resolver.go
type Resolver struct {
    CreateOrderUseCase usecase.CreateOrderUseCase
}
```

### Testando com o Playground

Acesse `http://localhost:8080` no browser e execute:

```graphql
# Criar uma ordem
mutation {
  createOrder(input: {
    id: "order-1"
    price: 100.0
    tax: 0.1
  }) {
    id
    price
    tax
    finalPrice
  }
}

# Listar todas as ordens
query {
  orders {
    id
    price
    tax
    finalPrice
  }
}
```

---

## 🌐 REST

A API REST roda na porta `8000`.

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/order` | Criar uma nova ordem |
| `GET` | `/order/list` | Listar todas as ordens |

---

## ▶️ Executando o Projeto

### Passo a passo completo

```bash
# 1. Subir infraestrutura
docker-compose up -d

# 2. Rodar migrations
make migrate-up

# 3. Rodar a aplicação (sempre de dentro da pasta!)
cd cmd/ordersystem
go run main.go wire_gen.go
```

Ao subir, você verá:

```
Starting web server on port :8000
Starting gRPC server on port 50051
Starting GraphQL server on port 8080
```

---

## 🛠️ Makefile — Referência Completa

```makefile
# Gera o código gRPC a partir dos .proto
gen:
	protoc --proto_path=internal\infra\grpc\protofiles \
	       internal/infra/grpc/protofiles/*.proto \
	       --go_out=internal/infra/grpc/pb \
	       --go-grpc_out=internal/infra/grpc/pb \
	       --go_opt=paths=source_relative \
	       --go-grpc_opt=paths=source_relative

# Gera o código GraphQL a partir do schema
graph:
	go run github.com/99designs/gqlgen generate

# Abre o Evans (cliente gRPC interativo)
evans:
	docker run --rm -it \
	  -v "C:/Users/LUISFP/go/src/luisfp/pos/cleanarchitecture:/mount:ro" \
	  ghcr.io/ktr0731/evans:latest \
	  --path ./internal/infra/grpc/protofiles/ \
	  --proto order.proto \
	  --host host.docker.internal \
	  --port 50051 repl
```

---

## 🔄 Fluxo de Desenvolvimento


### Resumo dos comandos de geração de código

| Situação | Comando |
|---|---|
| Alterou o `.proto` | `make gen` |
| Alterou o `schema.graphqls` | `make graph` |
| Alterou dependências (Wire) | `wire` |

---

## 📦 Dependências Principais

```go
require (
    github.com/google/wire                      // Injeção de dependências
    github.com/99designs/gqlgen                 // GraphQL
    google.golang.org/grpc                      // gRPC
    google.golang.org/protobuf                  // Protobuf
    github.com/spf13/viper                      // Leitura de configuração (.env)
    github.com/golang-migrate/migrate/v4        // Migrations
    github.com/streadway/amqp                   // RabbitMQ
    github.com/lib/pq                           // Driver PostgreSQL
)
```