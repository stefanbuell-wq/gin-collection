import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';
import { router } from './routes';
import { useAuthStore } from './stores/authStore';

function App() {
  const initializeAuth = useAuthStore((state) => state.initializeAuth);

  // Initialize auth and fetch CSRF token on app load
  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  return <RouterProvider router={router} />;
}

export default App;
