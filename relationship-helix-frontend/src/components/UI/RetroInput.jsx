import React from 'react';
import '../../styles/retro.css';

/**
 * RetroInput Component
 * 
 * Input stilizat în tema retro pentru aplicație
 */
const RetroInput = ({ 
  type = 'text', 
  placeholder = '', 
  value, 
  onChange, 
  name,
  required = false,
  maxLength,
  autoFocus = false,
  className = ''
}) => {
  return (
    <div className="retro-input-container">
      <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        name={name}
        required={required}
        maxLength={maxLength}
        autoFocus={autoFocus}
        className={`retro-input ${className}`}
      />
      <div className="retro-input-glow"></div>
    </div>
  );
};

export default RetroInput;