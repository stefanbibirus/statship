let socket = null;

/**
 * Configurează conexiunea WebSocket
 * 
 * @param {string} relationshipId ID-ul relației
 * @param {Function} onPositionUpdate Callback pentru actualizări de poziție
 */
export const setupWebsocket = (relationshipId, onPositionUpdate) => {
  // Închide orice conexiune existentă
  if (socket) {
    closeWebsocket();
  }
  
  // Obține token-ul de autentificare
  const token = localStorage.getItem('auth_token');
  
  if (!token) {
    console.error('Autentificare necesară pentru WebSocket');
    return;
  }
  
  // Construiește URL-ul WebSocket
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsHost = process.env.REACT_APP_WS_HOST || window.location.hostname + ':8080';
  const wsUrl = `${wsProtocol}//${wsHost}/ws/relationship/${relationshipId}?token=${token}`;
  
  // Creează conexiunea WebSocket
  socket = new WebSocket(wsUrl);
  
  // Handler pentru deschiderea conexiunii
  socket.onopen = () => {
    console.log('Conexiune WebSocket stabilită');
  };
  
  // Handler pentru primirea mesajelor
  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      
      // Verifică tipul mesajului
      if (data.type === 'position_update') {
        onPositionUpdate(data.payload);
      }
    } catch (error) {
      console.error('Eroare la procesarea mesajului WebSocket:', error);
    }
  };
  
  // Handler pentru erori
  socket.onerror = (error) => {
    console.error('Eroare WebSocket:', error);
  };
  
  // Handler pentru închiderea conexiunii
  socket.onclose = (event) => {
    console.log(`Conexiune WebSocket închisă: ${event.code} ${event.reason}`);
    
    // Reconectare automată după 5 secunde în cazul unei închideri anormale
    if (event.code !== 1000) {
      setTimeout(() => {
        setupWebsocket(relationshipId, onPositionUpdate);
      }, 5000);
    }
  };
  
  return socket;
};

/**
 * Închide conexiunea WebSocket
 */
export const closeWebsocket = () => {
  if (socket) {
    socket.close(1000, 'Închidere normală');
    socket = null;
  }
};