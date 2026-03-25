import { ArrowRight, FolderKanban, ListTodo, Users } from 'lucide-react';
import { Link } from 'react-router-dom';
import { PageHeader } from '@/components/layout/page-header';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useWorkspace } from '@/hooks/use-workspace';

export function DashboardPage() {
  const { team, project, members, stats, isLoading, message, clearMessage } = useWorkspace();

  const cards = [
    { label: 'Projeto ativo', value: project?.name ?? 'Nenhum', icon: FolderKanban },
    { label: 'Pessoas na equipe', value: team ? String(members.length || 1) : '0', icon: Users },
    { label: 'Tarefas em andamento', value: String(stats.DOING), icon: ListTodo },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description="Acompanhe a equipe, o projeto ativo e a distribuição das tarefas."
        action={<Button asChild><Link to="/tasks">Abrir quadro</Link></Button>}
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {isLoading && <div className="rounded-2xl border bg-white p-6 text-sm text-slate-600 shadow-soft">Carregando dados do workspace...</div>}

      <div className="grid gap-4 md:grid-cols-3">
        {cards.map((card) => {
          const Icon = card.icon;
          return (
            <Card key={card.label}>
              <CardHeader className="flex-row items-center justify-between space-y-0">
                <div>
                  <CardDescription>{card.label}</CardDescription>
                  <CardTitle className="mt-2 text-lg">{card.value}</CardTitle>
                </div>
                <div className="rounded-xl bg-slate-100 p-3">
                  <Icon className="h-5 w-5 text-slate-700" />
                </div>
              </CardHeader>
            </Card>
          );
        })}
      </div>

      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <Card>
          <CardHeader>
            <CardTitle>{project?.name ?? 'Nenhum projeto ativo'}</CardTitle>
            <CardDescription>
              {team
                ? 'Use o quadro para acompanhar tarefas, revisar entregas e concluir o fluxo da equipe.'
                : 'Crie sua equipe para começar a organizar o trabalho.'}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm text-slate-600">
            <div className="flex flex-wrap gap-2">
              <Badge variant={project?.settings.allowDropTask ? 'success' : 'secondary'}>
                Retorno para To Do: {project?.settings.allowDropTask ? 'Permitido' : 'Bloqueado'}
              </Badge>
              <Badge variant={project?.active ? 'success' : 'secondary'}>
                {project?.active ? 'Projeto em andamento' : 'Sem projeto ativo'}
              </Badge>
            </div>
            <p>
              {project
                ? 'A equipe pode criar tarefas, assumir a execução e enviar itens para revisão dentro do mesmo fluxo.'
                : 'Depois de criar a equipe, você pode abrir um projeto e começar a registrar as tarefas.'}
            </p>
            <Button asChild variant="outline">
              <Link to="/projects" className="gap-2">Ver projeto <ArrowRight className="h-4 w-4" /></Link>
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Resumo do fluxo</CardTitle>
            <CardDescription>Contagem atual de tarefas por etapa.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>To Do</span><strong>{stats.TODO}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>Doing</span><strong>{stats.DOING}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>In Review</span><strong>{stats.IN_REVIEW}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>Done</span><strong>{stats.DONE}</strong></div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
