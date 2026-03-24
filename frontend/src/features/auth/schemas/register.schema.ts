import { z } from 'zod';

export const registerSchema = z.object({
  name: z.string().min(2, 'Informe seu nome'),
  email: z.string().email('Informe um e-mail válido'),
  password: z.string().min(6, 'A senha deve ter ao menos 6 caracteres'),
});

export type RegisterSchema = z.infer<typeof registerSchema>;
