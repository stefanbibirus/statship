import React, { createContext, useState, useContext, useCallback } from 'react';
import { api } from '../services/api';
import { setupWebsocket, closeWebsocket } from '../services/websocket';

// Crearea contextului de relație
const RelationshipContext = createContext(null);

// Hook personalizat pentru utilizarea contextului de relație
export const useRelationship = () => {
  const context = useContext(RelationshipContext);
  if (!context) {
    throw new Error('useRelationship trebuie utilizat în cadrul unui RelationshipProvider');
  }
  return context;
};

/**
 * RelationshipProvider Component
 * 
 * Furnizor de context pentru relația utilizatorului
 */
export const RelationshipProvider = ({ children }) => {
  const [relationship, setRelationship] = useState(null);
  const [userCurvePosition, setUserCurvePosition] = useState(0);
  const [partnerCurvePosition, setPartnerCurvePosition] = useState(0);
  const [wsConnected, setWsConnected] = useState(false);
  
  // Obține informații despre relația utilizatorului
  const fetchRelationship = useCallback(async () => {
    try {
      const response = await api.get('/api/relationship');
      
      if (response.data && response.data.relationship) {
        setRelationship(response.data.relationship);
        setUserCurvePosition(response.data.userCurvePosition || 0);
        setPartnerCurvePosition(response.data.partnerCurvePosition || 0);
        
        // Inițializează conexiunea WebSocket dacă există o relație
        if (response.data.relationship.id && !wsConnected) {
          initWebsocket(response.data.relationship.id);
        }
      } else {
        setRelationship(null);
        setUserCurvePosition(0);
        setPartnerCurvePosition(0);
        
        // Închide conexiunea WebSocket dacă nu există o relație
        if (wsConnected) {
          closeWebsocket();
          setWsConnected(false);
        }
      }
      
      return response.data;
    } catch (error) {
      console.error('Eroare la obținerea relației:', error);
      throw error;
    }
  }, [wsConnected]);
  
  // Generează un cod de invitație
  const generateInviteCode = async () => {
    try {
      const response = await api.post('/api/relationship/invite');
      return response.data.inviteCode;
    } catch (error) {
      console.error('Eroare la generarea codului de invitație:', error);
      throw error;
    }
  };
  
  // Utilizează un cod de invitație - REDENUMIT pentru claritate
  const joinWithInviteCode = async (inviteCode) => {
    try {
      const response = await api.post('/api/relationship/join', { inviteCode });
      
      if (response.data && response.data.relationship) {
        setRelationship(response.data.relationship);
        setUserCurvePosition(response.data.userCurvePosition || 0);
        setPartnerCurvePosition(response.data.partnerCurvePosition || 0);
        
        // Inițializează conexiunea WebSocket
        initWebsocket(response.data.relationship.id);
      }
      
      return response.data;
    } catch (error) {
      console.error('Eroare la utilizarea codului de invitație:', error);
      throw error;
    }
  };
  
  // Actualizează poziția curbei utilizatorului
  const updateUserPosition = async (position) => {
    try {
      const response = await api.post('/api/relationship/position', { position });
      setUserCurvePosition(position);
      return response.data;
    } catch (error) {
      console.error('Eroare la actualizarea poziției:', error);
      throw error;
    }
  };
  
  // Șterge relația
  const deleteRelationship = async () => {
    try {
      await api.delete('/api/relationship');
      setRelationship(null);
      setUserCurvePosition(0);
      setPartnerCurvePosition(0);
      
      // Închide conexiunea WebSocket
      if (wsConnected) {
        closeWebsocket();
        setWsConnected(false);
      }
    } catch (error) {
      console.error('Eroare la ștergerea relației:', error);
      throw error;
    }
  };
  
  // Inițializează WebSocket pentru actualizări în timp real
  const initWebsocket = (relationshipId) => {
    // Callback pentru actualizarea poziției partenerului
    const handlePositionUpdate = (data) => {
      if (data && data.partnerId && data.position !== undefined) {
        setPartnerCurvePosition(data.position);
      }
    };
    
    // Configurează conexiunea WebSocket
    setupWebsocket(relationshipId, handlePositionUpdate);
    setWsConnected(true);
  };
  
  // Valoarea contextului
  const value = {
    relationship,
    userCurvePosition,
    partnerCurvePosition,
    fetchRelationship,
    generateInviteCode,
    joinWithInviteCode, // REDENUMIT pentru claritate
    updateUserPosition,
    deleteRelationship
  };
  
  return (
    <RelationshipContext.Provider value={value}>
      {children}
    </RelationshipContext.Provider>
  );
};