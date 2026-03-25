import { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '@/features/auth/components/auth-provider';
import { loginSchema, type LoginSchema } from '@/features/auth/schemas/login.schema';
import { getErrorMessage } from '@/lib/api-error';

export function LoginForm() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [message, setMessage] = useState('');
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginSchema>({ resolver: zodResolver(loginSchema) });

  const onSubmit = async (values: LoginSchema) => {
    try {
      setMessage('');
      await login(values);
      navigate('/');
    } catch (error) {
      setMessage(getErrorMessage(error, 'Não foi possível fazer login.'));
    }
  };

  return (
    <form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
      <div className="space-y-2">
        <Label htmlFor="email">E-mail</Label>
        <Input id="email" placeholder="voce@email.com" {...register('email')} />
        {errors.email && <p className="text-sm text-rose-600">{errors.email.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="password">Senha</Label>
        <Input id="password" type="password" placeholder="******" {...register('password')} />
        {errors.password && <p className="text-sm text-rose-600">{errors.password.message}</p>}
      </div>
      {message && <p className="rounded-xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">{message}</p>}
      <Button className="w-full" type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Entrando...' : 'Entrar'}
      </Button>
    </form>
  );
}
