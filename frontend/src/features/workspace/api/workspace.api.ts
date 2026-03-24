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
  description: string;
  settings: ProjectSettings;
}

interface CreateTaskInput {
  title: string;
  description: string;
}

interface AddMemberInput {
  name: string;
  email: string;
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
    id: pickString(data, ['id', 'memberId', 'member_id', 'userId', 'user_id', 'ID'], crypto.randomUUID()),
    name: pickString(data, ['name', 'nome', 'Name'], 'Membro'),
    email: pickString(data, ['email', 'Email']),
    role: roleValue === 'LEADER' ? 'LEADER' : 'MEMBER',
  };
}

function normalizeSettings(input: unknown): ProjectSettings {
  const data = asObject(input);
  return {
    allowDropTask: pickBoolean(data, ['allowDropTask', 'permitirSoltarDoingParaTodo', 'PermitirSoltarDoingParaTodo'], false),
    leaderApprovalRequired: pickBoolean(data, ['leaderApprovalRequired', 'LeaderApprovalRequired'], true),
    leaderTaskNeedsApproval: pickBoolean(data, ['leaderTaskNeedsApproval', 'LeaderTaskNeedsApproval'], true),
    peerApprovalAllowed: pickBoolean(data, ['peerApprovalAllowed', 'PeerApprovalAllowed'], false),
  };
}

function normalizeProject(input: unknown): Project {
  const data = asObject(input);
  return {
    id: pickString(data, ['id', 'projectId', 'project_id', 'ProjectID', 'ID'], crypto.randomUUID()),
    name: pickString(data, ['name', 'nome', 'Name'], 'Projeto sem nome'),
    description: pickString(data, ['description', 'descricao', 'Description']),
    active: !pickString(data, ['closedAt', 'closed_at', 'fechadoEm', 'FechadoEm']),
    createdAt: pickString(data, ['createdAt', 'created_at', 'criadoEm', 'CriadoEm'], new Date().toISOString()),
    closedAt: pickString(data, ['closedAt', 'closed_at', 'fechadoEm', 'FechadoEm']) || undefined,
    settings: normalizeSettings(pickValue(data, 'settings', 'configuracoes', 'configuracoesProject')),
  };
}

function normalizeTask(input: unknown): Task {
  const data = asObject(input);
  return {
    id: pickString(data, ['id', 'taskId', 'task_id', 'TaskID', 'ID'], crypto.randomUUID()),
    title: pickString(data, ['title', 'titulo', 'Title'], 'Tarefa sem título'),
    description: pickString(data, ['description', 'descricao', 'Description']),
    status: normalizeTaskStatus(pickValue(data, 'status', 'estado', 'Status')),
    assigneeId: pickString(data, ['assigneeId', 'assignee_id', 'responsavelId', 'executorId']) || undefined,
    selectedUserId: pickString(data, ['selectedUserId', 'selected_user_id', 'sugeridoParaId']) || undefined,
    paused: pickBoolean(data, ['paused', 'pausada', 'Paused'], false),
  };
}

function normalizeTeam(input: unknown): Team {
  const data = asObject(input);
  const membersRaw = pickValue(data, 'members', 'membros', 'Members');
  const members = Array.isArray(membersRaw) ? membersRaw.map(normalizeMember) : [];

  return {
    id: pickString(data, ['id', 'teamId', 'team_id', 'TeamID', 'ID'], crypto.randomUUID()),
    name: pickString(data, ['name', 'nome', 'Name'], 'Equipe'),
    members,
  };
}

function normalizeProjects(input: unknown): Project[] {
  if (Array.isArray(input)) return input.map(normalizeProject);
  const data = asObject(input);
  const items = pickValue(data, 'projects', 'projetos', 'items', 'data');
  return Array.isArray(items) ? items.map(normalizeProject) : [];
}

function normalizeTasks(input: unknown): Task[] {
  if (Array.isArray(input)) return input.map(normalizeTask);
  const data = asObject(input);
  const items = pickValue(data, 'tasks', 'tarefas', 'items', 'data');
  return Array.isArray(items) ? items.map(normalizeTask) : [];
}

function optimisticTeam(name: string, currentUser: AuthUser): Team {
  return {
    id: crypto.randomUUID(),
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
}

function optimisticProject(id: string, input: CreateProjectInput): Project {
  return {
    id,
    name: input.name,
    description: input.description,
    active: true,
    createdAt: new Date().toISOString(),
    settings: input.settings,
  };
}

function optimisticTask(id: string, input: CreateTaskInput): Task {
  return {
    id,
    title: input.title,
    description: input.description,
    status: 'TODO',
    paused: false,
  };
}

export const workspaceApi = {
  async bootstrap(): Promise<WorkspaceSnapshot> {
    const teamId = workspaceStorage.getTeamId();
    const selectedProjectId = workspaceStorage.getProjectId();

    if (!teamId) {
      return { team: null, project: null, tasks: [], members: [] };
    }

    let team: Team | null = null;
    let project: Project | null = null;
    let tasks: Task[] = [];

    try {
      const teamResponse = await apiClient.get(`/teams/${teamId}`);
      team = normalizeTeam(teamResponse.data);
    } catch {
      team = null;
    }

    try {
      const projectsResponse = await apiClient.get(`/teams/${teamId}/projects`);
      const projects = normalizeProjects(projectsResponse.data);
      project = projects.find((item) => item.id === selectedProjectId) ?? projects.find((item) => item.active) ?? projects[0] ?? null;
      if (project) {
        workspaceStorage.setProjectId(project.id);
      }
    } catch {
      project = selectedProjectId
        ? {
            id: selectedProjectId,
            name: 'Projeto atual',
            description: '',
            active: true,
            createdAt: new Date().toISOString(),
            settings: normalizeSettings({}),
          }
        : null;
    }

    if (project?.id) {
      try {
        const tasksResponse = await apiClient.get(`/projects/${project.id}/tasks`);
        tasks = normalizeTasks(tasksResponse.data);
      } catch {
        tasks = [];
      }
    }

    const members = team?.members ?? [];
    return { team, project, tasks, members };
  },

  async createTeam(name: string, currentUser: AuthUser): Promise<WorkspaceSnapshot> {
    const { data } = await apiClient.post('/teams', { nome: name });
    const teamId = extractId(data, 'teamId', 'team_id', 'TeamID', 'id', 'ID');
    workspaceStorage.setTeamId(teamId);
    workspaceStorage.clearProjectId();

    const createdTeam = optimisticTeam(name, currentUser);
    createdTeam.id = teamId;

    return {
      team: createdTeam,
      project: null,
      tasks: [],
      members: createdTeam.members,
    };
  },

  async createProject(teamId: string, input: CreateProjectInput): Promise<Project> {
    const { data } = await apiClient.post(`/teams/${teamId}/projects`, {
      nome: input.name,
      descricao: input.description,
      permitirSoltarDoingParaTodo: input.settings.allowDropTask,
      leaderApprovalRequired: input.settings.leaderApprovalRequired,
      leaderTaskNeedsApproval: input.settings.leaderTaskNeedsApproval,
      peerApprovalAllowed: input.settings.peerApprovalAllowed,
    });

    const projectId = extractId(data, 'projectId', 'project_id', 'ProjectID', 'id', 'ID');
    workspaceStorage.setProjectId(projectId);
    return optimisticProject(projectId, input);
  },

  async closeProject(projectId: string): Promise<void> {
    await apiClient.post(`/projects/${projectId}/close`);
  },

  async updateProjectSettings(projectId: string, settings: ProjectSettings): Promise<void> {
    await apiClient.put(`/projects/${projectId}/settings`, {
      allowDropTask: settings.allowDropTask,
      leaderApprovalRequired: settings.leaderApprovalRequired,
      leaderTaskNeedsApproval: settings.leaderTaskNeedsApproval,
      peerApprovalAllowed: settings.peerApprovalAllowed,
      permitirSoltarDoingParaTodo: settings.allowDropTask,
    });
  },

  async createTask(projectId: string, input: CreateTaskInput): Promise<Task> {
    const { data } = await apiClient.post(`/projects/${projectId}/tasks`, {
      titulo: input.title,
      descricao: input.description,
    });

    const taskId = extractId(data, 'taskId', 'task_id', 'TaskID', 'id', 'ID');
    return optimisticTask(taskId, input);
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

  async moveTaskToReview(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/in-review`);
  },

  async approveTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/approve`);
  },

  async rejectTask(taskId: string): Promise<void> {
    await apiClient.post(`/tasks/${taskId}/reject`);
  },

  async addMember(teamId: string, input: AddMemberInput): Promise<TeamMember> {
    const { data } = await apiClient.post(`/teams/${teamId}/members`, {
      nome: input.name,
      email: input.email,
    });

    return normalizeMember({
      ...asObject(data),
      name: input.name,
      email: input.email,
    });
  },

  async removeMember(teamId: string, memberId: string): Promise<void> {
    await apiClient.delete(`/teams/${teamId}/members/${memberId}`);
  },

  async transferLeadership(teamId: string, memberId: string): Promise<void> {
    await apiClient.post(`/teams/${teamId}/leader/transfer`, {
      memberId,
      member_id: memberId,
      novoLeaderId: memberId,
    });
  },
};
