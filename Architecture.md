# Arquitetura — saas_kanban

## Visão Geral

O backend do `saas_kanban` foi estruturado com foco em separação clara de responsabilidades, baixo acoplamento e facilidade de evolução.

A arquitetura é inspirada em Clean Architecture, organizada por bounded contexts, com aplicação prática de princípios SOLID.

O objetivo principal é manter o domínio isolado de detalhes externos como banco de dados, transporte HTTP ou frameworks.

---

## Camadas

### 1. Domain

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

### 2. Application

Responsável por:

- Casos de uso
- Orquestração do domínio
- Definição de interfaces (Ports)
- DTOs

Contém:

- `usecase/`
- `ports/`
- `dto/`

A camada Application depende do Domain, mas não conhece implementações concretas da Infrastructure.

---

### 3. Infrastructure

Contém implementações concretas das interfaces definidas na Application.

Exemplo atual:

- Persistência em memória

Essa camada depende de Application e Domain, mas nunca o contrário.

---

### 4. Tests

Testes focados principalmente em:

- Validação de comportamento do domínio
- Garantia das regras críticas

---

## Estrutura de Pastas


---

## Organização por Bounded Context

O projeto está dividido em contextos independentes:

- `task`
- `project`
- `team`

Cada contexto possui sua própria organização interna dentro das camadas.

No estágio atual (MVP), pode haver dependência entre contextos no nível de Application.

Refinamentos futuros podem introduzir interfaces leitoras específicas para reduzir acoplamento direto entre contexts.

---

## Dependências Entre Camadas

Fluxo permitido:

Infrastructure → Application → Domain

Fluxo não permitido:

Domain → Application  
Domain → Infrastructure  
Application → Infrastructure (implementações concretas)

O domínio permanece isolado.

---

## Decisões Arquiteturais Atuais

- Persistência em memória para manter foco nas regras de negócio
- API HTTP ainda não implementada
- Banco de dados real e Docker planejados para próxima etapa
- Estrutura modular definida desde o início para evitar refatorações estruturais futuras

---

## Diretrizes de Código

- `usecase` permanece no singular (padronizado)
- `ports` contém apenas interfaces
- Implementações concretas vivem exclusivamente na Infrastructure
- Entidades concentram comportamento, não apenas dados
- Application orquestra, Domain decide

---

## Objetivo Evolutivo

A arquitetura foi pensada para permitir:

- Adição futura de camada HTTP sem alterar o domínio
- Troca de persistência (memória → PostgreSQL) sem impactar regras
- Crescimento do sistema mantendo previsibilidade estrutural
