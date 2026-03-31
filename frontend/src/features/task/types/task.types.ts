export type TaskStatus = 'TODO' | 'DOING' | 'IN_REVIEW' | 'DONE';

export interface Task {
  id: string;
  title: string;
  description: string;
  deliveryComment?: string;
  reviewComment?: string;
  status: TaskStatus;
  assigneeId?: string;
  selectedUserId?: string;
  paused?: boolean;
}
