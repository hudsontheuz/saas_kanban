import { PageHeader } from '@/components/layout/page-header';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { useWorkspace } from '@/hooks/use-workspace';

export function SettingsPage() {
  const { project, updateProjectSettings, isLoading, message, clearMessage } = useWorkspace();

  return (
    <div className="space-y-6">
      <PageHeader
        title="Configurações do projeto"
        description="Essas opções agora estão apontadas para endpoints reais de configuração do projeto."
      />

      {message && (
        <button type="button" onClick={clearMessage} className="w-full rounded-xl border bg-white px-4 py-3 text-left text-sm text-slate-600 shadow-soft">
          {message}
        </button>
      )}

      {!project && !isLoading && (
        <Card>
          <CardHeader>
            <CardTitle>Nenhum projeto carregado</CardTitle>
            <CardDescription>Crie um projeto antes de ajustar as regras do fluxo.</CardDescription>
          </CardHeader>
        </Card>
      )}

      {project && (
        <Card>
          <CardHeader>
            <CardTitle>Regras do fluxo</CardTitle>
            <CardDescription>Quando o backend expor o endpoint de settings, essa tela já estará pronta para consumir.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center justify-between rounded-xl border p-4">
              <div>
                <Label>Permitir soltar tarefa</Label>
                <p className="text-sm text-slate-500">Decide se o usuário pode devolver uma tarefa de Doing para ToDo.</p>
              </div>
              <Switch checked={project.settings.allowDropTask} onCheckedChange={(checked) => void updateProjectSettings({ ...project.settings, allowDropTask: checked })} />
            </div>
            <div className="flex items-center justify-between rounded-xl border p-4">
              <div>
                <Label>Aprovação do líder obrigatória</Label>
                <p className="text-sm text-slate-500">Controla se tarefas em revisão exigem aprovação do líder.</p>
              </div>
              <Switch checked={project.settings.leaderApprovalRequired} onCheckedChange={(checked) => void updateProjectSettings({ ...project.settings, leaderApprovalRequired: checked })} />
            </div>
            <div className="flex items-center justify-between rounded-xl border p-4">
              <div>
                <Label>Tarefas do líder exigem aprovação</Label>
                <p className="text-sm text-slate-500">Permite flexibilizar o fluxo quando o líder executa a própria tarefa.</p>
              </div>
              <Switch checked={project.settings.leaderTaskNeedsApproval} onCheckedChange={(checked) => void updateProjectSettings({ ...project.settings, leaderTaskNeedsApproval: checked })} />
            </div>
            <div className="flex items-center justify-between rounded-xl border p-4">
              <div>
                <Label>Aprovação por qualquer membro</Label>
                <p className="text-sm text-slate-500">Base para um time mais sênior, onde a aprovação não depende de uma única pessoa.</p>
              </div>
              <Switch checked={project.settings.peerApprovalAllowed} onCheckedChange={(checked) => void updateProjectSettings({ ...project.settings, peerApprovalAllowed: checked })} />
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
