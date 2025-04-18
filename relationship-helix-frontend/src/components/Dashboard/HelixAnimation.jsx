import React, { useRef, useEffect, useState, useCallback } from 'react';
import { useRelationship } from '../../context/RelationshipContext';
import '../../styles/animations.css';

/**
 * HelixAnimation Component
 * 
 * Acest component renderează animația cu dubla elice (double helix)
 * care reprezintă vizual relația dintre doi utilizatori.
 */
const HelixAnimation = () => {
  const canvasRef = useRef(null);
  const { relationship, userCurvePosition, partnerCurvePosition, updateUserPosition } = useRelationship();
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const animationRef = useRef(null);

  // Configurare pentru animație
  const config = {
    radius: 30,           // Raza helixului
    verticalSpacing: 20,  // Spațiul vertical între spirale
    speed: 0.5,           // Viteza de rotație
    lineWidth: 4,         // Grosimea liniei
    userColor: '#ff00ff', // Neon roz
    partnerColor: '#00ffff', // Neon cyan
    maxDistance: 40,      // Distanța maximă de îndepărtare
    cycles: 6             // Numărul de cicluri complete
  };

  // Funcție pentru redimensionarea canvas-ului
  const handleResize = () => {
    if (canvasRef.current) {
      const canvas = canvasRef.current;
      canvas.width = window.innerWidth * 0.8;
      canvas.height = window.innerHeight * 0.6;
      setDimensions({ width: canvas.width, height: canvas.height });
    }
  };

  // Inițializare și ascultător de redimensionare
  useEffect(() => {
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Funcția de desenare a helixului
  const drawHelix = useCallback((ctx, timestamp) => {
    const { width, height } = dimensions;
    ctx.clearRect(0, 0, width, height);
    
    const centerX = width / 2;
    const centerY = height / 2;
    const time = timestamp * 0.001 * config.speed;
    
    // Calculăm distanța utilizatorului actual - convertim din procente (0-100) în pixeli
    const userDistance = (config.maxDistance * userCurvePosition) / 100;
    
    // Calculăm distanța partenerului
    const partnerDistance = (config.maxDistance * partnerCurvePosition) / 100;
    
    // Calculăm înălțimea totală a helixului
    const helixHeight = config.cycles * config.verticalSpacing;
    const startY = centerY - helixHeight / 2;
    
    // Desenăm curba utilizatorului curent (curba 1)
    ctx.beginPath();
    ctx.strokeStyle = config.userColor;
    ctx.lineWidth = config.lineWidth;
    ctx.shadowBlur = 10;
    ctx.shadowColor = config.userColor;
    
    for (let i = 0; i <= config.cycles * 10; i += 0.1) {
      const progress = i / (config.cycles * 10);
      const y = startY + progress * helixHeight;
      
      // Formula pentru helix vertical cu distanțare
      const x = centerX + Math.sin(i + time) * (config.radius - userDistance);
      
      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }
    ctx.stroke();
    
    // Desenăm curba partenerului (curba 2)
    ctx.beginPath();
    ctx.strokeStyle = config.partnerColor;
    ctx.shadowColor = config.partnerColor;
    
    for (let i = 0; i <= config.cycles * 10; i += 0.1) {
      const progress = i / (config.cycles * 10);
      const y = startY + progress * helixHeight;
      
      // Formula pentru helix vertical cu distanțare, defazat cu 180 grade (PI radiani)
      const x = centerX + Math.sin(i + time + Math.PI) * (config.radius - partnerDistance);
      
      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }
    ctx.stroke();
    
    // Adăugăm efect de scanline retro
    drawScanlines(ctx, width, height);
    
    animationRef.current = requestAnimationFrame((timestamp) => drawHelix(ctx, timestamp));
  }, [dimensions, userCurvePosition, partnerCurvePosition]);

  // Funcție pentru desenarea scanline-urilor (efect CRT)
  const drawScanlines = (ctx, width, height) => {
    ctx.globalAlpha = 0.1;
    ctx.fillStyle = '#000000';
    
    for (let y = 0; y < height; y += 4) {
      ctx.fillRect(0, y, width, 2);
    }
    
    ctx.globalAlpha = 1.0;
  };

  // Efect pentru gestionarea animației
  useEffect(() => {
    if (!canvasRef.current || !dimensions.width) return;
    
    const ctx = canvasRef.current.getContext('2d');
    
    // Pornește animația
    animationRef.current = requestAnimationFrame((timestamp) => drawHelix(ctx, timestamp));
    
    // Curăță animația la demontat
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [dimensions, drawHelix]);

  // Dacă nu există o relație, afișăm un mesaj
  if (!relationship || !relationship.id) {
    return (
      <div className="helix-placeholder retro-text">
        <h2>Nicio relație activă</h2>
        <p>Adaugă un partener folosind codul de invitație pentru a vedea animația</p>
      </div>
    );
  }

  return (
    <div className="helix-container">
      <canvas ref={canvasRef} className="helix-canvas"></canvas>
      <div className="helix-controls">
        <div className="user-info">
          <h3>Tu</h3>
          <div className="distance-control">
            <button 
              className="retro-button"
              onClick={() => updateUserPosition(Math.max(0, userCurvePosition - 10))}
            >
              Apropie
            </button>
            <button 
              className="retro-button"
              onClick={() => updateUserPosition(Math.min(100, userCurvePosition + 10))}
            >
              Îndepărtează
            </button>
          </div>
        </div>
        <div className="partner-info">
          <h3>{relationship.partnerName || 'Partener'}</h3>
          <div className="distance-indicator">
            {partnerCurvePosition < 30 ? (
              <span className="status close">Aproape</span>
            ) : partnerCurvePosition < 70 ? (
              <span className="status neutral">Neutru</span>
            ) : (
              <span className="status distant">Distant</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default HelixAnimation;