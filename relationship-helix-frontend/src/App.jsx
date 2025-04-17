import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

import { AuthProvider, useAuth } from './context/AuthContext';
import { RelationshipProvider } from './context/RelationshipContext';

import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Dashboard from './components/Dashboard/Dashboard';
import InviteCode from './components/Relationship/InviteCode';
import RelationshipSettings from './components/Relationship/RelationshipSettings';

import './styles/retro.css';

// ProtectedRoute - redirecționează la login dacă utilizatorul nu este autentificat
const ProtectedRoute = ({ children }) => {
  const { user, loading } = useAuth();
  
  // În timpul verificării autentificării, nu afișăm nimic
  if (loading) {
    return <div className="loading-screen retro-bg">
      <h1 className="retro-text neon-glow">Se încarcă...</h1>
    </div>;
  }
  
  // Dacă utilizatorul nu este autentificat, redirecționăm la login
  if (!user) {
    return <Navigate to="/login" />;
  }
  
  // Dacă utilizatorul este autentificat, afișăm conținutul protejat
  return children;
};

// PublicRoute - redirecționează la dashboard dacă utilizatorul este deja autentificat
const PublicRoute = ({ children }) => {
  const { user, loading } = useAuth();
  
  // În timpul verificării autentificării, nu afișăm nimic
  if (loading) {
    return <div className="loading-screen retro-bg">
      <h1 className="retro-text neon-glow">Se încarcă...</h1>
    </div>;
  }
  
  // Dacă utilizatorul este autentificat, redirecționăm la dashboard
  if (user) {
    return <Navigate to="/dashboard" />;
  }
  
  // Dacă utilizatorul nu este autentificat, afișăm conținutul public
  return children;
};

function App() {
  return (
    <Router>
      <AuthProvider>
        <RelationshipProvider>
          <Routes>
            {/* Rute publice */}
            <Route path="/login" element={
              <PublicRoute>
                <Login />
              </PublicRoute>
            } />
            <Route path="/register" element={
              <PublicRoute>
                <Register />
              </PublicRoute>
            } />
            
            {/* Rute protejate */}
            <Route path="/dashboard" element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } />
            <Route path="/invite-code" element={
              <ProtectedRoute>
                <InviteCode />
              </ProtectedRoute>
            } />
            <Route path="/relationship-settings" element={
              <ProtectedRoute>
                <RelationshipSettings />
              </ProtectedRoute>
            } />
            
            {/* Redirecționare implicită */}
            <Route path="/" element={<Navigate to="/login" />} />
            <Route path="*" element={<Navigate to="/login" />} />
          </Routes>
        </RelationshipProvider>
      </AuthProvider>
    </Router>
  );
}

export default App;