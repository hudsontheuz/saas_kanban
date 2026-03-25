import * as React from 'react';
import { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { authApi } from '@/features/auth/api/auth.api';
import type { AuthUser, LoginInput, RegisterInput } from '@/features/auth/types/auth.types';
import { storage } from '@/lib/storage';
import { workspaceStorage } from '@/lib/workspace-storage';

interface AuthContextValue {
  user: AuthUser | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (input: LoginInput) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const storedToken = storage.getToken();
    const storedUser = storage.getUser();
    if (storedToken && storedUser) {
      setToken(storedToken);
      setUser(storedUser);
    }
  }, []);

  const persist = (nextToken: string, nextUser: AuthUser) => {
    storage.setToken(nextToken);
    storage.setUser(nextUser);
    setToken(nextToken);
    setUser(nextUser);
  };

  const login = async (input: LoginInput) => {
    const response = await authApi.login(input);
    persist(response.token, response.user);
  };

  const register = async (input: RegisterInput) => {
    const response = await authApi.register(input);
    persist(response.token, response.user);
  };

  const logout = () => {
    storage.clearToken();
    storage.clearUser();
    workspaceStorage.clear();
    setToken(null);
    setUser(null);
  };

  const value = useMemo(
    () => ({ user, token, isAuthenticated: Boolean(user && token), login, register, logout }),
    [user, token],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used inside AuthProvider');
  return context;
}
