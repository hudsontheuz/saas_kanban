BEGIN;

ALTER TABLE tarefa
  DROP COLUMN IF EXISTS comentario_entrega;

ALTER TABLE tarefa
  DROP COLUMN IF EXISTS comentario_review;

COMMIT;