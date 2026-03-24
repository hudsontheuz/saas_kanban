import { Link, Navigate } from 'react-router-dom';
import { AuthShell } from '@/features/auth/components/auth-shell';
import { LoginForm } from '@/features/auth/components/login-form';
import { useAuth } from '@/features/auth/components/auth-provider';

export function SignInPage() {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) return <Navigate to="/" replace />;

  return (
    <AuthShell
      title="Entrar"
      description="Acesse seu workspace para acompanhar projeto, equipe e tarefas."
      footer={<>Não tem conta? <Link className="font-medium text-primary" to="/sign-up">Criar conta</Link></>}
    >
      <LoginForm />
    </AuthShell>
  );
}
