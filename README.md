# saas_kanban

Sistema full stack de gestão de tarefas no modelo **Kanban**, desenvolvido com foco em **regras de negócio**, **manutenibilidade**, **testabilidade** e **organização arquitetural**.

O projeto foi construído como um **SaaS Kanban MVP**, com autenticação JWT, fluxo completo de tarefas e separação de responsabilidades inspirada em **Clean Architecture**, **SOLID** e conceitos iniciais de **DDD (Domain-Driven Design)**.

---

## 🚀 Como testar o sistema

### Pré-requisitos

Você precisa ter instalado na sua máquina:

- Git
- Docker
- Docker Compose

---

### 1. Clonar o repositório

```bash
git clone https://github.com/hudsontheuz/saas_kanban.git
cd saas_kanban
```

---

### 2. Criar o arquivo `.env`

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

No Windows PowerShell:

```powershell
copy .env.example .env
```

---

### 3. Subir o projeto

```bash
docker compose up --build -d
```

Esse comando irá:

- subir o PostgreSQL
- aguardar o banco ficar saudável
- rodar as migrations automaticamente
- subir o backend
- subir o frontend

---

### 4. Acessar a aplicação

**Frontend**  
`http://localhost:5173`

**Backend**  
`http://localhost:8080`

---

### 5. Primeiro uso

Ao subir pela primeira vez:

1. registre um usuário
2. faça login
3. crie uma equipe
4. crie um projeto
5. crie tarefas
6. teste o fluxo Kanban

---

### Persistência de dados

O banco roda em container, mas os dados ficam persistidos em um **volume Docker**.

**Parar o projeto sem apagar dados**
```bash
docker compose down
```

**Subir novamente mantendo os dados**
```bash
docker compose up -d
```

**Resetar tudo do zero**
```bash
docker compose down -v
```

Esse comando remove também o volume do banco.

---

## 🛠️ Tecnologias utilizadas

### Backend
- Go
- Chi
- GORM
- PostgreSQL
- JWT

### Frontend
- React
- Vite
- TypeScript
- Tailwind CSS
- shadcn/ui

### Infra
- Docker
- Docker Compose

---

## 🧠 Como o projeto foi pensado

Este projeto não foi pensado como um CRUD simples. A ideia principal foi construir um sistema de negócio pequeno, mas com regras explícitas e comportamento real.

O foco foi praticar:

- modelagem de domínio
- separação de responsabilidades
- regras de negócio no núcleo da aplicação
- arquitetura organizada para evolução
- fluxo full stack com backend, frontend e banco integrados

A estrutura segue uma abordagem inspirada em **Clean Architecture + DDD leve**, organizada em contextos e camadas.

### Contextos principais

- `auth`
- `user`
- `team`
- `project`
- `task`

### Camadas

- **Domain** → regras puras de negócio
- **Application** → casos de uso, DTOs e ports
- **Delivery** → handlers, rotas e middlewares HTTP
- **Infrastructure** → persistência e integrações concretas
- **Migrations** → versionamento do schema SQL
- **Tests** → validação de regras e fluxos

---

## 📋 Regras de negócio centrais

### Team
- uma equipe possui membros
- uma equipe pode ter liderança
- uma equipe pode ter **apenas 1 projeto ativo por vez**

### Project
- pertence a uma equipe
- possui configurações de fluxo e aprovação
- contém tarefas

### Task
- pertence a um projeto
- só é assumida oficialmente quando alguém faz **self-assign**
- cada usuário pode ter **apenas 1 task em Doing por vez**
- pode ser pausada e retomada
- se for reprovada em `InReview`, retorna para `ToDo`
- suporta **soft delete** para preservação histórica

---

## 🧩 Visão geral do sistema

O sistema permite:

- cadastro e login de usuários
- criação e gestão de equipes
- criação de projetos dentro de equipes
- criação e movimentação de tarefas em fluxo Kanban
- aprovação e reprovação de tarefas em revisão

### Fluxo de status das tasks

- `ToDo`
- `Doing`
- `InReview`
- `Done`

---

## 📁 Estrutura principal

```text
backend/
├── cmd/
├── delivery/
├── infrastructure/
├── internal/
├── migrations/
└── tests/

frontend/
├── src/
└── ...
```

---

## 🌐 Endpoints principais

### Registro
`POST /auth/register`

Exemplo de body:

```json
{
  "nome": "Hudson",
  "email": "hudson@teste.com",
  "senha": "123456"
}
```

### Login
`POST /auth/login`

Exemplo de body:

```json
{
  "email": "hudson@teste.com",
  "senha": "123456"
}
```

---

## 🧪 Como rodar os testes do backend

Dentro da pasta `backend`:

```bash
cd backend
go test ./...
```

---

## 📦 Objetivo deste repositório

Este projeto foi desenvolvido com foco em portfólio técnico, buscando demonstrar evolução em:

- backend orientado a regras
- separação por responsabilidades
- consistência de fluxo
- arquitetura para sistemas de negócio
- integração full stack

---

## ✅ Status

Projeto funcional e pronto para execução local com Docker Compose.

---

## Autor

**Hudson**
