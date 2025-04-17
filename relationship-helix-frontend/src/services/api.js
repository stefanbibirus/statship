import axios from 'axios';

// URL de bază pentru API
const BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Instanță Axios configurată pentru API
export const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Interceptor pentru gestionarea erorilor de răspuns
api.interceptors.response.use(
  response => response,
  error => {
    // Gestionează erorile de autorizare (401)
    if (error.response && error.response.status === 401) {
      // Șterge token-ul din localStorage
      localStorage.removeItem('auth_token');
      
      // Șterge token-ul din header-ele de autentificare
      delete api.defaults.headers.common['Authorization'];
      
      // Redirecționează la pagina de login dacă nu suntem deja acolo
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);