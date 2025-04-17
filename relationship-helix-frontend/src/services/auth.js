import { api } from './api';

/**
 * Autentifică un utilizator
 * 
 * @param {string} email Adresa de email
 * @param {string} password Parola
 * @returns {Promise} Datele utilizatorului și token-ul
 */
export const login = async (email, password) => {
  try {
    const response = await api.post('/api/auth/login', { email, password });
    return response.data;
  } catch (error) {
    if (error.response && error.response.data && error.response.data.message) {
      throw new Error(error.response.data.message);
    }
    throw error;
  }
};

/**
 * Înregistrează un utilizator nou
 * 
 * @param {string} username Numele de utilizator
 * @param {string} email Adresa de email
 * @param {string} password Parola
 * @returns {Promise} Datele utilizatorului și token-ul
 */
export const register = async (username, email, password) => {
  try {
    const response = await api.post('/api/auth/register', { username, email, password });
    return response.data;
  } catch (error) {
    if (error.response && error.response.data && error.response.data.message) {
      throw new Error(error.response.data.message);
    }
    throw error;
  }
};

/**
 * Verifică autentificarea utilizatorului curent
 * 
 * @returns {Promise} Datele utilizatorului
 */
export const checkAuth = async () => {
  try {
    const response = await api.get('/api/auth/me');
    return response.data.user;
  } catch (error) {
    throw error;
  }
};