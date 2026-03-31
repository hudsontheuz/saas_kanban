package task

import "errors"

var (
	ErrProjetoObrigatorio               = errors.New("projeto é obrigatório")
	ErrTituloObrigatorio                = errors.New("título da task é obrigatório")
	ErrDescricaoObrigatoria             = errors.New("descrição da task é obrigatória")
	ErrMotivoReviewObrigatorio          = errors.New("motivo do review é obrigatório")
	ErrOutcomeSomenteEmDone             = errors.New("outcome só pode existir quando status = DONE")
	ErrTransicaoInvalida                = errors.New("transição de status inválida")
	ErrSemAssignee                      = errors.New("task sem assignee")
	ErrSomenteDoingPodePausar           = errors.New("somente task em Doing pode pausar/retomar")
	ErrJaPausada                        = errors.New("task já está pausada")
	ErrNaoEstaPausada                   = errors.New("task não está pausada")
	ErrSomenteAssigneePodePausar        = errors.New("somente o assignee pode pausar/retomar")
	ErrSomenteAssigneePodeMoverInReview = errors.New("somente o assignee pode mover para InReview")
	ErrMoverParaInReviewSomenteDoing    = errors.New("somente task em Doing pode mover para InReview")
	ErrTaskPausadaNaoPodeIrReview       = errors.New("task pausada não pode ir para InReview")
	ErrRejeitarSomenteEmTodo            = errors.New("rejeitar em Done só é permitido quando a task está em ToDo")
	ErrReprovarSomenteInReview          = errors.New("reprovar para ajustes só é permitido em InReview")
	ErrComentarioEntregaMuitoCurto      = errors.New("comentário de entrega muito curto")
	ErrComentarioEntregaObrigatorio     = errors.New("comentário de entrega é obrigatório")
)
