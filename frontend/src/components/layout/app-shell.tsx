import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import { FolderKanban, LayoutDashboard, LogOut, Settings, ShieldCheck, Users } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { useAuth } from '@/features/auth/components/auth-provider';
import { useWorkspace } from '@/hooks/use-workspace';

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/projects', label: 'Projetos', icon: FolderKanban },
  { to: '/tasks', label: 'Tarefas', icon: ShieldCheck },
  { to: '/team', label: 'Equipe', icon: Users },
  { to: '/settings', label: 'Configurações', icon: Settings },
];

export function AppShell() {
  const { user, logout } = useAuth();
  const { team, project } = useWorkspace();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/sign-in');
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="mx-auto grid min-h-screen max-w-7xl grid-cols-1 gap-6 p-4 lg:grid-cols-[260px_1fr] lg:p-6">
        <aside className="rounded-2xl border bg-white p-4 shadow-soft">
          <Link to="/" className="mb-8 flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground">
              <FolderKanban className="h-5 w-5" />
            </div>
            <div>
              <p className="font-semibold">SaaS Kanban</p>
              <p className="text-sm text-slate-500">Gestão de tarefas</p>
            </div>
          </Link>

          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              return (
                <NavLink
                  key={item.to}
                  to={item.to}
                  end={item.to === '/'}
                  className={({ isActive }) =>
                    cn(
                      'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors',
                      isActive ? 'bg-primary text-primary-foreground' : 'text-slate-600 hover:bg-slate-100',
                    )
                  }
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </NavLink>
              );
            })}
          </nav>

          <div className="mt-8 space-y-3 rounded-2xl bg-slate-100 p-4 text-sm">
            <div>
              <p className="font-medium text-slate-900">{user?.name}</p>
              <p className="text-xs text-slate-500">{user?.email}</p>
            </div>
            <div className="space-y-1 border-t pt-3 text-xs text-slate-500">
              <p>Equipe: <span className="font-medium text-slate-700">{team?.name ?? 'Não criada'}</span></p>
              <p>Projeto: <span className="font-medium text-slate-700">{project?.name ?? 'Nenhum ativo'}</span></p>
            </div>
          </div>

          <Button variant="outline" className="mt-4 w-full justify-start gap-2" onClick={handleLogout}>
            <LogOut className="h-4 w-4" />
            Sair
          </Button>
        </aside>

        <main className="space-y-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
