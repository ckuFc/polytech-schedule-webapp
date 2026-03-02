import React, { useState, useMemo, useRef, useEffect } from 'react';
import GlassCard from '../components/GlassCard';
import { MoonIcon, SunIcon, SearchIcon, ChevronDownIcon } from '../components/Icons';

interface SettingsScreenProps {
  group: string;
  setGroup: (group: string) => void;
  groups: string[];
  theme: 'light' | 'dark';
  toggleTheme: () => void;
  notificationsEnabled: boolean;
  onToggleNotifications: (val: boolean) => void;
}

const SettingsScreen: React.FC<SettingsScreenProps> = ({ 
  group, setGroup, groups, theme, toggleTheme, 
  notificationsEnabled, onToggleNotifications 
}) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const dropdownRef = useRef<HTMLDivElement>(null);

  const filteredGroups = useMemo(() => {
    if (!searchTerm) return groups;
    const q = searchTerm.toLowerCase();
    return groups.filter(g => g.toLowerCase().includes(q));
  }, [searchTerm, groups]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="flex flex-col h-full w-full overflow-hidden">
      <div className="px-6 safe-pt pb-6 shrink-0">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white tracking-tight">Настройки</h1>
      </div>

      <div className="flex-1 overflow-y-auto px-5 pb-24 no-scrollbar">
        <div className="space-y-6">
          
          {/* Группа */}
          <section className="relative z-30">
            <h3 className="text-xs font-bold text-emerald-700/60 dark:text-emerald-400/60 uppercase tracking-wider mb-2 ml-1">
              Учебная группа
            </h3>
            <div ref={dropdownRef} className="relative">
              <GlassCard className="p-0" onClick={() => setIsDropdownOpen(!isDropdownOpen)}>
                <div className="flex items-center w-full px-4 py-3 cursor-pointer">
                  <div className="mr-3 text-emerald-500/60 dark:text-emerald-400/50 shrink-0">
                    <SearchIcon className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    {isDropdownOpen ? (
                      <input
                        autoFocus
                        type="text"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        placeholder="Поиск..."
                        className="w-full bg-transparent outline-none text-gray-800 dark:text-white font-medium p-0"
                      />
                    ) : (
                      <span className="block w-full text-gray-800 dark:text-white font-medium truncate text-left">
                        {group || 'Выберите группу'}
                      </span>
                    )}
                  </div>
                  <ChevronDownIcon className={`ml-3 w-4 h-4 transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`} />
                </div>
              </GlassCard>

              {isDropdownOpen && (
                <div className="absolute top-full left-0 right-0 mt-2 max-h-56 overflow-y-auto rounded-xl bg-emerald-50/95 dark:bg-[#0a2e23]/95 backdrop-blur-xl border border-emerald-300/25 dark:border-emerald-500/15 shadow-2xl z-50">
                  {filteredGroups.map((g) => (
                    <button
                      key={g}
                      onClick={() => { setGroup(g); setIsDropdownOpen(false); }}
                      className={`w-full text-left px-4 py-3 text-sm ${group === g ? 'bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 font-bold' : 'text-gray-700 dark:text-gray-200'}`}
                    >
                      {g}
                    </button>
                  ))}
                </div>
              )}
            </div>
          </section>

          {/* Уведомления */}
          <section>
            <h3 className="text-xs font-bold text-emerald-700/60 dark:text-emerald-400/60 uppercase tracking-wider mb-2 ml-1">
              Оповещения
            </h3>
            <GlassCard>
              <div className="flex items-center justify-between px-4 py-3">
                <div className="flex flex-col">
                  <span className="font-medium text-gray-800 dark:text-gray-200 text-sm">Изменения в парах</span>
                  <span className="text-[10px] text-gray-500">Уведомлять о заменах и отменах</span>
                </div>
                <button
                  onClick={() => onToggleNotifications(!notificationsEnabled)}
                  className={`relative w-11 h-6 flex items-center rounded-full p-1 transition-colors duration-300 ${notificationsEnabled ? 'bg-emerald-600' : 'bg-gray-400/50'}`}
                >
                  <div className={`bg-white w-4 h-4 rounded-full shadow-md transform transition-transform duration-300 ${notificationsEnabled ? 'translate-x-5' : 'translate-x-0'}`} />
                </button>
              </div>
            </GlassCard>
          </section>

          {/* Тема */}
          <section>
            <h3 className="text-xs font-bold text-emerald-700/60 dark:text-emerald-400/60 uppercase tracking-wider mb-2 ml-1">
              Внешний вид
            </h3>
            <GlassCard>
              <div className="flex items-center justify-between px-4 py-3" onClick={toggleTheme}>
                <div className="flex items-center gap-3">
                  {theme === 'dark' ? <MoonIcon className="w-5 h-5 text-emerald-400" /> : <SunIcon className="w-5 h-5 text-emerald-600" />}
                  <span className="font-medium text-gray-800 dark:text-gray-200 text-sm">Темная тема</span>
                </div>
                <div className={`relative w-11 h-6 flex items-center rounded-full p-1 transition-colors ${theme === 'dark' ? 'bg-emerald-600' : 'bg-emerald-300'}`}>
                   <div className={`bg-white w-4 h-4 rounded-full shadow transform transition-transform ${theme === 'dark' ? 'translate-x-5' : 'translate-x-0'}`} />
                </div>
              </div>
            </GlassCard>
          </section>

        </div>
      </div>
    </div>
  );
};

export default SettingsScreen;