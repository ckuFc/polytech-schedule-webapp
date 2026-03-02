import React from 'react';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="relative h-[100dvh] w-full overflow-hidden bg-[#e0eddf] dark:bg-[#022c22] text-gray-800 dark:text-gray-100 font-sans selection:bg-emerald-500 selection:text-white transition-colors duration-300">
      {/* Background gradient */}
      <div className="fixed inset-0 w-full h-full pointer-events-none z-0 bg-gradient-to-br from-emerald-200/20 via-emerald-100/10 to-emerald-200/15 dark:from-black/40 dark:via-emerald-950/20 dark:to-black/60"></div>
      
      {/* Decorative Blur Blobs */}
      <div className="fixed top-[-10%] left-[-10%] w-[50%] h-[40%] bg-emerald-400/10 dark:bg-emerald-600/10 rounded-full blur-[80px] pointer-events-none"></div>
      <div className="fixed bottom-[-10%] right-[-10%] w-[50%] h-[40%] bg-emerald-300/10 dark:bg-emerald-500/10 rounded-full blur-[80px] pointer-events-none"></div>

      {/* Content Container */}
      <div className="relative z-10 w-full h-full flex flex-col max-w-md mx-auto">
         {children}
      </div>
    </div>
  );
};

export default Layout;