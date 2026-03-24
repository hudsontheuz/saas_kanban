# Contrato de integração Front + Back

## O que o frontend agora consome

### Auth
- `POST /auth/register`
- `POST /auth/login`

### Team
- `POST /teams`
- `GET /teams/{id}`
- `POST /teams/{id}/members`
- `DELETE /teams/{id}/members/{memberId}`
- `POST /teams/{id}/leader/transfer`

### Project
- `POST /teams/{id}/projects`
- `GET /teams/{id}/projects`
- `POST /projects/{id}/close`
- `PUT /projects/{id}/settings`

### Task
- `GET /projects/{id}/tasks`
- `POST /projects/{id}/tasks`
- `POST /tasks/{id}/self-assign`
- `POST /tasks/{id}/pause`
- `POST /tasks/{id}/resume`
- `POST /tasks/{id}/in-review`
- `POST /tasks/{id}/approve`
- `POST /tasks/{id}/reject`

## Observações importantes

- O frontend foi ajustado para parar de depender de mocks e `localStorage` para dados de domínio.
- `localStorage` ficou apenas para token, usuário autenticado e ids atuais de equipe/projeto.
- Enquanto `GET /teams/{id}` e `GET /projects/{id}/tasks` não existirem no backend, a interface funciona de forma otimista após criação/ações na sessão atual, mas não consegue reconstruir tudo após refresh com fidelidade completa.
- O frontend está tolerante a respostas com campos em camelCase, snake_case ou nomes em português vindos do Go.

## Prioridade sugerida no backend

1. Expor `project` no delivery HTTP
2. Criar leituras:
   - `GET /teams/{id}`
   - `GET /teams/{id}/projects`
   - `GET /projects/{id}/tasks`
3. Criar gestão de equipe:
   - adicionar membro
   - remover membro
   - transferir liderança
4. Criar `PUT /projects/{id}/settings`

## Resultado esperado

Depois desses endpoints, o frontend já preparado nesta entrega passa a operar sem mock e com ciclo real de:
- login
- criação de equipe
- criação de projeto
- criação de tarefas
- self-assign
- pause/resume
- revisão
- aprovação/reprovação
- gestão de membros
- configurações do projeto
