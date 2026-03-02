import React from 'react';

interface GlassCardProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

const GlassCard: React.FC<GlassCardProps> = ({ children, className = '', onClick }) => {
  return (
    <div 
      onClick={onClick}
      className={`
        relative overflow-hidden
        bg-emerald-100/40 dark:bg-emerald-950/30 
        backdrop-blur-xl 
        border border-emerald-300/25 dark:border-emerald-500/10
        shadow-md shadow-emerald-900/[0.04] dark:shadow-black/20
        rounded-2xl
        transition-all duration-300
        ${onClick ? 'active:scale-[0.98] cursor-pointer' : ''}
        ${className}
      `}
    >
      {/* Shine effect */}
      <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-white/25 to-transparent pointer-events-none opacity-40 dark:opacity-10 rounded-2xl"></div>
      
      <div className="relative z-10">
        {children}
      </div>
    </div>
  );
};

export default GlassCard;