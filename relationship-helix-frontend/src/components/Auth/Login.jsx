import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import RetroButton from '../UI/RetroButton';
import RetroInput from '../UI/RetroInput';
import CRTEffect from '../UI/CRTEffect';
import '../../styles/retro.css';

/**
 * Login Component
 * 
 * Pagina de autentificare a utilizatorilor
 */
const Login = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  
  // Actualizează datele din formular
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };
  
  // Procesează autentificarea
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    
    try {
      await login(formData.email, formData.password);
      navigate('/dashboard');
    } catch (err) {
      console.error('Eroare la autentificare:', err);
      setError(err.message || 'Eroare la autentificare. Verifică datele și încearcă din nou.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container retro-bg">
      <CRTEffect />
      
      <div className="auth-card">
        <div className="auth-header">
          <h1 className="retro-text neon-glow">LOGIN</h1>
        </div>
        
        <form className="auth-form" onSubmit={handleSubmit}>
          {error && <div className="auth-error retro-text">{error}</div>}
          
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
          
          <div className="form-actions">
            <RetroButton type="submit" disabled={loading}>
              {loading ? 'Se încarcă...' : 'Autentificare'}
            </RetroButton>
          </div>
        </form>
        
        <div className="auth-footer">
          <p className="retro-text">
            Nu ai cont? <Link to="/register" className="neon-link">Înregistrează-te</Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;