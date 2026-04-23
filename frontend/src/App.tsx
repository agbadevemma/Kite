import './App.css';
import Login from './pages/Login';
import Signup from './pages/Signup';
import NotFound from './pages/NotFound';
import { TooltipProvider } from '@/components/ui/tooltip';
import { Toaster as Sonner } from '@/components/ui/sonner';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './pages/Dashboard';
import ProtectedRoute from './components/ProtectedRoute';
import Deposit from './pages/Deposit';
import Payout from './pages/Payout';
import Convert from './pages/Convert';
import Transactions from './pages/Transactions';

function App() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: 1, refetchOnWindowFocus: false, staleTime: 30_000 },
    },
  });
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <TooltipProvider>
          <Sonner richColors position="top-right" />
          <BrowserRouter>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/signup" element={<Signup />} />
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/deposit"
                element={
                  <ProtectedRoute>
                    <Deposit />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/convert"
                element={
                  <ProtectedRoute>
                    <Convert />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/payout"  
                element={
                  <ProtectedRoute>
                    <Payout />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/transactions"
                element={
                  <ProtectedRoute>
                    <Transactions />
                  </ProtectedRoute>
                }
              />
              <Route path="/dashboard" element={<Navigate to="/" replace />} />
              <Route path="*" element={<NotFound />} />
            </Routes>
          </BrowserRouter>
        </TooltipProvider>
      </QueryClientProvider>
    </>
  );
}

export default App;
