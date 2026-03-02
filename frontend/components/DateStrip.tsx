import React, { useRef } from 'react';
import { CalendarPlusIcon } from './Icons';

interface DateStripProps {
  selectedDate: Date;
  onSelectDate: (date: Date) => void;
}

const DateStrip: React.FC<DateStripProps> = ({ selectedDate, onSelectDate }) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  
  const days = [];
  const today = new Date();
  
  for (let i = -1; i <= 3; i++) {
    const d = new Date();
    d.setDate(today.getDate() + i);
    days.push(d);
  }

  const isSameDay = (d1: Date, d2: Date) => 
    d1.getDate() === d2.getDate() && 
    d1.getMonth() === d2.getMonth() &&
    d1.getFullYear() === d2.getFullYear();

  const isSelectedInRange = days.some(d => isSameDay(d, selectedDate));

  const handleManualDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.value) {
      const newDate = new Date(e.target.value);
      onSelectDate(newDate);
      e.target.value = ''; 
    }
  };

  return (
    <div className="w-full z-20 shrink-0 mb-1">
      <div className="px-1">
        <div 
          ref={scrollRef}
          className="flex overflow-x-auto gap-3 px-4 py-6 -my-4 no-scrollbar snap-x items-center"
          style={{ scrollPaddingLeft: '1rem' }}
        >
          {days.map((date, index) => {
            const isSelected = isSameDay(date, selectedDate);
            const isToday = isSameDay(date, new Date());
            
            return (
              <button
                key={index}
                onClick={() => onSelectDate(date)}
                className={`
                  flex flex-col items-center justify-between
                  w-[64px] h-[76px] shrink-0 py-2.5
                  rounded-2xl snap-center transition-all duration-300 relative
                  ${isSelected 
                    ? 'bg-emerald-700 text-white shadow-lg shadow-emerald-900/30 scale-105 z-10 ring-1 ring-emerald-400/20' 
                    : isToday 
                      ? 'bg-emerald-200/50 dark:bg-emerald-900/20 border-2 border-emerald-400/40 text-emerald-800 dark:text-emerald-400' 
                      : 'bg-emerald-100/40 dark:bg-white/5 border border-emerald-300/25 dark:border-white/5 text-gray-500 dark:text-gray-400 hover:bg-emerald-100/60 dark:hover:bg-white/10'
                  }
                `}
              >
                <div className={`
                  text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full
                  ${isSelected ? 'bg-emerald-600/50 text-white/90' : 'opacity-70'}
                `}>
                  {isToday ? 'СЕГ' : date.toLocaleDateString('ru-RU', { weekday: 'short' })}
                </div>
                
                <span className={`text-2xl font-bold leading-none tracking-tight mb-1
                  ${isSelected ? 'text-white' : 'text-gray-800 dark:text-gray-200'}
                `}>
                  {date.getDate()}
                </span>
                
                {isToday && !isSelected && (
                   <div className="absolute bottom-1.5 w-6 h-1 rounded-full bg-emerald-500/50"></div>
                )}
              </button>
            );
          })}

          <div className="w-px h-10 bg-emerald-300/40 dark:bg-white/10 shrink-0 mx-1"></div>

          <div className="relative shrink-0 group">
              <div
                  className={`
                      relative flex items-center justify-center 
                      w-[64px] h-[76px] rounded-2xl
                      transition-all duration-300
                      ${!isSelectedInRange 
                          ? 'bg-emerald-700 text-white shadow-lg shadow-emerald-900/30' 
                          : 'bg-emerald-100/40 dark:bg-white/5 text-gray-500 dark:text-gray-400 border border-emerald-300/25 dark:border-white/5 hover:bg-emerald-100/60'
                      }
                  `}
              >
                  <CalendarPlusIcon className={`w-6 h-6 pointer-events-none ${!isSelectedInRange ? 'text-white' : ''}`} />
                  
                  <input 
                      type="date" 
                      className="absolute inset-0 w-full h-full opacity-0 z-10 cursor-pointer date-input-full-trigger"
                      onChange={handleManualDateChange}
                  />
              </div>
          </div>
          
          <div className="w-2 shrink-0"></div>
        </div>
      </div>
    </div>
  );
};

export default DateStrip;