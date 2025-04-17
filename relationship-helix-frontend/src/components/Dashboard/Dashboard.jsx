import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useRelationship } from '../../context/RelationshipContext';
import HelixAnimation from './HelixAnimation';
import RetroButton from '../UI/RetroButton';
import CRTEffect from '../UI/CRTEffect';
import '../../styles/retro.css';

/**
 * Dashboard Component
 * 
 * Pagina principală a aplicației care conține animația helixului
 * și controalele pentru gestionarea relației.
 */
const Dashboard = () => {
  const { user, logout } = useAuth();
  const { relationship, fetchRelationship, generateInviteCode } = useRelationship();
  const [inviteCode, setInviteCode] = useState('');
  const [codeCopied, setCodeCopied] = useState(false);
  const navigate = useNavigate();
  
  // Obține informații despre relația utilizatorului
  useEffect(() => {
    if (user && user.id) {
      fetchRelationship();
    }
  }, [user, fetchRelationship]);
  
  // Generează un cod nou de invitație
  const handleGenerateCode = async () => {
    const code = await generateInviteCode();
    setInviteCode(code);
    setCodeCopied(false);
  };
  
  // Copiază codul de invitație în clipboard
  const handleCopyCode = () => {
    if (inviteCode) {
      navigator.clipboard.writeText(inviteCode)
        .then(() => {
          setCodeCopied(true);
          setTimeout(() => setCodeCopied(false), 3000);
        })
        .catch(err => console.error('Eroare la copiere:', err));
    }
  };
  
  // Navighează către pagina de partener
  const handleAddPartner = () => {
    navigate('/invite-code');
  };
  
  // Navighează către pagina de setări relație
  const handleRelationshipSettings = () => {
    navigate('/relationship-settings');
  };

  return (
    <div className="dashboard-container retro-bg">
      <CRTEffect />
      
      <header className="dashboard-header">
        <h1 className="retro-text neon-glow">RETRO SYNTH RELATIONSHIP</h1>
        <div className="user-controls">
          <span className="retro-text">Bine ai venit, {user?.username || 'Utilizator'}</span>
          <RetroButton onClick={logout} size="small">Logout</RetroButton>
        </div>
      </header>
      
      <main className="dashboard-content">
        <HelixAnimation />
        
        <div className="dashboard-controls">
          {!relationship || !relationship.id ? (
            <div className="no-relationship-controls">
              <h2 className="retro-text">Nicio relație activă</h2>
              <div className="invite-controls">
                <RetroButton onClick={handleGenerateCode}>Generează Cod de Invitație</RetroButton>
                
                {inviteCode && (
                  <div className="invite-code-display">
                    <div className="code-container">
                      <span className="retro-text invite-code">{inviteCode}</span>
                      <RetroButton 
                        onClick={handleCopyCode} 
                        size="small"
                      >
                        {codeCopied ? 'Copiat!' : 'Copiază'}
                      </RetroButton>
                    </div>
                    <p className="retro-text">Trimite acest cod partenerului tău pentru a vă conecta</p>
                  </div>
                )}
                
                <div className="separator">sau</div>
                
                <RetroButton onClick={handleAddPartner}>
                  Introdu Cod Partener
                </RetroButton>
              </div>
            </div>
          ) : (
            <div className="relationship-controls">
              <h2 className="retro-text">
                Conectat cu {relationship.partnerName || 'Partener'}
              </h2>
              <div className="days-counter">
                <span className="retro-text">
                  {relationship.daysSinceStart} 
                  {relationship.daysSinceStart === 1 ? ' zi' : ' zile'} împreună
                </span>
              </div>
              <RetroButton onClick={handleRelationshipSettings}>
                Setări Relație
              </RetroButton>
            </div>
          )}
        </div>
      </main>
      
      <footer className="dashboard-footer">
        <p className="retro-text">© 2025 RETRO SYNTH RELATIONSHIP</p>
      </footer>
    </div>
  );
};

export default Dashboard;