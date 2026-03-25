import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AppShell } from '@/components/layout/app-shell';
import { ProtectedRoute } from '@/components/shared/protected-route';
import { SignInPage } from '@/features/auth/pages/sign-in-page';
import { SignUpPage } from '@/features/auth/pages/sign-up-page';
import { DashboardPage } from '@/features/dashboard/pages/dashboard-page';
import { ProjectsPage } from '@/features/project/pages/projects-page';
import { TasksPage } from '@/features/task/pages/tasks-page';
import { TeamPage } from '@/features/team/pages/team-page';
import { SettingsPage } from '@/features/settings/pages/settings-page';

export const router = createBrowserRouter([
  { path: '/', element: <ProtectedRoute><AppShell /></ProtectedRoute>, children: [
    { index: true, element: <DashboardPage /> },
    { path: 'projects', element: <ProjectsPage /> },
    { path: 'tasks', element: <TasksPage /> },
    { path: 'team', element: <TeamPage /> },
    { path: 'settings', element: <SettingsPage /> },
  ] },
  { path: '/sign-in', element: <SignInPage /> },
  { path: '/sign-up', element: <SignUpPage /> },
  { path: '*', element: <Navigate to="/" replace /> },
]);
