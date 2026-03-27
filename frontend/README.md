# SaaS Kanban Frontend

Frontend refeito do zero para o projeto **saas_kanban**, seguindo a stack combinada:

- **Vite**
- **React 18**
- **TypeScript**
- **Tailwind CSS v3**
- **shadcn/ui style components** (com Radix + utilitários no padrão shadcn)
- **React Router DOM**
- **React Hook Form + Zod**
- **Axios**
- **Lucide React**

## Objetivo desta versão

Esta base substitui a mistura anterior de:

- JavaScript + TypeScript
- Material Tailwind + outras bibliotecas visuais
- estrutura herdada de template dashboard

O foco aqui é uma base **limpa, tipada, organizada por feature e alinhada ao backend em Go**.

---

## Arquitetura adotada

A estrutura foi organizada por **feature/domain no frontend**, separando layout compartilhado, componentes de UI e módulos de negócio.

```text
src/
  app/
    globals.css
    router.tsx

  components/
    layout/
    shared/
    ui/

  features/
    auth/
      api/
      components/
      pages/
      schemas/
      types/

    dashboard/
      pages/

    project/
      components/
      pages/
      types/

    task/
      api/
      components/
      pages/
      types/

    team/
      components/
      pages/
      types/

    settings/
      pages/

  hooks/
  lib/
  types/
  main.tsx
```

## Como essa arquitetura conversa com o backend

O backend atual já foi pensado em bounded contexts como:

- auth
- project
- task
- team
- user

Por isso o frontend foi desenhado para seguir a mesma direção conceitual. Isso facilita:

- evolução por domínio
- menor acoplamento
- manutenção mais simples
- crescimento mais previsível

---

## Rotas prontas

As rotas do frontend já estão funcionando:

- `/sign-in`
- `/sign-up`
- `/`
- `/projects`
- `/tasks`
- `/team`
- `/settings`

### Fluxos implementados na interface

- autenticação com formulário de login
- cadastro com formulário tipado
- proteção de rotas privadas
- dashboard inicial
- tela de projeto ativo
- quadro de tarefas com colunas:
  - ToDo
  - Doing
  - In Review
  - Done
- limite visual/regra local de **1 tarefa em Doing por usuário**
- pausa e retomada de tarefa
- tela de equipe com convite, remoção e transferência de liderança
- tela de configurações do projeto com switches para regras do fluxo

---

## Integração com backend

Hoje o backend enviado já expõe de forma clara principalmente:

- `POST /auth/register`
- `POST /auth/login`
- `POST /tasks/{id}/self-assign`

Por isso este frontend já tenta integrar autenticação com o backend via `VITE_API_URL`.

### Importante

Se o backend ainda não estiver rodando ou se algum endpoint ainda não existir, o frontend entra em **modo demo local** para não travar a navegação. Isso foi feito porque o backend ainda está em evolução.

Ou seja:

- **auth** tenta usar API real primeiro
- se falhar, cria sessão local demo
- project/task/team/settings usam estado local persistido em `localStorage` como base temporária

Isso permite continuar a construção do frontend sem depender de toda a API pronta.

---

## Como rodar

### 1. Instalar dependências

```bash
npm install
```

### 2. Configurar ambiente

Copie o arquivo `.env.example` para `.env`:

```bash
cp .env.example .env
```

### 3. Rodar em desenvolvimento

```bash
npm run dev
```

### 4. Build de produção

```bash
npm run build
```

---

## Observações importantes

### 1. Esta versão já remove a base conceitual do template antigo

Aqui não existe mais dependência de:

- Material Tailwind
- Heroicons do template antigo
- widgets/charts/cards herdados do dashboard pronto

### 2. Componentes no padrão shadcn/ui

Para manter a stack escolhida, os componentes foram escritos no estilo do ecossistema shadcn/ui:

- `Button`
- `Input`
- `Card`
- `Badge`
- `Dialog`
- `Switch`
- `Tabs`

### 3. Próximos passos naturais

Quando o backend avançar, a evolução recomendada é:

- trocar stores locais por chamadas reais de API
- criar camada de query/mutation por feature
- mapear DTOs exatamente conforme contratos do backend
- ligar criação de projeto, criação de tarefa, revisão e aprovação reais
- adicionar feedback de loading/error por operação

---

## Resumo

Esse frontend foi montado para ser:

- mais limpo que a base anterior
- totalmente em **TypeScript**
- consistente com **Tailwind + shadcn/ui**
- preparado para crescer junto do backend
- simples o bastante para você manter e evoluir
