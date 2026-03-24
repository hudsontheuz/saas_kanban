const TEAM_KEY = 'saas-kanban.current-team-id';
const PROJECT_KEY = 'saas-kanban.current-project-id';

export const workspaceStorage = {
  getTeamId: () => localStorage.getItem(TEAM_KEY),
  setTeamId: (teamId: string) => localStorage.setItem(TEAM_KEY, teamId),
  clearTeamId: () => localStorage.removeItem(TEAM_KEY),
  getProjectId: () => localStorage.getItem(PROJECT_KEY),
  setProjectId: (projectId: string) => localStorage.setItem(PROJECT_KEY, projectId),
  clearProjectId: () => localStorage.removeItem(PROJECT_KEY),
  clear: () => {
    localStorage.removeItem(TEAM_KEY);
    localStorage.removeItem(PROJECT_KEY);
  },
};
