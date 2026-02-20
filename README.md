# SaaS Kanban (Backend)

Backend em **Go** para um sistema de gestão de tarefas no modelo Kanban, desenvolvido com foco em **clareza de domínio**, **manutenibilidade** e **testabilidade**.

Este projeto representa uma evolução consciente de arquitetura: sair de um CRUD funcional para uma estrutura orientada a regras de negócio, aplicando **Clean Architecture**, princípios **SOLID** e conceitos iniciais de **DDD** (Domain-Driven Design).

---

## 🎯 Objetivo do Projeto

Construir um backend sólido e bem estruturado que:

- Separe claramente domínio e infraestrutura
- Modele comportamento (não apenas tabelas)
- Permita evolução sem alto acoplamento
- Seja fácil de testar
- Demonstre maturidade técnica em entrevistas

Este não é apenas um projeto que “funciona”, mas um projeto pensado para **evoluir com consistência**.

---

## 📚 Regras de Negócio (MVP)

### 📌 Kanban (Fluxo de Status)

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
- Pode ter sugestão de responsável (SelectedAssignee).
- Só é assumida oficialmente quando alguém realiza **self-assign**.
- Cada usuário pode ter **apenas 1 task em Doing por vez**.
- Se reprovada em `InReview`, retorna para `ToDo`.
- Suporta conceito de soft delete (preservação histórica).

---

## 🧠 Arquitetura

Este projeto segue uma abordagem inspirada em **Clean Architecture**:

- **Domain** → Regras puras de negócio (não depende de banco ou framework)
- **Application** → Casos de uso, DTOs e Ports (interfaces)
- **Infrastructure** → Implementações concretas (ex: persistência em memória, futuro Postgres)
- **Tests** → Foco em validar comportamento do domínio

### Regra central:

> O domínio não depende de nada.  
> O restante do sistema depende do domínio.

---

## 🗂 Estrutura do Projeto

```text
internal/
  domain/
    project/
    task/
    team/
    shared/

  application/
    project/
      dto/
      ports/
      usecases/
    task/
      dto/
      ports/
      usecases/
    team/
      dto/
      ports/
      usecases/

  infrastructure/
    persistence/
      memory/

tests/


🧪 Testes

Rodar todos os testes:

go test ./...
