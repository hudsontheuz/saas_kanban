import { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '@/features/auth/components/auth-provider';
import { registerSchema, type RegisterSchema } from '@/features/auth/schemas/register.schema';
import { getErrorMessage } from '@/lib/api-error';

export function RegisterForm() {
  const navigate = useNavigate();
  const { register: signUp } = useAuth();
  const [message, setMessage] = useState('');
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterSchema>({ resolver: zodResolver(registerSchema) });

  const onSubmit = async (values: RegisterSchema) => {
    try {
      setMessage('');
      await signUp(values);
      navigate('/');
    } catch (error) {
      setMessage(getErrorMessage(error, 'Não foi possível criar a conta.'));
    }
  };

  return (
    <form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
      <div className="space-y-2">
        <Label htmlFor="name">Nome</Label>
        <Input id="name" placeholder="Seu nome" {...register('name')} />
        {errors.name && <p className="text-sm text-rose-600">{errors.name.message}</p>}
      </div>
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
        {isSubmitting ? 'Criando conta...' : 'Criar conta'}
      </Button>
    </form>
  );
}
