# saas_kanban

Backend em Go para um sistema de gestão de tarefas no modelo Kanban, construído com foco em **regras de negócio**, **manutenibilidade** e **testabilidade**.

Este projeto representa uma evolução consciente de arquitetura: sair de um CRUD funcional para uma estrutura orientada a domínio, aplicando princípios de **Clean Architecture**, **SOLID** e conceitos iniciais de **DDD (Domain-Driven Design)**.

---

## 🎯 Visão do MVP

Sistema Kanban com regras explícitas e comportamento modelado no domínio.

### 📌 Fluxo de Status

- `ToDo`
- `Doing`
- `InReview`
- `Done`

### 👥 Team (Equipe)

- Uma equipe possui membros.
- Pode haver líder(es) com permissões administrativas.
- Uma equipe pode ter **apenas 1 projeto ativo por vez**.

### 📁 Project (Projeto)

- Pertence a uma equipe.
- Contém tarefas.
- Pode ter configurações que influenciam fluxo e aprovação.

### 📝 Task (Tarefa)

- Pertence a um projeto.
- Pode ter sugestão de responsável (**SelectedAssignee**).
- Só é assumida oficialmente quando alguém realiza **self-assign**.
- Cada usuário pode ter **apenas 1 task em Doing por vez**.
- Se reprovada em `InReview`, retorna para `ToDo`.
- Suporta conceito de soft delete (preservação histórica).

---

## 🧠 Arquitetura (Bounded Contexts)

Este projeto segue uma abordagem inspirada em Clean Architecture + DDD leve.

### Camadas

- **Domain** → Regras puras de negócio (sem dependência externa)
- **Application** → Casos de uso, DTOs e Ports (interfaces)
- **Delivery (HTTP)** → Router/handlers/middlewares (chi)
- **Infrastructure** → Implementações concretas (memória, PostgreSQL via GORM)
- **Migrations** → Schema SQL versionado (golang-migrate)
- **Tests** → Validação de comportamento e regras

---

## ✅ Status atual

- PostgreSQL via Docker (porta **5434**)
- Migrations com `golang-migrate`
- Infra GORM criada (`infrastructure/persistence/gorm`)
- Delivery HTTP preparada (`delivery/http`)
- Primeira rota preparada: `POST /tasks/{id}/self-assign`
- Autenticação real ainda não existe: será usado **fake auth** via middleware (dev)

---

## 🐳 Subir Postgres (Docker)

Exemplo (ajuste se você já tiver seu `docker-compose.yml`):

```bash
docker run --name saas_kanban_pg \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=saas_kanban \
  -p 5434:5432 \
  -d postgres:16
🗄️ Rodar migrations (golang-migrate)

Exemplo:

migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable" up

Rollback:

migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable" down 1
▶️ Rodar a API

A API usa DB_URL:

export DB_URL="postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable"
go run ./cmd/api
🌐 Endpoints (inicial)
POST /tasks/{id}/self-assign

Assume a task para o usuário autenticado (via fake auth).

Header (dev):

X-User-Id: <id>

Exemplo:

curl -X POST \
  -H "X-User-Id: 1" \
  http://localhost:8080/tasks/10/self-assign

Resposta esperada:

{"ok":true}

Erros comuns:

404 se task não existe

409 se o usuário já tem uma task em Doing não pausada

400 para violações de regra de negócio

🧪 Como rodar os testes
go test ./...