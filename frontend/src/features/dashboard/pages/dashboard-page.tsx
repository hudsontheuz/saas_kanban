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
    { label: 'Membros no time', value: team ? `${members.length}/5` : '0/5', icon: Users },
    { label: 'Tarefas em Doing', value: String(stats.DOING), icon: ListTodo },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description="Visão geral do workspace e do projeto ativo da equipe."
        action={<Button asChild><Link to="/tasks">Ir para o quadro</Link></Button>}
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {isLoading && <div className="rounded-2xl border bg-white p-6 text-sm text-slate-600 shadow-soft">Carregando workspace...</div>}

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
            <CardTitle>{project?.name ?? 'Nenhum projeto carregado'}</CardTitle>
            <CardDescription>{project?.description ?? 'Crie uma equipe e depois um projeto para começar o fluxo real.'}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm text-slate-600">
            <div className="flex flex-wrap gap-2">
              <Badge variant={project?.settings.allowDropTask ? 'success' : 'secondary'}>
                Soltar tarefa: {project?.settings.allowDropTask ? 'Sim' : 'Não'}
              </Badge>
              <Badge variant={project?.settings.leaderApprovalRequired ? 'warning' : 'secondary'}>
                Aprovação do líder: {project?.settings.leaderApprovalRequired ? 'Obrigatória' : 'Livre'}
              </Badge>
            </div>
            <p>
              O frontend agora está preparado para refletir chamadas reais de backend para equipe, projeto e tarefas,
              deixando o mock de lado.
            </p>
            <Button asChild variant="outline">
              <Link to="/projects" className="gap-2">Ver projeto <ArrowRight className="h-4 w-4" /></Link>
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Resumo do fluxo</CardTitle>
            <CardDescription>Contagem atual de tarefas por status.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>ToDo</span><strong>{stats.TODO}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>Doing</span><strong>{stats.DOING}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>In Review</span><strong>{stats.IN_REVIEW}</strong></div>
            <div className="flex items-center justify-between rounded-xl bg-slate-50 p-3"><span>Done</span><strong>{stats.DONE}</strong></div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
