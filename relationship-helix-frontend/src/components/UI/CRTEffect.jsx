import React from 'react';
import '../../styles/retro.css';

/**
 * CRTEffect Component
 * 
 * Adaugă efecte vizuale de ecran CRT pentru a îmbunătăți aspectul retro
 */
const CRTEffect = () => {
  return (
    <div className="crt-effect">
      <div className="crt-scanlines"></div>
      <div className="crt-glow"></div>
      <div className="crt-flicker"></div>
    </div>
  );
};

export default CRTEffect;