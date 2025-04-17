import React, { createContext, useState, useContext, useEffect } from 'react';
import { api } from '../services/api';
import { login as loginApi, register as registerApi, checkAuth } from '../services/auth';

// Crearea contextului de autentificare
const AuthContext = createContext(null);

// Hook personalizat pentru utilizarea contextului de autentificare
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth trebuie utilizat în cadrul unui AuthProvider');
  }
  return context;
};

/**
 * AuthProvider Component
 * 
 * Furnizor de context pentru autentificare
 */
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [initialized, setInitialized] = useState(false);
  
  // Verifică autentificarea la încărcarea aplicației
  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('auth_token');
      
      if (token) {
        try {
          // Setează token-ul în header-ele de autentificare
          api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
          
          // Verifică dacă token-ul este valid
          const userData = await checkAuth();
          setUser(userData);
        } catch (error) {
          console.error('Token invalid sau expirat', error);
          localStorage.removeItem('auth_token');
          delete api.defaults.headers.common['Authorization'];
        }
      }
      
      setLoading(false);
      setInitialized(true);
    };
    
    initAuth();
  }, []);
  
  // Funcția de autentificare
  const login = async (email, password) => {
    try {
      const { user: userData, token } = await loginApi(email, password);
      
      // Salvează token-ul în localStorage
      localStorage.setItem('auth_token', token);
      
      // Setează token-ul în header-ele de autentificare
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
      // Actualizează starea utilizatorului
      setUser(userData);
      
      return userData;
    } catch (error) {
      throw error;
    }
  };
  
  // Funcția de înregistrare
  const register = async (username, email, password) => {
    try {
      const { user: userData, token } = await registerApi(username, email, password);
      
      // Salvează token-ul în localStorage
      localStorage.setItem('auth_token', token);
      
      // Setează token-ul în header-ele de autentificare
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
      // Actualizează starea utilizatorului
      setUser(userData);
      
      return userData;
    } catch (error) {
      throw error;
    }
  };
  
  // Funcția de delogare
  const logout = () => {
    // Șterge token-ul din localStorage
    localStorage.removeItem('auth_token');
    
    // Șterge token-ul din header-ele de autentificare
    delete api.defaults.headers.common['Authorization'];
    
    // Resetează starea utilizatorului
    setUser(null);
  };
  
  // Valoarea contextului
  const value = {
    user,
    loading,
    initialized,
    login,
    register,
    logout
  };
  
  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};