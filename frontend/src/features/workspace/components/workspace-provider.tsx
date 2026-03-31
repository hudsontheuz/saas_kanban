import * as React from 'react';
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useAuth } from '@/features/auth/components/auth-provider';
import type { Project, ProjectSettings } from '@/features/project/types/project.types';
import type { Task, TaskStatus } from '@/features/task/types/task.types';
import type { Team, TeamMember } from '@/features/team/types/team.types';
import { getErrorMessage } from '@/lib/api-error';
import { workspaceApi } from '@/features/workspace/api/workspace.api';
import { workspaceStorage } from '@/lib/workspace-storage';

interface WorkspaceContextValue {
  team: Team | null;
  project: Project | null;
  tasks: Task[];
  members: TeamMember[];
  stats: Record<TaskStatus, number>;
  isLoading: boolean;
  message: string | null;
  clearMessage: () => void;
  refresh: () => Promise<void>;
  createTeam: (name: string) => Promise<void>;
  createProject: (name: string) => Promise<void>;
  closeProject: () => Promise<void>;
  updateProjectSettings: (settings: ProjectSettings) => Promise<void>;
  createTask: (title: string, description: string) => Promise<void>;
  selfAssignTask: (taskId: string) => Promise<void>;
  pauseTask: (taskId: string) => Promise<void>;
  resumeTask: (taskId: string) => Promise<void>;
  moveTaskToReview: (taskId: string, deliveryComment: string) => Promise<void>;
  approveTask: (taskId: string) => Promise<void>;
  rejectTask: (taskId: string, reason: string) => Promise<void>;
}

const defaultStats: Record<TaskStatus, number> = {
  TODO: 0,
  DOING: 0,
  IN_REVIEW: 0,
  DONE: 0,
};

const defaultSettings: ProjectSettings = {
  allowDropTask: false,
  leaderApprovalRequired: true,
  leaderTaskNeedsApproval: true,
  peerApprovalAllowed: false,
};

const WorkspaceContext = createContext<WorkspaceContextValue | null>(null);

export function WorkspaceProvider({ children }: { children: React.ReactNode }) {
  const { user, isAuthenticated } = useAuth();
  const [team, setTeam] = useState<Team | null>(null);
  const [project, setProject] = useState<Project | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [members, setMembers] = useState<TeamMember[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [message, setMessage] = useState<string | null>(null);

  const resetWorkspace = useCallback(() => {
    setTeam(null);
    setProject(null);
    setTasks([]);
    setMembers([]);
  }, []);

  const refresh = useCallback(async () => {
    if (!isAuthenticated) {
      resetWorkspace();
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    try {
      const snapshot = await workspaceApi.bootstrap(user);
      setTeam(snapshot.team);
      setProject(snapshot.project);
      setTasks(snapshot.tasks);
      setMembers(snapshot.members);
    } catch (error) {
      setMessage(getErrorMessage(error, 'Não foi possível carregar o workspace.'));
    } finally {
      setIsLoading(false);
    }
  }, [isAuthenticated, resetWorkspace, user]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  useEffect(() => {
    if (!isAuthenticated) {
      workspaceStorage.clear();
      resetWorkspace();
    }
  }, [isAuthenticated, resetWorkspace]);

  const runAction = useCallback(async (action: () => Promise<void>, successMessage: string) => {
    try {
      await action();
      setMessage(successMessage);
    } catch (error) {
      setMessage(getErrorMessage(error));
      throw error;
    }
  }, []);

  const createTeam = useCallback(async (name: string) => {
    if (!user) throw new Error('Usuário não autenticado.');

    await runAction(async () => {
      const snapshot = await workspaceApi.createTeam(name, user);
      setTeam(snapshot.team);
      setProject(null);
      setTasks([]);
      setMembers(snapshot.members);
    }, 'Equipe criada com sucesso.');
  }, [runAction, user]);

  const createProject = useCallback(async (name: string) => {
    if (!team) throw new Error('Crie uma equipe antes de criar um projeto.');

    await runAction(async () => {
      const createdProject = await workspaceApi.createProject(team.id, {
        name,
        settings: defaultSettings,
      });
      setProject(createdProject);
      setTasks([]);
    }, 'Projeto criado com sucesso.');
  }, [runAction, team]);

  const closeProject = useCallback(async () => {
    if (!project) throw new Error('Nenhum projeto ativo para fechar.');

    await runAction(async () => {
      await workspaceApi.closeProject(project.id);
      setProject((current) =>
        current
          ? {
              ...current,
              active: false,
              closedAt: new Date().toISOString(),
            }
          : current,
      );
    }, 'Projeto fechado com sucesso.');
  }, [project, runAction]);

  const updateProjectSettings = useCallback(async (settings: ProjectSettings) => {
    if (!project) throw new Error('Nenhum projeto carregado.');

    await runAction(async () => {
      await workspaceApi.updateProjectSettings(project.id, settings);
      setProject((current) => (current ? { ...current, settings: { ...current.settings, ...settings } } : current));
    }, 'Configurações atualizadas com sucesso.');
  }, [project, runAction]);

  const createTask = useCallback(async (title: string, description: string) => {
    if (!project) throw new Error('Crie um projeto antes de criar tarefas.');

    await runAction(async () => {
      const task = await workspaceApi.createTask(project.id, { title, description });
      setTasks((current) => [task, ...current]);
    }, 'Tarefa criada com sucesso.');
  }, [project, runAction]);

  const selfAssignTask = useCallback(async (taskId: string) => {
    if (!user) throw new Error('Usuário não autenticado.');

    await runAction(async () => {
      await workspaceApi.selfAssignTask(taskId);
      setTasks((current) =>
        current.map((task) =>
          task.id === taskId
            ? { ...task, assigneeId: user.id, status: 'DOING', paused: false }
            : task,
        ),
      );
    }, 'Tarefa assumida com sucesso.');
  }, [runAction, user]);

  const pauseTask = useCallback(async (taskId: string) => {
    await runAction(async () => {
      await workspaceApi.pauseTask(taskId);
      setTasks((current) => current.map((task) => (task.id === taskId ? { ...task, paused: true } : task)));
    }, 'Tarefa pausada com sucesso.');
  }, [runAction]);

  const resumeTask = useCallback(async (taskId: string) => {
    await runAction(async () => {
      await workspaceApi.resumeTask(taskId);
      setTasks((current) => current.map((task) => (task.id === taskId ? { ...task, paused: false } : task)));
    }, 'Tarefa retomada com sucesso.');
  }, [runAction]);

  const moveTaskToReview = useCallback(async (taskId: string, deliveryComment: string) => {
    await runAction(async () => {
      await workspaceApi.moveTaskToReview(taskId, deliveryComment);
      setTasks((current) => current.map((task) => (task.id === taskId ? { ...task, status: 'IN_REVIEW', paused: false, deliveryComment } : task)));
    }, 'Tarefa enviada para revisão.');
  }, [runAction]);

  const approveTask = useCallback(async (taskId: string) => {
    await runAction(async () => {
      await workspaceApi.approveTask(taskId);
      setTasks((current) => current.map((task) => (task.id === taskId ? { ...task, status: 'DONE', paused: false, reviewComment: undefined } : task)));
    }, 'Tarefa aprovada com sucesso.');
  }, [runAction]);

  const rejectTask = useCallback(async (taskId: string, reason: string) => {
    await runAction(async () => {
      await workspaceApi.rejectTask(taskId, reason);
      setTasks((current) => current.map((task) => (task.id === taskId ? { ...task, status: 'TODO', paused: false, reviewComment: reason } : task)));
    }, 'Tarefa reprovada e devolvida para To Do.');
  }, [runAction]);

  const stats = useMemo(() => {
    return tasks.reduce<Record<TaskStatus, number>>((acc, task) => {
      acc[task.status] += 1;
      return acc;
    }, { ...defaultStats });
  }, [tasks]);

  const value = useMemo<WorkspaceContextValue>(() => ({
    team,
    project,
    tasks,
    members,
    stats,
    isLoading,
    message,
    clearMessage: () => setMessage(null),
    refresh,
    createTeam,
    createProject,
    closeProject,
    updateProjectSettings,
    createTask,
    selfAssignTask,
    pauseTask,
    resumeTask,
    moveTaskToReview,
    approveTask,
    rejectTask,
  }), [
    team,
    project,
    tasks,
    members,
    stats,
    isLoading,
    message,
    refresh,
    createTeam,
    createProject,
    closeProject,
    updateProjectSettings,
    createTask,
    selfAssignTask,
    pauseTask,
    resumeTask,
    moveTaskToReview,
    approveTask,
    rejectTask,
  ]);

  return <WorkspaceContext.Provider value={value}>{children}</WorkspaceContext.Provider>;
}

export function useWorkspace() {
  const context = useContext(WorkspaceContext);
  if (!context) throw new Error('useWorkspace must be used inside WorkspaceProvider');
  return context;
}
