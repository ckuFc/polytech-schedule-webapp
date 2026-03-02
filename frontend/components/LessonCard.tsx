import React, { useRef, useState, useEffect } from 'react';
import { Lesson } from '../types';
import { LESSON_TIMES, getSubjectStyles } from '../constants';
import GlassCard from './GlassCard';
import { MapPinIcon } from './Icons';

interface LessonCardProps {
  lessons: Lesson[];
  isNext?: boolean;
}

const isCancelled = (lesson: Lesson) => {
  const s = lesson.subject.trim().toUpperCase();
  return s === 'ОТМЕНА' || s === 'ОТМЕНЕНА' || s.startsWith('ОТМЕН');
};

const isReplacement = (lesson: Lesson) => {
  return (lesson.zam || 0) >= 2 && !isCancelled(lesson);
};

const LessonContent = ({ lesson, isSplit }: { lesson: Lesson; isSplit: boolean }) => {
  const styles = getSubjectStyles(lesson.subject);
  const isSubgroup = lesson.sub_group > 0;
  const cancelled = isCancelled(lesson);
  const replacement = isReplacement(lesson);

  const subjectRef = useRef<HTMLHeadingElement>(null);
  const [fontClass, setFontClass] = useState('text-[16px]');

  // Адаптивный размер шрифта: начинаем с 16px, уменьшаем если не влезает
  useEffect(() => {
    if (!isSplit) {
      setFontClass('text-[16px]');
      return;
    }

    const el = subjectRef.current;
    if (!el) return;

    // Начинаем с нормального размера
    setFontClass('text-[16px]');

    // Проверяем после рендера
    requestAnimationFrame(() => {
      if (!el) return;
      const isOverflowing = el.scrollHeight > el.clientHeight + 2;
      if (isOverflowing) {
        setFontClass('text-[13px]');
      }
    });
  }, [lesson.subject, isSplit]);

  const getSubgroupStyle = (groupNum: number) => {
    if (groupNum === 1) return 'bg-indigo-100/50 text-indigo-500 border-indigo-200/40 dark:bg-indigo-500/20 dark:text-indigo-300 dark:border-indigo-500/20';
    if (groupNum === 2) return 'bg-orange-100/50 text-orange-500 border-orange-200/40 dark:bg-orange-500/20 dark:text-orange-300 dark:border-orange-500/20';
    return 'bg-gray-100/50 text-gray-500 border-gray-200/40 dark:bg-white/10 dark:text-gray-300 dark:border-white/5';
  };

  return (
    <div className={`flex flex-col h-full min-w-0 relative justify-between ${cancelled ? 'opacity-45' : ''}`}>
      {/* Бейджики */}
      <div className="flex items-start gap-1.5 mb-1.5 flex-wrap">
        {cancelled ? (
          <>
            <span className="px-2.5 py-[3px] rounded-md text-[10px] font-bold uppercase border tracking-wide bg-red-100/60 text-red-500 border-red-200/40 dark:bg-red-500/15 dark:text-red-400 dark:border-red-500/20">
              Отмена
            </span>
            {isSubgroup && (
              <span className={`px-2 py-[3px] rounded-md text-[9px] font-bold uppercase tracking-wide border ${getSubgroupStyle(lesson.sub_group)}`}>
                {lesson.sub_group} подгр.
              </span>
            )}
          </>
        ) : isSubgroup ? (
           <span className={`
             px-2.5 py-[3px] rounded-md text-[10px] font-bold uppercase tracking-wide border
             ${getSubgroupStyle(lesson.sub_group)}
           `}>
             {lesson.sub_group} подгр.
           </span>
        ) : (
           <span className={`px-2.5 py-[3px] rounded-md text-[10px] font-bold uppercase border tracking-wide ${styles.badge}`}>
             {lesson.lesson_num} пара
           </span>
        )}
        {replacement && (
          <span className="px-2 py-[3px] rounded-md text-[9px] font-bold uppercase border tracking-wide bg-amber-100/50 text-amber-600 border-amber-200/40 dark:bg-amber-500/15 dark:text-amber-300 dark:border-amber-500/20">
            Замена
          </span>
        )}
      </div>

      {/* Предмет — адаптивный размер шрифта */}
      <h3
        ref={subjectRef}
        className={`font-semibold leading-snug break-words my-auto
        ${cancelled
          ? 'line-through text-gray-400 dark:text-gray-500'
          : 'text-gray-700 dark:text-gray-100'}
        ${isSplit ? fontClass : 'text-[17px] pr-4'}`}
        style={isSplit ? { maxHeight: '3.6em', overflow: 'hidden' } : undefined}
      >
        {lesson.subject}
      </h3>

      {/* Преподаватель и Аудитория */}
      <div className={`mt-auto flex ${isSplit ? 'flex-col gap-1.5 items-start' : 'flex-row items-center justify-between'}`}>
        <span className={`font-medium truncate max-w-full ${cancelled ? 'text-gray-400 dark:text-gray-600' : 'text-gray-500 dark:text-gray-400'} ${isSplit ? 'text-[12px]' : 'text-[13px]'}`}>
          {lesson.teacher}
        </span>

        {lesson.room && !cancelled && (
          <div className={`
              flex items-center gap-1
              bg-emerald-50/40 dark:bg-emerald-900/20 backdrop-blur-md
              px-2 py-0.5 rounded-md
              border border-emerald-200/30 dark:border-emerald-500/10
              shrink-0
              ${isSplit ? '' : 'ml-2'}
          `}>
            <MapPinIcon className={`w-3 h-3 ${styles.room} opacity-80`} />
            <span className={`text-[11px] font-bold ${styles.room} dark:text-emerald-200`}>
              {lesson.room}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

const LessonCard: React.FC<LessonCardProps> = ({ lessons, isNext = false }) => {
  if (!lessons || lessons.length === 0) return null;

  const sortedLessons = [...lessons].sort((a, b) => a.sub_group - b.sub_group);
  const mainLesson = sortedLessons[0];
  const isSplit = sortedLessons.length > 1;
  const allCancelled = sortedLessons.every(l => isCancelled(l));

  const times = LESSON_TIMES[mainLesson.lesson_num];
  const timeStart = times?.start || `${mainLesson.lesson_num}`;
  const timeEnd = times?.end || '';

  const timeColWidth = 'w-[68px]';

  return (
    <GlassCard className={`
        mb-0 overflow-visible transition-transform duration-300 min-h-[110px]
        ${allCancelled ? 'opacity-55' : ''}
        ${isNext
        ? 'ring-2 ring-emerald-400/30 dark:ring-emerald-400/25 shadow-lg shadow-emerald-500/5'
        : 'hover:border-emerald-300/40 dark:hover:border-emerald-500/15'}
    `}>
      {isNext && !allCancelled && (
        <div className="absolute -top-px -right-px px-2.5 py-0.5 bg-emerald-600 text-white text-[9px] font-bold uppercase tracking-wider rounded-bl-xl rounded-tr-xl z-20 shadow-sm">
          Сейчас
        </div>
      )}

      <div className="flex flex-row items-stretch rounded-2xl overflow-hidden min-h-[110px]">
        
        {/* КОЛОНКА ВРЕМЕНИ */}
        <div className={`
            flex flex-col justify-center items-center text-center
            ${timeColWidth} shrink-0 
            border-r border-emerald-200/25 dark:border-emerald-500/10 
            bg-emerald-50/40 dark:bg-emerald-950/30 backdrop-blur-sm
            py-3
        `}>
          <div className={`text-[16px] font-extrabold leading-none tracking-tight ${allCancelled ? 'text-gray-400 dark:text-gray-500' : 'text-gray-800 dark:text-white'}`}>
            {timeStart}
          </div>
          {timeEnd && (
            <>
              <div className={`h-8 w-[2px] rounded-full my-2 ${allCancelled
                ? 'bg-gradient-to-b from-gray-200/20 via-gray-300/20 to-gray-200/20 dark:from-gray-700/20 dark:via-gray-600/20 dark:to-gray-700/20'
                : 'bg-gradient-to-b from-emerald-200/20 via-emerald-400/40 to-emerald-200/20 dark:from-emerald-800/20 dark:via-emerald-500/30 dark:to-emerald-800/20'
              }`}></div>
              <div className={`text-[13px] font-semibold leading-none ${allCancelled ? 'text-gray-400 dark:text-gray-600' : 'text-emerald-500/70 dark:text-emerald-400'}`}>
                {timeEnd}
              </div>
            </>
          )}
        </div>

        {/* ПРАВАЯ ЧАСТЬ */}
        <div className="flex-1 min-w-0 bg-transparent">
          {isSplit ? (
            <div className="flex flex-row h-full divide-x divide-emerald-200/25 dark:divide-emerald-500/10">
              {sortedLessons.map((lesson, idx) => (
                <div
                  key={`${lesson.sub_group}-${idx}`}
                  className={`flex-1 min-w-0 p-3 flex flex-col ${idx > 0 ? 'bg-emerald-50/15 dark:bg-emerald-900/10' : ''}`}
                >
                  <LessonContent lesson={lesson} isSplit={true} />
                </div>
              ))}
            </div>
          ) : (
            <div className="h-full p-4 flex flex-col">
              <LessonContent lesson={mainLesson} isSplit={false} />
            </div>
          )}
        </div>
      </div>
    </GlassCard>
  );
};

export default LessonCard;