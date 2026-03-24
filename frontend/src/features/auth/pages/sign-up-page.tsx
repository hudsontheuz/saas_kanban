import { Link, Navigate } from 'react-router-dom';
import { AuthShell } from '@/features/auth/components/auth-shell';
import { RegisterForm } from '@/features/auth/components/register-form';
import { useAuth } from '@/features/auth/components/auth-provider';

export function SignUpPage() {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) return <Navigate to="/" replace />;

  return (
    <AuthShell
      title="Criar conta"
      description="Cadastre um usuário e já entre com a sessão pronta."
      footer={<>Já possui conta? <Link className="font-medium text-primary" to="/sign-in">Entrar</Link></>}
    >
      <RegisterForm />
    </AuthShell>
  );
}
