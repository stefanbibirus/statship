/**
 * Add this code to your src/services/api.js file
 * and use it to debug network connection issues
 */

import axios from 'axios';

// URL de bază pentru API - with debugging
const BASE_URL = process.env.REACT_APP_API_URL;
console.log('API Base URL:', BASE_URL);

// Verify URL has protocol
function ensureProtocol(url) {
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    return 'https://' + url;
  }
  return url;
}

// Instanță Axios configurată pentru API
export const api = axios.create({
  baseURL: ensureProtocol(BASE_URL),
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add request logging for debugging
api.interceptors.request.use(request => {
  console.log('API Request:', request.method, request.url);
  return request;
});

// Interceptor pentru gestionarea erorilor de răspuns
api.interceptors.response.use(
  response => {
    console.log('API Response Success:', response.status);
    return response;
  },
  error => {
    // Log detailed error information
    console.error('API Response Error:', {
      status: error.response?.status,
      statusText: error.response?.statusText,
      url: error.config?.url,
      method: error.config?.method,
      message: error.message,
      response: error.response?.data
    });
    
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