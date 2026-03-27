import { useState } from 'react';
import { Crown } from 'lucide-react';
import { PageHeader } from '@/components/layout/page-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useWorkspace } from '@/hooks/use-workspace';

export function TeamPage() {
  const { team, members, isLoading, message, clearMessage, createTeam } = useWorkspace();
  const [teamName, setTeamName] = useState('');

  return (
    <div className="space-y-6">
      <PageHeader
        title="Equipe"
        description="Crie a equipe responsável pelo projeto e acompanhe o responsável pelo workspace."
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {!team && !isLoading && (
        <Card>
          <CardHeader>
            <CardTitle>Criar equipe</CardTitle>
            <CardDescription>Defina o nome da equipe para começar a organizar o projeto.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3 md:flex-row">
            <Input value={teamName} onChange={(event) => setTeamName(event.target.value)} placeholder="Nome da equipe" />
            <Button onClick={() => void createTeam(teamName)} disabled={!teamName.trim()}>Criar equipe</Button>
          </CardContent>
        </Card>
      )}

      {team && (
        <>
          <Card>
            <CardHeader>
              <CardTitle>{team.name}</CardTitle>
              <CardDescription>Equipe atual vinculada ao workspace.</CardDescription>
            </CardHeader>
          </Card>

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
                  {member.role === 'LEADER' && (
                    <div className="inline-flex items-center gap-2 rounded-xl bg-slate-100 px-3 py-2 text-sm text-slate-600">
                      <Crown className="h-4 w-4 text-amber-600" /> Responsável pela equipe
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
