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
        description="Ajuste as regras disponíveis para o fluxo do projeto ativo."
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
            <CardTitle>Fluxo de trabalho</CardTitle>
            <CardDescription>Escolha como a equipe deve tratar tarefas que já estão em andamento.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center justify-between rounded-xl border p-4">
              <div>
                <Label>Permitir retorno para To Do</Label>
                <p className="text-sm text-slate-500">Define se uma tarefa em andamento pode voltar para a etapa inicial.</p>
              </div>
              <Switch checked={project.settings.allowDropTask} onCheckedChange={(checked) => void updateProjectSettings({ ...project.settings, allowDropTask: checked })} />
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
