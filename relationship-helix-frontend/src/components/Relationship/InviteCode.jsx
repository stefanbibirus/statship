import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRelationship } from '../../context/RelationshipContext';
import RetroButton from '../UI/RetroButton';
import RetroInput from '../UI/RetroInput';
import CRTEffect from '../UI/CRTEffect';
import '../../styles/retro.css';

/**
 * InviteCode Component
 * 
 * Pagina pentru utilizarea unui cod de invitație pentru a stabili o relație
 */
const InviteCode = () => {
  const [inviteCode, setInviteCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { useInviteCode } = useRelationship();
  const navigate = useNavigate();
  
  // Actualizează codul de invitație
  const handleChange = (e) => {
    setInviteCode(e.target.value.toUpperCase());
  };
  
  // Verifică și utilizează codul de invitație
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    
    if (!inviteCode.trim()) {
      return setError('Te rugăm să introduci codul de invitație');
    }
    
    setLoading(true);
    
    try {
      await useInviteCode(inviteCode);
      navigate('/dashboard');
    } catch (err) {
      console.error('Eroare la utilizarea codului de invitație:', err);
      setError(err.message || 'Cod invalid sau expirat. Verifică și încearcă din nou.');
    } finally {
      setLoading(false);
    }
  };
  
  // Navighează înapoi la dashboard
  const handleBack = () => {
    navigate('/dashboard');
  };

  return (
    <div className="invite-code-container retro-bg">
      <CRTEffect />
      
      <div className="invite-code-card">
        <div className="invite-code-header">
          <h1 className="retro-text neon-glow">COD INVITAȚIE</h1>
        </div>
        
        <form className="invite-code-form" onSubmit={handleSubmit}>
          {error && <div className="invite-code-error retro-text">{error}</div>}
          
          <div className="form-info">
            <p className="retro-text">
              Introdu codul de invitație primit de la partenerul tău pentru a vă conecta.
            </p>
          </div>
          
          <div className="form-group">
            <RetroInput
              type="text"
              placeholder="Cod Invitație"
              value={inviteCode}
              onChange={handleChange}
              maxLength={8}
              className="invite-code-input"
              autoFocus
              required
            />
          </div>
          
          <div className="form-actions">
            <RetroButton type="button" onClick={handleBack} variant="secondary">
              Înapoi
            </RetroButton>
            <RetroButton type="submit" disabled={loading}>
              {loading ? 'Se procesează...' : 'Conectează'}
            </RetroButton>
          </div>
        </form>
        
        <div className="invite-code-footer">
          <p className="retro-text">
            Codurile de invitație sunt valabile 24 de ore.
          </p>
        </div>
      </div>
    </div>
  );
};

export default InviteCode;