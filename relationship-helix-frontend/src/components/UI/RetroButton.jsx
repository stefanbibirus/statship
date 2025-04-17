import React from 'react';
import '../../styles/retro.css';

/**
 * RetroButton Component
 * 
 * Buton stilizat în tema retro pentru aplicație
 */
const RetroButton = ({ 
  children, 
  onClick, 
  type = 'button', 
  disabled = false, 
  variant = 'primary', 
  size = 'medium',
  className = ''
}) => {
  const variantClass = `retro-button-${variant}`;
  const sizeClass = `retro-button-${size}`;
  
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`retro-button ${variantClass} ${sizeClass} ${className} ${disabled ? 'retro-button-disabled' : ''}`}
    >
      {children}
    </button>
  );
};

export default RetroButton;