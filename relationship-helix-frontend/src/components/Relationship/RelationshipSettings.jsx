import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRelationship } from '../../context/RelationshipContext';
import RetroButton from '../UI/RetroButton';
import CRTEffect from '../UI/CRTEffect';
import '../../styles/retro.css';

/**
 * RelationshipSettings Component
 * 
 * Pagina pentru gestionarea setărilor relației
 */
const RelationshipSettings = () => {
  const { relationship, deleteRelationship, fetchRelationship } = useRelationship();
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  
  // Obține informații despre relația curentă
  useEffect(() => {
    fetchRelationship();
  }, [fetchRelationship]);
  
  // Gestionează ștergerea relației
  const handleDeleteRelationship = async () => {
    if (!confirmDelete) {
      return setConfirmDelete(true);
    }
    
    setLoading(true);
    
    try {
      await deleteRelationship();
      navigate('/dashboard');
    } catch (err) {
      console.error('Eroare la ștergerea relației:', err);
    } finally {
      setLoading(false);
      setConfirmDelete(false);
    }
  };
  
  // Anulează ștergerea
  const handleCancelDelete = () => {
    setConfirmDelete(false);
  };
  
  // Navighează înapoi la dashboard
  const handleBack = () => {
    navigate('/dashboard');
  };

  // Dacă nu există o relație, redirecționăm la dashboard
  if (!relationship || !relationship.id) {
    return (
      <div className="relationship-settings-container retro-bg">
        <CRTEffect />
        <div className="relationship-settings-card">
          <h2 className="retro-text">Nicio relație activă</h2>
          <RetroButton onClick={handleBack}>
            Înapoi la Dashboard
          </RetroButton>
        </div>
      </div>
    );
  }

  return (
    <div className="relationship-settings-container retro-bg">
      <CRTEffect />
      
      <div className="relationship-settings-card">
        <div className="relationship-settings-header">
          <h1 className="retro-text neon-glow">SETĂRI RELAȚIE</h1>
        </div>
        
        <div className="relationship-settings-content">
          <div className="relationship-info">
            <h2 className="retro-text">
              Relația cu {relationship.partnerName || 'Partener'}
            </h2>
            <p className="retro-text">
              Conectați de {relationship.daysSinceStart} 
              {relationship.daysSinceStart === 1 ? ' zi' : ' zile'}
            </p>
          </div>
          
          <div className="relationship-actions">
            {!confirmDelete ? (
              <RetroButton 
                onClick={handleDeleteRelationship}
                variant="danger"
              >
                Încheie Relația
              </RetroButton>
            ) : (
              <div className="confirm-delete">
                <p className="retro-text warning">Ești sigur că vrei să închei această relație?</p>
                <div className="confirm-actions">
                  <RetroButton 
                    onClick={handleCancelDelete}
                    variant="secondary"
                  >
                    Anulează
                  </RetroButton>
                  <RetroButton 
                    onClick={handleDeleteRelationship}
                    variant="danger"
                    disabled={loading}
                  >
                    {loading ? 'Se procesează...' : 'Confirmă Încheierea'}
                  </RetroButton>
                </div>
              </div>
            )}
          </div>
        </div>
        
        <div className="relationship-settings-footer">
          <RetroButton onClick={handleBack} variant="secondary">
            Înapoi la Dashboard
          </RetroButton>
        </div>
      </div>
    </div>
  );
};

export default RelationshipSettings;