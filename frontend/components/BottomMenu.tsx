import React from 'react';
import { Tab } from '../types';
import { CalendarIcon, SettingsIcon } from './Icons';

interface BottomMenuProps {
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
}

const BottomMenu: React.FC<BottomMenuProps> = ({ activeTab, onTabChange }) => {
  return (
    <div className="fixed bottom-6 left-0 right-0 z-50 flex justify-center pointer-events-none">
      <div className="
        pointer-events-auto
        flex items-center gap-1 p-1.5
        rounded-2xl
        bg-emerald-100/70 dark:bg-[#0a2e23]/90
        backdrop-blur-xl
        border border-emerald-300/30 dark:border-emerald-500/15
        shadow-lg shadow-emerald-900/[0.08] dark:shadow-black/50
        transform transition-all duration-300 hover:scale-105
      ">
        <button
          onClick={() => onTabChange(Tab.SCHEDULE)}
          className={`
            relative flex items-center justify-center w-14 h-12 rounded-xl transition-all duration-300
            ${activeTab === Tab.SCHEDULE 
              ? 'text-white' 
              : 'text-emerald-700/40 dark:text-emerald-400/40 hover:text-emerald-700 dark:hover:text-emerald-300 hover:bg-emerald-200/40 dark:hover:bg-emerald-900/30'}
          `}
        >
           {activeTab === Tab.SCHEDULE && (
             <div className="absolute inset-0 bg-emerald-600 rounded-xl shadow-lg shadow-emerald-500/30 animate-fade-in -z-10" />
           )}
          <CalendarIcon className="w-6 h-6" />
        </button>

        <div className="w-px h-6 bg-emerald-300/40 dark:bg-emerald-700/40 mx-1"></div>

        <button
          onClick={() => onTabChange(Tab.SETTINGS)}
          className={`
            relative flex items-center justify-center w-14 h-12 rounded-xl transition-all duration-300
            ${activeTab === Tab.SETTINGS 
              ? 'text-white' 
              : 'text-emerald-700/40 dark:text-emerald-400/40 hover:text-emerald-700 dark:hover:text-emerald-300 hover:bg-emerald-200/40 dark:hover:bg-emerald-900/30'}
          `}
        >
          {activeTab === Tab.SETTINGS && (
             <div className="absolute inset-0 bg-emerald-600 rounded-xl shadow-lg shadow-emerald-500/30 animate-fade-in -z-10" />
           )}
          <SettingsIcon className="w-6 h-6" />
        </button>
      </div>
    </div>
  );
};

export default BottomMenu;