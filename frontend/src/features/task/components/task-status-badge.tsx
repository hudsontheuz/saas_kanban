import { Badge } from '@/components/ui/badge';
import type { TaskStatus } from '@/features/task/types/task.types';

export function TaskStatusBadge({ status }: { status: TaskStatus }) {
  const map = {
    TODO: { label: 'ToDo', variant: 'secondary' as const },
    DOING: { label: 'Doing', variant: 'warning' as const },
    IN_REVIEW: { label: 'In Review', variant: 'default' as const },
    DONE: { label: 'Done', variant: 'success' as const },
  };
  const current = map[status];
  return <Badge variant={current.variant}>{current.label}</Badge>;
}
