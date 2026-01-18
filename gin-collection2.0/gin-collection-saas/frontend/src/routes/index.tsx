import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Layout } from '../components/Layout';
import { ProtectedRoute } from '../components/ProtectedRoute';

// Pages (lazy loaded for code splitting)
import { lazy, Suspense } from 'react';

const Login = lazy(() => import('../pages/Login'));
const Register = lazy(() => import('../pages/Register'));
const ForgotPassword = lazy(() => import('../pages/ForgotPassword'));
const ResetPassword = lazy(() => import('../pages/ResetPassword'));
const Dashboard = lazy(() => import('../pages/Dashboard'));
const GinList = lazy(() => import('../pages/GinList'));
const GinDetail = lazy(() => import('../pages/GinDetail'));
const GinCreate = lazy(() => import('../pages/GinCreate'));
const Subscription = lazy(() => import('../pages/Subscription'));
const SubscriptionSuccess = lazy(() => import('../pages/SubscriptionSuccess'));
const SubscriptionCancel = lazy(() => import('../pages/SubscriptionCancel'));
const Settings = lazy(() => import('../pages/Settings'));
const Users = lazy(() => import('../pages/Users'));

// Loading component
const LoadingFallback = () => (
  <div className="flex items-center justify-center min-h-screen">
    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
  </div>
);

// Wrapper component for lazy loaded pages
const LazyPage = ({ children }: { children: React.ReactNode }) => (
  <Suspense fallback={<LoadingFallback />}>{children}</Suspense>
);

export const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <LazyPage>
        <Login />
      </LazyPage>
    ),
  },
  {
    path: '/register',
    element: (
      <LazyPage>
        <Register />
      </LazyPage>
    ),
  },
  {
    path: '/forgot-password',
    element: (
      <LazyPage>
        <ForgotPassword />
      </LazyPage>
    ),
  },
  {
    path: '/reset-password',
    element: (
      <LazyPage>
        <ResetPassword />
      </LazyPage>
    ),
  },
  {
    path: '/subscription/success',
    element: (
      <ProtectedRoute>
        <LazyPage>
          <SubscriptionSuccess />
        </LazyPage>
      </ProtectedRoute>
    ),
  },
  {
    path: '/subscription/cancel',
    element: (
      <ProtectedRoute>
        <LazyPage>
          <SubscriptionCancel />
        </LazyPage>
      </ProtectedRoute>
    ),
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <Layout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: 'dashboard',
        element: (
          <LazyPage>
            <Dashboard />
          </LazyPage>
        ),
      },
      {
        path: 'gins',
        children: [
          {
            index: true,
            element: (
              <LazyPage>
                <GinList />
              </LazyPage>
            ),
          },
          {
            path: 'new',
            element: (
              <LazyPage>
                <GinCreate />
              </LazyPage>
            ),
          },
          {
            path: ':id',
            element: (
              <LazyPage>
                <GinDetail />
              </LazyPage>
            ),
          },
        ],
      },
      {
        path: 'subscription',
        element: (
          <LazyPage>
            <Subscription />
          </LazyPage>
        ),
      },
      {
        path: 'settings',
        element: (
          <LazyPage>
            <Settings />
          </LazyPage>
        ),
      },
      {
        path: 'users',
        element: (
          <LazyPage>
            <Users />
          </LazyPage>
        ),
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/dashboard" replace />,
  },
]);
