import { useState } from 'react';
import { Crown, Trash2 } from 'lucide-react';
import { PageHeader } from '@/components/layout/page-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { InviteMemberDialog } from '@/features/team/components/invite-member-dialog';
import { useWorkspace } from '@/hooks/use-workspace';

export function TeamPage() {
  const {
    team,
    members,
    isLoading,
    message,
    clearMessage,
    createTeam,
    addMember,
    removeMember,
    transferLeadership,
  } = useWorkspace();
  const [teamName, setTeamName] = useState('');

  return (
    <div className="space-y-6">
      <PageHeader
        title="Equipe"
        description="Gestão real de membros, liderança e limite máximo de 5 pessoas."
        action={<InviteMemberDialog onInvite={addMember} disabled={!team} />}
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {!team && !isLoading && (
        <Card>
          <CardHeader>
            <CardTitle>Crie sua primeira equipe</CardTitle>
            <CardDescription>Depois disso o restante do fluxo passa a conversar com projeto e tarefas do backend.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3 md:flex-row">
            <Input value={teamName} onChange={(event) => setTeamName(event.target.value)} placeholder="Nome da equipe" />
            <Button onClick={() => void createTeam(teamName)}>Criar equipe</Button>
          </CardContent>
        </Card>
      )}

      {team && (
        <Card>
          <CardHeader>
            <CardTitle>{team.name}</CardTitle>
            <CardDescription>{members.length} membro(s) conectado(s) ao workspace atual.</CardDescription>
          </CardHeader>
        </Card>
      )}

      <div className="grid gap-4">
        {members.map((member) => (
          <Card key={member.id}>
            <CardContent className="flex flex-col gap-4 p-4 md:flex-row md:items-center md:justify-between">
              <div>
                <div className="flex items-center gap-2">
                  <p className="font-medium">{member.name}</p>
                  {member.role === 'LEADER' && <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700">Líder</span>}
                </div>
                <p className="text-sm text-slate-500">{member.email}</p>
              </div>
              <div className="flex flex-wrap gap-2">
                {member.role !== 'LEADER' && (
                  <Button variant="outline" onClick={() => void transferLeadership(member.id)}>
                    <Crown className="mr-2 h-4 w-4" /> Transferir liderança
                  </Button>
                )}
                {member.role !== 'LEADER' && (
                  <Button variant="destructive" onClick={() => void removeMember(member.id)}>
                    <Trash2 className="mr-2 h-4 w-4" /> Remover
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
