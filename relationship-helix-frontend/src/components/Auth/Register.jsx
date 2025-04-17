import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import RetroButton from '../UI/RetroButton';
import RetroInput from '../UI/RetroInput';
import CRTEffect from '../UI/CRTEffect';
import '../../styles/retro.css';

/**
 * Register Component
 * 
 * Pagina de înregistrare pentru utilizatori noi
 */
const Register = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();
  
  // Actualizează datele din formular
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };
  
  // Procesează înregistrarea
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    
    // Validare parole
    if (formData.password !== formData.confirmPassword) {
      return setError('Parolele nu se potrivesc');
    }
    
    setLoading(true);
    
    try {
      await register(formData.username, formData.email, formData.password);
      navigate('/dashboard');
    } catch (err) {
      console.error('Eroare la înregistrare:', err);
      setError(err.message || 'Eroare la înregistrare. Verifică datele și încearcă din nou.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container retro-bg">
      <CRTEffect />
      
      <div className="auth-card">
        <div className="auth-header">
          <h1 className="retro-text neon-glow">ÎNREGISTRARE</h1>
        </div>
        
        <form className="auth-form" onSubmit={handleSubmit}>
          {error && <div className="auth-error retro-text">{error}</div>}
          
          <div className="form-group">
            <RetroInput
              type="text"
              name="username"
              placeholder="Nume utilizator"
              value={formData.username}
              onChange={handleChange}
              required
            />
          </div>
          
          <div className="form-group">
            <RetroInput
              type="email"
              name="email"
              placeholder="E-mail"
              value={formData.email}
              onChange={handleChange}
              required
            />
          </div>
          
          <div className="form-group">
            <RetroInput
              type="password"
              name="password"
              placeholder="Parolă"
              value={formData.password}
              onChange={handleChange}
              required
            />
          </div>
          
          <div className="form-group">
            <RetroInput
              type="password"
              name="confirmPassword"
              placeholder="Confirmă parola"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
            />
          </div>
          
          <div className="form-actions">
            <RetroButton type="submit" disabled={loading}>
              {loading ? 'Se încarcă...' : 'Înregistrare'}
            </RetroButton>
          </div>
        </form>
        
        <div className="auth-footer">
          <p className="retro-text">
            Ai deja cont? <Link to="/login" className="neon-link">Autentifică-te</Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Register;