export interface TeamMember {
  id: string;
  name: string;
  email: string;
  role: 'LEADER' | 'MEMBER';
}

export interface Team {
  id: string;
  name: string;
  members: TeamMember[];
}
