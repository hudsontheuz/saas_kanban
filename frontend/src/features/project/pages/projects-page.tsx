import { PageHeader } from '@/components/layout/page-header';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useWorkspace } from '@/hooks/use-workspace';
import { CreateProjectDialog } from '@/features/project/components/create-project-dialog';

export function ProjectsPage() {
  const { team, project, isLoading, message, clearMessage, createProject, closeProject } = useWorkspace();

  return (
    <div className="space-y-6">
      <PageHeader
        title="Projetos"
        description="A tela agora está preparada para criar e fechar projetos via API."
        action={<CreateProjectDialog onCreate={createProject} disabled={!team || Boolean(project?.active)} />}
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {!team && !isLoading && (
        <Card>
          <CardHeader>
            <CardTitle>Crie sua equipe primeiro</CardTitle>
            <CardDescription>O backend exige que o projeto pertença a uma equipe e respeita a regra de um projeto ativo por vez.</CardDescription>
          </CardHeader>
        </Card>
      )}

      {project ? (
        <Card>
          <CardHeader>
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
              <div>
                <CardTitle>{project.name}</CardTitle>
                <CardDescription>{project.description}</CardDescription>
              </div>
              {project.active && (
                <Button variant="outline" onClick={() => void closeProject()}>
                  Fechar projeto
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent className="space-y-4 text-sm text-slate-600">
            <div className="flex flex-wrap gap-2">
              <Badge variant={project.active ? 'success' : 'secondary'}>{project.active ? 'Projeto ativo' : 'Projeto fechado'}</Badge>
              <Badge variant={project.settings.allowDropTask ? 'success' : 'secondary'}>
                Soltar tarefa: {project.settings.allowDropTask ? 'Permitido' : 'Bloqueado'}
              </Badge>
              <Badge variant={project.settings.leaderApprovalRequired ? 'warning' : 'secondary'}>
                Aprovação do líder: {project.settings.leaderApprovalRequired ? 'Obrigatória' : 'Livre'}
              </Badge>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="rounded-xl bg-slate-50 p-4">
                <p className="text-xs uppercase tracking-wide text-slate-500">Criado em</p>
                <p className="mt-1 font-medium text-slate-900">{new Date(project.createdAt).toLocaleDateString('pt-BR')}</p>
              </div>
              <div className="rounded-xl bg-slate-50 p-4">
                <p className="text-xs uppercase tracking-wide text-slate-500">Fechado em</p>
                <p className="mt-1 font-medium text-slate-900">{project.closedAt ? new Date(project.closedAt).toLocaleDateString('pt-BR') : 'Ainda ativo'}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      ) : (
        team && !isLoading && (
          <Card>
            <CardHeader>
              <CardTitle>Nenhum projeto ativo</CardTitle>
              <CardDescription>A equipe já existe. Agora você pode criar o primeiro projeto real ligado ao backend.</CardDescription>
            </CardHeader>
          </Card>
        )
      )}
    </div>
  );
}
