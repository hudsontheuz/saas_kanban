import { apiClient } from '@/lib/api-client';
import type { AuthResponse, LoginInput, RegisterInput } from '@/features/auth/types/auth.types';

function decodeToken(token: string): Record<string, unknown> {
  try {
    const [, payload] = token.split('.');
    if (!payload) return {};
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/');
    const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=');
    const json = atob(padded);
    return JSON.parse(json) as Record<string, unknown>;
  } catch {
    return {};
  }
}

function normalizeAuthResponse(data: unknown, fallbackName?: string, fallbackEmail?: string): AuthResponse {
  const payload = data && typeof data === 'object' ? (data as Record<string, unknown>) : {};
  const token = typeof payload.token === 'string' ? payload.token : '';
  const decoded = token ? decodeToken(token) : {};
  const user = payload.user && typeof payload.user === 'object' ? (payload.user as Record<string, unknown>) : {};

  return {
    token,
    user: {
      id:
        (typeof user.id === 'string' && user.id) ||
        (typeof user.userId === 'string' && user.userId) ||
        (typeof decoded.sub === 'string' && decoded.sub) ||
        crypto.randomUUID(),
      name:
        (typeof user.name === 'string' && user.name) ||
        (typeof user.nome === 'string' && user.nome) ||
        fallbackName ||
        'Usuário',
      email:
        (typeof user.email === 'string' && user.email) ||
        (typeof decoded.email === 'string' && decoded.email) ||
        fallbackEmail ||
        '',
    },
  };
}

export const authApi = {
  async login(input: LoginInput): Promise<AuthResponse> {
    const { data } = await apiClient.post('/auth/login', {
      email: input.email,
      senha: input.password,
      password: input.password,
    });
    return normalizeAuthResponse(data, undefined, input.email);
  },
  async register(input: RegisterInput): Promise<AuthResponse> {
    const { data } = await apiClient.post('/auth/register', {
      nome: input.name,
      email: input.email,
      senha: input.password,
      password: input.password,
    });
    return normalizeAuthResponse(data, input.name, input.email);
  },
};
