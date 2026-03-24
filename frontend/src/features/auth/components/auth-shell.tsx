import * as React from 'react';
import { Link } from 'react-router-dom';
import { FolderKanban } from 'lucide-react';

export function AuthShell({ title, description, children, footer }: { title: string; description: string; children: React.ReactNode; footer: React.ReactNode }) {
  return (
    <div className="grid min-h-screen place-items-center bg-slate-100 p-4">
      <div className="grid w-full max-w-5xl gap-6 lg:grid-cols-[1.1fr_0.9fr]">
        <div className="hidden rounded-[2rem] bg-slate-900 p-10 text-white shadow-soft lg:block">
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/10">
              <FolderKanban className="h-6 w-6" />
            </div>
            <div>
              <p className="text-lg font-semibold">SaaS Kanban</p>
              <p className="text-sm text-slate-300">Projeto alinhado ao backend em Go</p>
            </div>
          </div>
          <div className="mt-16 max-w-md space-y-4">
            <h2 className="text-4xl font-semibold leading-tight">Frontend limpo, em TypeScript, com shadcn/ui e rotas prontas.</h2>
            <p className="text-slate-300">
              Estrutura modular por feature, foco em auth, projeto, time, tarefas e configurações do fluxo de trabalho.
            </p>
          </div>
        </div>

        <div className="rounded-[2rem] border bg-white p-8 shadow-soft">
          <Link to="/sign-in" className="text-sm font-medium text-primary">SaaS Kanban</Link>
          <div className="mt-6 space-y-2">
            <h1 className="text-3xl font-semibold">{title}</h1>
            <p className="text-sm text-slate-500">{description}</p>
          </div>
          <div className="mt-8">{children}</div>
          <div className="mt-6 text-sm text-slate-500">{footer}</div>
        </div>
      </div>
    </div>
  );
}
