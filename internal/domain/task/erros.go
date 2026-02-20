package task

import "errors"

var (
	ErrProjetoObrigatorio               = errors.New("project é obrigatório")
	ErrTituloObrigatorio                = errors.New("título é obrigatório")
	ErrTransicaoInvalida                = errors.New("transição inválida")
	ErrOutcomeSomenteEmDone             = errors.New("Resultado somente quando a tarefa for concluida")
	ErrSemAssignee                      = errors.New("task sem assignee")
	ErrSomenteDoingPodePausar           = errors.New("só pode pausar/retomar quando status=DOING")
	ErrJaPausada                        = errors.New("task já está pausada")
	ErrNaoEstaPausada                   = errors.New("task não está pausada")
	ErrRejeitarSomenteEmTodo            = errors.New("rejeitar só é permitido em ToDo")
	ErrReprovarSomenteInReview          = errors.New("reprovar para ajustes só é permitido em InReview")
	ErrSomenteAssigneePodePausar        = errors.New("somente o assignee pode pausar/retomar")
	ErrSomenteAssigneePodeMoverInReview = errors.New("somente o assignee pode mover para InReview")
)
