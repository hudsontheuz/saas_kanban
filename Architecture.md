# Arquitetura — saas_kanban

## Visão Geral

O backend do `saas_kanban` foi estruturado com foco em separação clara de responsabilidades, baixo acoplamento e facilidade de evolução.

A arquitetura é inspirada em **Clean Architecture** com um **DDD leve**, organizada por bounded contexts, aplicando princípios SOLID e mantendo o **domínio isolado** de detalhes externos (HTTP, banco, frameworks).

O objetivo é permitir evolução previsível: trocar persistência, trocar transporte, adicionar autenticação real, sem reescrever regras de negócio.

---

## Camadas

### 1. Domain (`internal/domain`)

Responsável por concentrar:

- Entidades
- Regras de negócio
- Invariantes
- Comportamentos do sistema

Características:

- Não depende de nenhuma outra camada
- Não conhece banco de dados
- Não conhece HTTP
- Não conhece frameworks

É o núcleo do sistema.

---

### 2. Application (`internal/application`)

Responsável por:

- Casos de uso (usecases)
- Orquestração do domínio
- Definição de interfaces (Ports)
- DTOs

Contém (por bounded context):

- `usecase/`
- `ports/`
- `dto/`

A camada Application depende do Domain, mas **não** conhece implementações concretas da Infrastructure.

---

### 3. Delivery (HTTP) (`delivery/http`)

Responsável por expor o sistema via HTTP, sem regras de negócio.

Contém:

- Router (chi)
- Handlers (controllers)
- Middlewares (ex.: fake auth)
- Mapeamento de erros para HTTP status

Características:

- Depende de Application (chama usecases)
- Não deve conter regra de negócio
- Não conhece detalhes de persistência (só interfaces/usecases)

---

### 4. Infrastructure (`infrastructure`)

Implementações concretas das interfaces definidas na Application.

Exemplos atuais:

- Persistência em memória (`infrastructure/persistence/memory`) — usada em testes/execução simples
- Persistência PostgreSQL via GORM (`infrastructure/persistence/gorm`)

Características:

- Depende de Application e Domain
- Nunca o contrário

---

### 5. Migrations (`migrations`)

Migrações SQL controlam o schema do Postgres.

O schema já inclui constraints que reforçam regras críticas do sistema, como:

- **1 projeto ACTIVE por equipe**
- **1 task DOING não pausada por usuário**

---

### 6. Tests

Testes focados principalmente em:

- Validação de comportamento do domínio
- Garantia das regras críticas
- Usecases (quando aplicável)

---

## Estrutura de Pastas

```text
cmd/
  api/
    main.go

delivery/
  http/
    router.go
    handlers/
      task_handler.go
    middleware/
      (fake auth, etc)

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
      usecase/
    task/
      dto/
      ports/
      usecase/
    team/
      dto/
      ports/
      usecase/

infrastructure/
  persistence/
    memory/
    gorm/
      db.go
      model/
      repo/

migrations/
  001_start_system.up.sql
  001_start_system.down.sql

tests/