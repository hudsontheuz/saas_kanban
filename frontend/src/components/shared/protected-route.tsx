import * as React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '@/features/auth/components/auth-provider';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/sign-in" replace />;
  return <>{children}</>;
}
