# Domain Rules — Task Management SaaS

Este documento descreve as regras de negócio do sistema.

O objetivo é centralizar a lógica de domínio e manter o sistema consistente, seguindo princípios de DDD, Clean Architecture e SOLID.

---

# Ubiquitous Language

## User

* Possui UserID único
* Pode participar de múltiplas Teams
* Pode ter no máximo 1 Task em Doing ativa no sistema inteiro

---

## Team

* Grupo de usuários
* Possui 1 Leader
* Pode ter múltiplos membros
* Pode ter apenas 1 Project ativo

### Regras

* Leader adiciona membros
* Leader remove membros
* Member pode sair da equipe

---

## Project

* Pertence a uma Team
* Pode estar Active ou Closed
* Apenas 1 Project ativo por Team

### Regras

* Se Project estiver Closed, nenhuma Task pode mudar de estado

---

# Task

Uma Task possui:

* Status
* Assignee
* SelectedAssignee
* Outcome
* Pause
* Soft delete

---

# Status Flow

Estados possíveis:

* ToDo
* Doing
* InReview
* Done

Fluxo obrigatório:

ToDo → Doing → InReview → Done

Não pode pular etapas.

---

# Assignment Rules

## SelectedAssignee

* Sugestão feita por Leader
* Não atribui automaticamente

## Assignee

* Definido via SelfAssign
* Apenas 1 Assignee por Task

---

# Self Assign

* Task deve estar em ToDo
* Usuário vira Assignee
* Task muda para Doing

---

# Global Work Limit

* Usuário pode ter apenas 1 Task em Doing ativa
* Para iniciar outra deve pausar a atual

---

# Pause / Resume

## Pause

* Apenas em Doing
* Não altera status
* Não conta como ativa

## Resume

* Apenas se não houver outra Doing ativa

---

# Review Rules

## Approve

InReview → Done

Outcome = Approved

## Reject em InReview

InReview → Doing

---

# Reject em ToDo

ToDo → Done

Outcome = Rejected

---

# Done

Estado final da Task

Outcome:

* Approved
* Rejected

---

# Soft Delete

* Não remove do banco
* Marca DeletedAt
* Marca DeletedBy
* Não aparece nas consultas padrão

---

# Invariantes de Domínio

* Task em Doing deve ter Assignee
* Task em InReview deve ter Assignee
* Task em Done Approved deve ter Assignee
* Task em Done Rejected pode não ter Assignee
* Usuário só pode ter 1 Doing ativa

---

# Objetivo do Domínio

O domínio foi projetado para:

* Garantir consistência de regras
* Facilitar manutenção
* Permitir evolução do sistema
* Separar regras de negócio da infraestrutura

Este projeto segue:

* Domain Driven Design (DDD)
* Clean Architecture
* SOLID Principles
