BEGIN;

CREATE TABLE IF NOT EXISTS usuario (
  id          BIGSERIAL PRIMARY KEY,
  nome        VARCHAR(255) NOT NULL,
  email       VARCHAR(255) NOT NULL UNIQUE,
  senha_hash  TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS equipe (
  id               BIGSERIAL PRIMARY KEY,
  nome             VARCHAR(255) NOT NULL,
  lider_usuario_id BIGINT NOT NULL,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_equipe_lider
    FOREIGN KEY (lider_usuario_id)
    REFERENCES usuario(id)
);

CREATE TABLE IF NOT EXISTS equipe_membro (
  id          BIGSERIAL PRIMARY KEY,
  equipe_id   BIGINT NOT NULL,
  usuario_id  BIGINT NOT NULL,
  role        VARCHAR(50) NOT NULL DEFAULT 'MEMBER',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT fk_equipe_membro_equipe
    FOREIGN KEY (equipe_id) REFERENCES equipe(id) ON DELETE CASCADE,

  CONSTRAINT fk_equipe_membro_usuario
    FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE
);

-- regra: não entra 2x na mesma equipe
CREATE UNIQUE INDEX IF NOT EXISTS uq_equipe_membro_equipe_usuario
  ON equipe_membro(equipe_id, usuario_id);

CREATE TABLE IF NOT EXISTS projeto (
  id                             BIGSERIAL PRIMARY KEY,
  equipe_id                       BIGINT NOT NULL,
  nome                            VARCHAR(255) NOT NULL,

  status                          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  permitir_soltar_doing           BOOLEAN NOT NULL DEFAULT false,

  fechado_em                      TIMESTAMPTZ NULL,

  created_at                      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at                      TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at                      TIMESTAMPTZ NULL,

  CONSTRAINT fk_projeto_equipe
    FOREIGN KEY (equipe_id) REFERENCES equipe(id) ON DELETE CASCADE,

  CONSTRAINT ck_projeto_status
    CHECK (status IN ('ACTIVE', 'CLOSED'))
);

-- regra: 1 projeto ACTIVE por equipe (ignorando soft delete)
CREATE UNIQUE INDEX IF NOT EXISTS uq_projeto_um_ativo_por_equipe
  ON projeto(equipe_id)
  WHERE status = 'ACTIVE' AND deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tarefa (
  id                  BIGSERIAL PRIMARY KEY,
  projeto_id          BIGINT NOT NULL,

  titulo              VARCHAR(255) NOT NULL,
  descricao           TEXT NULL,

  status              VARCHAR(20) NOT NULL DEFAULT 'TODO',

  usuario_atribuido_id BIGINT NULL,

  pausada             BOOLEAN NOT NULL DEFAULT false,

  outcome             VARCHAR(20) NULL,

  deleted_at          TIMESTAMPTZ NULL,
  deleted_by          BIGINT NULL,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT fk_tarefa_projeto
    FOREIGN KEY (projeto_id) REFERENCES projeto(id) ON DELETE CASCADE,

  CONSTRAINT fk_tarefa_usuario_atribuido
    FOREIGN KEY (usuario_atribuido_id) REFERENCES usuario(id),

  CONSTRAINT fk_tarefa_deleted_by
    FOREIGN KEY (deleted_by) REFERENCES usuario(id),

  CONSTRAINT ck_tarefa_status
    CHECK (status IN ('TODO', 'DOING', 'INREVIEW', 'DONE')),

  CONSTRAINT ck_tarefa_outcome
    CHECK (outcome IS NULL OR outcome IN ('APPROVED', 'REJECTED'))
);

-- regra: 1 tarefa DOING (não pausada) por usuário (ignorando soft delete)
CREATE UNIQUE INDEX IF NOT EXISTS uq_tarefa_um_doing_por_usuario
  ON tarefa(usuario_atribuido_id)
  WHERE status = 'DOING'
    AND pausada = false
    AND usuario_atribuido_id IS NOT NULL
    AND deleted_at IS NULL;

COMMIT;