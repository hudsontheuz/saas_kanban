export interface ProjectSettings {
  allowDropTask: boolean;
  leaderApprovalRequired: boolean;
  leaderTaskNeedsApproval: boolean;
  peerApprovalAllowed: boolean;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  active: boolean;
  createdAt: string;
  closedAt?: string;
  settings: ProjectSettings;
}
