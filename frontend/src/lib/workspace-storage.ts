const TEAM_KEY = 'saas-kanban.current-team-id';
const PROJECT_KEY = 'saas-kanban.current-project-id';
const TEAM_SNAPSHOT_KEY = 'saas-kanban.current-team';

interface StoredTeamSnapshot {
  id: string;
  name: string;
  members: Array<{
    id: string;
    name: string;
    email: string;
    role: 'LEADER' | 'MEMBER';
  }>;
}

function readTeamSnapshot(): StoredTeamSnapshot | null {
  const raw = localStorage.getItem(TEAM_SNAPSHOT_KEY);
  if (!raw) return null;

  try {
    return JSON.parse(raw) as StoredTeamSnapshot;
  } catch {
    localStorage.removeItem(TEAM_SNAPSHOT_KEY);
    return null;
  }
}

export const workspaceStorage = {
  getTeamId: () => localStorage.getItem(TEAM_KEY),
  setTeamId: (teamId: string) => localStorage.setItem(TEAM_KEY, teamId),
  clearTeamId: () => localStorage.removeItem(TEAM_KEY),
  getProjectId: () => localStorage.getItem(PROJECT_KEY),
  setProjectId: (projectId: string) => localStorage.setItem(PROJECT_KEY, projectId),
  clearProjectId: () => localStorage.removeItem(PROJECT_KEY),
  getTeamSnapshot: () => readTeamSnapshot(),
  setTeamSnapshot: (team: StoredTeamSnapshot) => localStorage.setItem(TEAM_SNAPSHOT_KEY, JSON.stringify(team)),
  clearTeamSnapshot: () => localStorage.removeItem(TEAM_SNAPSHOT_KEY),
  clear: () => {
    localStorage.removeItem(TEAM_KEY);
    localStorage.removeItem(PROJECT_KEY);
    localStorage.removeItem(TEAM_SNAPSHOT_KEY);
  },
};
