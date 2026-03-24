import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import { router } from '@/app/router';
import '@/app/globals.css';
import { AuthProvider } from '@/features/auth/components/auth-provider';
import { WorkspaceProvider } from '@/features/workspace/components/workspace-provider';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider>
      <WorkspaceProvider>
        <RouterProvider router={router} />
      </WorkspaceProvider>
    </AuthProvider>
  </React.StrictMode>,
);
