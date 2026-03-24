import { Check, Pause, Play, Send, UserPlus, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { TaskStatusBadge } from '@/features/task/components/task-status-badge';
import type { Task } from '@/features/task/types/task.types';
import type { TeamMember } from '@/features/team/types/team.types';

export function TaskCard({
  task,
  members,
  currentUserId,
  canReview,
  onSelfAssign,
  onPause,
  onResume,
  onMoveToReview,
  onApprove,
  onReject,
}: {
  task: Task;
  members: TeamMember[];
  currentUserId: string;
  canReview: boolean;
  onSelfAssign: (taskId: string) => void;
  onPause: (taskId: string) => void;
  onResume: (taskId: string) => void;
  onMoveToReview: (taskId: string) => void;
  onApprove: (taskId: string) => void;
  onReject: (taskId: string) => void;
}) {
  const assignee = members.find((member) => member.id === task.assigneeId);
  const selected = members.find((member) => member.id === task.selectedUserId);

  return (
    <Card className="rounded-xl">
      <CardContent className="space-y-4 p-4">
        <div className="space-y-2">
          <TaskStatusBadge status={task.status} />
          <h3 className="font-medium leading-tight">{task.title}</h3>
          <p className="text-sm text-slate-500">{task.description}</p>
        </div>

        <div className="space-y-1 text-xs text-slate-500">
          <p>Responsável: {assignee?.name ?? 'Ninguém'}</p>
          {selected && <p>Sugerido para: {selected.name}</p>}
          {task.paused && <p className="font-medium text-amber-700">Tarefa pausada</p>}
        </div>

        <div className="flex flex-wrap gap-2">
          {task.status === 'TODO' && (
            <Button size="sm" onClick={() => onSelfAssign(task.id)}>
              <UserPlus className="mr-2 h-4 w-4" /> Pegar
            </Button>
          )}
          {task.status === 'DOING' && task.assigneeId === currentUserId && !task.paused && (
            <>
              <Button size="sm" variant="outline" onClick={() => onPause(task.id)}>
                <Pause className="mr-2 h-4 w-4" /> Pausar
              </Button>
              <Button size="sm" onClick={() => onMoveToReview(task.id)}>
                <Send className="mr-2 h-4 w-4" /> Revisão
              </Button>
            </>
          )}
          {task.status === 'DOING' && task.assigneeId === currentUserId && task.paused && (
            <Button size="sm" variant="outline" onClick={() => onResume(task.id)}>
              <Play className="mr-2 h-4 w-4" /> Retomar
            </Button>
          )}
          {task.status === 'IN_REVIEW' && canReview && (
            <>
              <Button size="sm" onClick={() => onApprove(task.id)}>
                <Check className="mr-2 h-4 w-4" /> Aprovar
              </Button>
              <Button size="sm" variant="destructive" onClick={() => onReject(task.id)}>
                <X className="mr-2 h-4 w-4" /> Reprovar
              </Button>
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
