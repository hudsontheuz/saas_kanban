# Arquitetura — saas_kanban

## Visão geral

O `saas_kanban` está organizado em **módulos verticais (bounded contexts)**, mantendo as regras de negócio no domínio e deixando HTTP/DB como detalhes externos.

A ideia é conseguir evoluir o projeto com segurança:

- regras e invariantes ficam no **Domain**
- orquestração fica na **Application (usecases/ports/dto)**
- HTTP fica na **Delivery**
- persistência e integrações ficam na **Infrastructure**

## Estrutura

```text
SAAS_KANBAN/
├── cmd/
│   └── api/
│       └── main.go                # wiring (injeção)
├── delivery/
│   └── http/
│       └── router.go              # monta rotas chamando Register() de cada módulo
├── infrastructure/
│   └── persistence/
│       └── gorm/
│           └── db.go              # conexão/tx comum (infra “plataforma”)
├── internal/
│   ├── shared/
│   │   ├── errors/
│   │   ├── ids/
│   │   └── httpx/
│   ├── auth/
│   │   ├── delivery/http/
│   │   │   ├── authctx/
│   │   │   └── middleware/
│   │   ├── application/{dto,ports,usecase}/
│   │   ├── domain/
│   │   └── infrastructure/jwt/
│   ├── team/
│   │   ├── delivery/http/
│   │   ├── application/{dto,ports,usecase}/
│   │   ├── domain/
│   │   └── infrastructure/persistence/{gorm,memory}/
│   ├── project/
│   │   ├── delivery/http/
│   │   ├── application/{dto,ports,usecase}/
│   │   ├── domain/
│   │   └── infrastructure/persistence/{gorm,memory}/
│   └── task/
│       ├── delivery/http/
│       ├── application/{dto,ports,usecase}/
│       ├── domain/
│       └── infrastructure/persistence/{gorm,memory}/
├── migrations/
└── tests/
```

## Convenções

- **Router principal** (`delivery/http/router.go`) chama `Register()` de cada módulo.
- Cada módulo pode expor `Register()` no seu `delivery/http` para montar rotas e aplicar middlewares.
- Implementações concretas de repos ficam em `infrastructure/persistence/...` dentro do módulo.
- Coisas realmente compartilhadas ficam em `internal/shared/*`.
