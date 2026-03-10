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
- **Migrations** → Schema SQL versionado
- **Tests** → Validação de comportamento e regras

### Contextos principais

- **user** → identidade global do sistema
- **auth** → cadastro, login, hash de senha, emissão e validação de JWT
- **team** → equipe e liderança
- **project** → projeto e configurações
- **task** → fluxo das tarefas e invariantes

---

## ✅ Status atual

- PostgreSQL via Docker (porta **5434**)
- Migrations SQL criadas
- Infra GORM criada (`infrastructure/persistence/gorm`)
- Delivery HTTP preparado (`delivery/http`)
- Contexto `user` separado de `team`
- Autenticação com JWT (`register` e `login`)
- Rota protegida: `POST /tasks/{id}/self-assign`

---

## 🐳 Subir Postgres (Docker)

Exemplo:

```bash
docker run --name saas_kanban_pg \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=saas_kanban \
  -p 5434:5432 \
  -d postgres:16
```

## 🗄️ Rodar migrations

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable" up
```

Rollback:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable" down 1
```

## ▶️ Rodar a API

```bash
export DB_URL="postgres://postgres:postgres@localhost:5434/saas_kanban?sslmode=disable"
export JWT_SECRET="dev-secret"
export JWT_ISSUER="saas_kanban"
go run ./cmd/api
```

---

## 🌐 Endpoints atuais

### Registrar usuário

`POST /auth/register`

Body:

```json
{
  "nome": "Matheus",
  "email": "matheus@teste.com",
  "senha": "123456"
}
```

### Login

`POST /auth/login`

Body:

```json
{
  "email": "matheus@teste.com",
  "senha": "123456"
}
```

### Self assign protegido por JWT

`POST /tasks/{id}/self-assign`

Header:

```text
Authorization: Bearer <token>
```

Exemplo:

```bash
curl -X POST \
  -H "Authorization: Bearer SEU_TOKEN" \
  http://localhost:8080/tasks/10/self-assign
```

Resposta esperada:

```json
{"ok":true}
```

Erros comuns:

- `404` se a task não existe
- `409` se o usuário já tem uma task em Doing não pausada
- `400` para violações de regra de negócio
- `401` para token ausente, inválido ou expirado

---

## 🧪 Como rodar os testes

```bash
go test ./...
```
