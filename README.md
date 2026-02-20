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

Este projeto segue uma abordagem inspirada em Clean Architecture.

### Camadas:

- **Domain** → Regras puras de negócio (sem dependência externa)
- **Application** → Casos de uso, DTOs e Ports (interfaces)
- **Infrastructure** → Implementações concretas (ex: persistência em memória)
- **Tests** → Validação de comportamento

### Estrutura de Pastas

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

🧪 Como rodar os testes

go test ./...