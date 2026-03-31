import { apiClient } from '@/lib/api-client';
import type { AuthUser } from '@/features/auth/types/auth.types';
import type { Project, ProjectSettings } from '@/features/project/types/project.types';
import type { Task, TaskStatus } from '@/features/task/types/task.types';
import type { Team, TeamMember } from '@/features/team/types/team.types';
import { workspaceStorage } from '@/lib/workspace-storage';

export interface WorkspaceSnapshot {
  team: Team | null;
  project: Project | null;
  tasks: Task[];
  members: TeamMember[];
}

interface CreateProjectInput {
  name: string;
  settings: ProjectSettings;
}

interface CreateTaskInput {
  title: string;
  description: string;
}

function asObject(value: unknown): Record<string, unknown> {
  return value && typeof value === 'object' ? (value as Record<string, unknown>) : {};
}

function pickValue(source: Record<string, unknown>, ...keys: string[]): unknown {
  for (const key of keys) {
    const value = source[key];
    if (value !== undefined && value !== null) return value;
  }
  return undefined;
}

function pickString(source: Record<string, unknown>, keys: string[], fallback = ''): string {
  const value = pickValue(source, ...keys);
  if (typeof value === 'string') return value;
  if (typeof value === 'number') return String(value);
  return fallback;
}

function pickBoolean(source: Record<string, unknown>, keys: string[], fallback = false): boolean {
  const value = pickValue(source, ...keys);
  if (typeof value === 'boolean') return value;
  return fallback;
}

function extractId(source: unknown, ...keys: string[]): string {
  const data = asObject(source);
  const id = pickString(data, keys);
  if (!id) throw new Error('A API não retornou um identificador válido.');
  return id;
}

function normalizeTaskStatus(status: unknown): TaskStatus {
  const value = typeof status === 'string' ? status.toUpperCase() : 'TODO';
  switch (value) {
    case 'DOING':
      return 'DOING';
    case 'IN_REVIEW':
    case 'INREVIEW':
      return 'IN_REVIEW';
    case 'DONE':
      return 'DONE';
    default:
      return 'TODO';
  }
}

function normalizeMember(input: unknown): TeamMember {
  const data = asObject(input);
  const roleValue = pickString(data, ['role', 'papel', 'Role'], 'MEMBER').toUpperCase();
  return {
    id: pickString(data, ['id', 'userId', 'user_id', 'memberId', 'member_id', 'ID'], crypto.randomUUID()),
    name: pickString(data, ['name', 'nome', 'Name'], 'Membro'),
    email: pickString(data, ['email', 'Email']),
    role: roleValue === 'LEADER' ? 'LEADER' : 'MEMBER',
  };
}

function normalizeSettings(input: unknown): ProjectSettings {
  const data = asObject(input);
  const allowDropTask = pickBoolean(
    data,
    ['allowDropTask', 'permitirSoltarDoingParaTodo', 'permitir_soltar_doing_para_todo', 'PermitirSoltarDoingParaTodo'],
    false,
  );

  return {
    allowDropTask,
    leaderApprovalRequired: true,
    leaderTaskNeedsApproval: true,
    peerApprovalAllowed: false,
  };
}

function normalizeProject(input: unknown): Project {
  const data = asObject(input);
  const status = pickString(data, ['status', 'Status']).toUpperCase();
  const closedAt = pickString(data, ['closedAt', 'closed_at', 'fechadoEm', 'FechadoEm']) || undefined;
  const settingsSource = pickValue(data, 'settings', 'configuracoes', 'configuracoesProject') ?? data;

  return {
    id: pickString(data, ['id', 'projectId', 'project_id', 'ProjectID', 'ID'], crypto.randomUUID()),
    name: pickString(data, ['name', 'nome', 'Name'], 'Projeto'),
    description: '',
    active: status ? status === 'ACTIVE' : !closedAt,
    createdAt: pickString(data, ['createdAt', 'created_at', 'criadoEm', 'CriadoEm'], new Date().toISOString()),
    closedAt,
    settings: normalizeSettings(settingsSource),
  };
}

function normalizeTask(input: unknown): Task {
  const data = asObject(input);
  return {
    id: pickString(data, ['id', 'taskId', 'task_id', 'TaskID', 'ID'], crypto.randomUUID()),
    title: pickString(data, ['title', 'titulo', 'Title'], 'Tarefa'),
    description: pickString(data, ['description', 'descricao', 'Description']),
    deliveryComment: pickString(data, ['deliveryComment', 'comentario_entrega', 'comentarioEntrega']) || undefined,
    reviewComment: pickString(data, ['reviewComment', 'comentario_review', 'comentarioReview']) || undefined,
    status: normalizeTaskStatus(pickValue(data, 'status', 'estado', 'Status')),
    assigneeId: pickString(data, ['assigneeId', 'assignee_id', 'responsavelId', 'executorId']) || undefined,
    selectedUserId: pickString(data, ['selectedUserId', 'selected_user_id', 'sugeridoParaId']) || undefined,
    paused: pickBoolean(data, ['paused', 'pausada', 'Paused'], false),
  };
}

function normalizeTasks(input: unknown): Task[] {
  if (Array.isArray(input)) return input.map(normalizeTask);
  const data = asObject(input);
  const items = pickValue(data, 'tasks', 'tarefas', 'items', 'data');
  return Array.isArray(items) ? items.map(normalizeTask) : [];
}

function normalizeTeam(input: unknown, currentUser: AuthUser | null): Team {
  const data = asObject(input);
  const membersRaw = pickValue(data, 'members', 'membros', 'items');
  const members = Array.isArray(membersRaw) ? membersRaw.map(normalizeMember) : [];
  const fallbackMember = currentUser
    ? [{ id: currentUser.id, name: currentUser.name, email: currentUser.email, role: 'LEADER' as const }]
    : [];

  return {
    id: pickString(data, ['id', 'teamId', 'team_id', 'TeamID', 'ID']),
    name: pickString(data, ['name', 'nome', 'Name'], 'Minha equipe'),
    members: members.length > 0 ? members : fallbackMember,
  };
}

async function resolveTeamId(currentUser: AuthUser | null): Promise<string | null> {
  const storedTeamId = workspaceStorage.getTeamId();
  if (storedTeamId) return storedTeamId;

  try {
    const { data } = await apiClient.get('/me/teams');
    const items = Array.isArray((data as { items?: unknown[] })?.items) ? (data as { items: unknown[] }).items : [];
    const firstTeam = items[0];
    if (!firstTeam) return null;
    const teamId = extractId(firstTeam, 'teamId', 'team_id', 'TeamID', 'id', 'ID');
    workspaceStorage.setTeamId(teamId);
    if (currentUser) {
      workspaceStorage.setTeamSnapshot({
        id: teamId,
        name: pickString(asObject(firstTeam), ['name', 'nome', 'Name'], 'Minha equipe'),
        members: [{ id: currentUser.id, name: currentUser.name, email: currentUser.email, role: 'LEADER' }],
      });
    }
    return teamId;
  } catch {
    return null;
  }
}

export const workspaceApi = {
  async bootstrap(currentUser: AuthUser | null): Promise<WorkspaceSnapshot> {
    const teamId = await resolveTeamId(currentUser);
    if (!teamId) {
      return { team: null, project: null, tasks: [], members: [] };
    }

    let team: Team | null = null;
    let members: TeamMember[] = [];

    try {
      const teamResponse = await apiClient.get(`/teams/${teamId}`);
      team = normalizeTeam(teamResponse.data, currentUser);
      members = team.members;
      workspaceStorage.setTeamSnapshot(team);
      workspaceStorage.setTeamId(team.id);
    } catch {
      const snapshot = workspaceStorage.getTeamSnapshot();
      team = snapshot ? { id: snapshot.id, name: snapshot.name, members: snapshot.members } : null;
      members = team?.members ?? [];
    }

    let project: Project | null = null;
    let tasks: Task[] = [];

    try {
      const projectResponse = await apiClient.get(`/teams/${teamId}/projects/active`);
      project = normalizeProject(projectResponse.data);
      workspaceStorage.setProjectId(project.id);
    } catch {
      workspaceStorage.clearProjectId();
      project = null;
    }

    if (project?.id) {
      try {
        const tasksResponse = await apiClient.get(`/projects/${project.id}/tasks`);
        tasks = normalizeTasks(tasksResponse.data);
      } catch {
        tasks = [];
      }
    }

    return { team, project, tasks, members };
  },

  async createTeam(name: string, currentUser: AuthUser): Promise<WorkspaceSnapshot> {
    const { data } = await apiClient.post('/teams', { nome: name });
    const teamId = extractId(data, 'teamId', 'team_id', 'TeamID', 'id', 'ID');
    const createdTeam: Team = {
      id: teamId,
      name,
      members: [
        {
          id: currentUser.id,
          name: currentUser.name,
          email: currentUser.email,
          role: 'LEADER',
        },
      ],
    };

    workspaceStorage.setTeamId(teamId);
    workspaceStorage.setTeamSnapshot(createdTeam);
    workspaceStorage.clearProjectId();

    return {
      team: createdTeam,
      project: null,
      tasks: [],
      members: createdTeam.members,
    };
  },

  async createProject(teamId: string, input: CreateProjectInput): Promise<Project> {
    await apiClient.post(`/teams/${teamId}/projects`, {
      nome: input.name,
      permitir_soltar_doing_para_todo: input.settings.allowDropTask,
    });

    const activeProjectResponse = await apiClient.get(`/teams/${teamId}/projects/active`);
    const project = normalizeProject(activeProjectResponse.data);
    workspaceStorage.setProjectId(project.id);
    return project;
  },

  async closeProject(projectId: string): Promise<void> {
    await apiClient.post(`/projects/${projectId}/close`);
    workspaceStorage.clearProjectId();
  },

  async updateProjectSettings(projectId: string, settings: ProjectSettings): Promise<void> {
    await apiClient.patch(`/projects/${projectId}/settings`, {
      permitir_soltar_doing_para_todo: settings.allowDropTask,
    });
  },

  async createTask(projectId: string, input: CreateTaskInput): Promise<Task> {
    const { data } = await apiClient.post(`/projects/${projectId}/tasks`, {
      titulo: input.title,
      descricao: input.description,
    });

    const taskId = extractId(data, 'taskId', 'task_id', 'TaskID', 'id', 'ID');
    return {
      id: taskId,
      title: input.title,
      description: input.description,
      status: 'TODO',
      paused: false,
    };
  },

  async selfAssignTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/self-assign`);
  },

  async pauseTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/pause`);
  },

  async resumeTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/resume`);
  },

  async moveTaskToReview(taskId: string, deliveryComment: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/in-review`, { comentario_entrega: deliveryComment });
  },

  async approveTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/approve`);
  },

  async rejectTask(taskId: string, reason: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/reject`, { motivo: reason });
  },
};
