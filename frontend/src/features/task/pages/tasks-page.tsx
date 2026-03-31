import { useMemo, useState } from 'react';
import { PageHeader } from '@/components/layout/page-header';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/features/auth/components/auth-provider';
import { TaskCard } from '@/features/task/components/task-card';
import { useWorkspace } from '@/hooks/use-workspace';
import type { TaskStatus } from '@/features/task/types/task.types';
import { CreateTaskDialog } from '@/features/task/components/create-task-dialog';

const columns: { status: TaskStatus; title: string }[] = [
  { status: 'TODO', title: 'To Do' },
  { status: 'DOING', title: 'Doing' },
  { status: 'IN_REVIEW', title: 'In Review' },
  { status: 'DONE', title: 'Done' },
];

export function TasksPage() {
  const { user } = useAuth();
  const {
    project,
    tasks,
    members,
    isLoading,
    message,
    clearMessage,
    createTask,
    selfAssignTask,
    pauseTask,
    resumeTask,
    moveTaskToReview,
    approveTask,
    rejectTask,
  } = useWorkspace();
  const [query, setQuery] = useState('');

  const filtered = useMemo(() => {
    const lower = query.toLowerCase();
    return tasks.filter((task) => task.title.toLowerCase().includes(lower) || task.description.toLowerCase().includes(lower));
  }, [query, tasks]);

  const isLeader = members.some((member) => member.id === user?.id && member.role === 'LEADER');

  return (
    <div className="space-y-6">
      <PageHeader
        title="Tarefas"
        description="Acompanhe o fluxo das tarefas do projeto ativo."
        action={
          <div className="flex w-full flex-col gap-3 md:w-auto md:flex-row">
            <Input className="w-full md:w-72" placeholder="Buscar tarefa" value={query} onChange={(event) => setQuery(event.target.value)} />
            <CreateTaskDialog onCreate={createTask} disabled={!project} />
          </div>
        }
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {!project && !isLoading && (
        <div className="rounded-2xl border border-dashed bg-white p-6 text-sm text-slate-600 shadow-soft">
          Crie um projeto para começar a organizar as tarefas da equipe.
        </div>
      )}

      {isLoading && <div className="rounded-2xl border bg-white p-6 text-sm text-slate-600 shadow-soft">Carregando quadro...</div>}

      {project && (
        <div className="grid gap-4 xl:grid-cols-4">
          {columns.map((column) => {
            const columnTasks = filtered.filter((task) => task.status === column.status);

            return (
              <Card key={column.status} className="bg-slate-100/60">
                <CardHeader>
                  <CardTitle className="text-base">{column.title}</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {columnTasks.map((task) => (
                    <TaskCard
                      key={task.id}
                      task={task}
                      members={members}
                      currentUserId={user?.id ?? ''}
                      canReview={isLeader}
                      onSelfAssign={(taskId) => void selfAssignTask(taskId)}
                      onPause={(taskId) => void pauseTask(taskId)}
                      onResume={(taskId) => void resumeTask(taskId)}
                      onMoveToReview={(taskId, comment) => void moveTaskToReview(taskId, comment)}
                      onApprove={(taskId) => void approveTask(taskId)}
                      onReject={(taskId, reason) => void rejectTask(taskId, reason)}
                    />
                  ))}
                  {columnTasks.length === 0 && (
                    <div className="rounded-xl border border-dashed border-slate-300 bg-white p-4 text-center text-sm text-slate-500">
                      Nenhuma tarefa nesta etapa.
                    </div>
                  )}
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
