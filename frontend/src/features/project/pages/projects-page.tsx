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
        description="Crie e acompanhe o projeto ativo da equipe."
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
            <CardDescription>O projeto ativo fica vinculado à equipe atual.</CardDescription>
          </CardHeader>
        </Card>
      )}

      {project ? (
        <Card>
          <CardHeader>
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
              <div>
                <CardTitle>{project.name}</CardTitle>
                <CardDescription>Projeto ativo da equipe</CardDescription>
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
                Retorno para To Do: {project.settings.allowDropTask ? 'Permitido' : 'Bloqueado'}
              </Badge>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="rounded-xl bg-slate-50 p-4">
                <p className="text-xs uppercase tracking-wide text-slate-500">Criado em</p>
                <p className="mt-1 font-medium text-slate-900">{new Date(project.createdAt).toLocaleDateString('pt-BR')}</p>
              </div>
              <div className="rounded-xl bg-slate-50 p-4">
                <p className="text-xs uppercase tracking-wide text-slate-500">Status</p>
                <p className="mt-1 font-medium text-slate-900">{project.active ? 'Em andamento' : 'Encerrado'}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      ) : (
        team && !isLoading && (
          <Card>
            <CardHeader>
              <CardTitle>Nenhum projeto ativo</CardTitle>
              <CardDescription>Crie um projeto para iniciar o quadro de tarefas da equipe.</CardDescription>
            </CardHeader>
          </Card>
        )
      )}
    </div>
  );
}
