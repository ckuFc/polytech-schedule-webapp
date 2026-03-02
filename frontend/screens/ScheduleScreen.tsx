import React, { useMemo } from 'react';
import DateStrip from '../components/DateStrip';
import LessonCard from '../components/LessonCard';
import { CalendarIcon, SettingsIcon } from '../components/Icons';
import { Lesson } from '../types';
import { LESSON_TIMES } from '../constants';

interface ScheduleScreenProps {
  currentDate: Date;
  onDateChange: (date: Date) => void;
  lessons: Lesson[];
  loading: boolean;
  needSetup: boolean;
  onGoToSettings: () => void;
}

const isCancelledSubject = (subject: string) => {
  const s = subject.trim().toUpperCase();
  return s === 'ОТМЕНА' || s === 'ОТМЕНЕНА' || s.startsWith('ОТМЕН');
};

function dedupByExternalId(lessons: Lesson[]): Lesson[] {
  const byId: Record<string, Lesson[]> = {};
  const result: Lesson[] = [];

  lessons.forEach(l => {
    const id = l.external_id;
    if (!id || id.startsWith('-') || id === '0') {
      result.push(l);
      return;
    }
    if (!byId[id]) byId[id] = [];
    byId[id].push(l);
  });

  for (const id in byId) {
    const group = byId[id];
    if (group.length === 1) {
      result.push(group[0]);
      continue;
    }
    const cancellation = group.find(l => (l.zam || 0) === 1 && isCancelledSubject(l.subject));
    if (cancellation) {
      result.push(cancellation);
      continue;
    }
    const sorted = [...group].sort((a, b) => (a.zam || 0) - (b.zam || 0));
    result.push(sorted[0]);
  }

  return result;
}

function resolveSlot(slotLessons: Lesson[]): Lesson | null {
  if (slotLessons.length === 0) return null;
  if (slotLessons.length === 1) return slotLessons[0];

  const changes = slotLessons.filter(l => (l.zam || 0) > 0);
  const normal = slotLessons.filter(l => (l.zam || 0) === 0);

  if (changes.length > 0) {
    const real = changes.filter(l => !isCancelledSubject(l.subject));
    if (real.length > 0) return real[0];
    return changes[0];
  }

  return normal[0];
}

const ScheduleScreen: React.FC<ScheduleScreenProps> = ({
  currentDate, onDateChange, lessons, loading, needSetup, onGoToSettings,
}) => {
  const dateKey = `${currentDate.getFullYear()}-${String(currentDate.getMonth() + 1).padStart(2, '0')}-${String(currentDate.getDate()).padStart(2, '0')}`;
  
  const now = new Date();
  const todayKey = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
  const isToday = todayKey === dateKey;

  const dayLessons = useMemo(() => {
    return lessons
      .filter(l => l.date.split('T')[0] === dateKey)
      .sort((a, b) => a.lesson_num - b.lesson_num);
  }, [lessons, dateKey]);

  const dedupedLessons = useMemo(() => dedupByExternalId(dayLessons), [dayLessons]);

  const resolvedLessons = useMemo(() => {
    const slots: Record<string, Lesson[]> = {};
    dedupedLessons.forEach(l => {
      const key = `${l.lesson_num}_${l.sub_group}`;
      if (!slots[key]) slots[key] = [];
      slots[key].push(l);
    });

    const result: Lesson[] = [];
    for (const key in slots) {
      const resolved = resolveSlot(slots[key]);
      if (resolved) result.push(resolved);
    }
    return result;
  }, [dedupedLessons]);

  const groupedLessons = useMemo(() => {
    const groups: Record<number, Lesson[]> = {};
    resolvedLessons.forEach(l => {
      if (!groups[l.lesson_num]) groups[l.lesson_num] = [];
      groups[l.lesson_num].push(l);
    });
    return groups;
  }, [resolvedLessons]);

  const sortedLessonNums = useMemo(() => {
    return Object.keys(groupedLessons).map(Number).sort((a, b) => a - b);
  }, [groupedLessons]);

  const nextLessonNum = useMemo(() => {
    if (!isToday || sortedLessonNums.length === 0) return -1;
    const nowMinutes = now.getHours() * 60 + now.getMinutes();
    for (const num of sortedLessonNums) {
      const times = LESSON_TIMES[num];
      if (!times) continue;
      const [eh, em] = times.end.split(':').map(Number);
      if (nowMinutes < eh * 60 + em) return num;
    }
    return -1;
  }, [isToday, sortedLessonNums, now]);

  const dayName = currentDate.toLocaleDateString('ru-RU', { weekday: 'long' });
  const capitalizedDayName = dayName.charAt(0).toUpperCase() + dayName.slice(1);
  const totalPairs = sortedLessonNums.length;

  const getPairsLabel = (count: number) => {
    if (count % 10 === 1 && count % 100 !== 11) return 'пара';
    if ([2, 3, 4].includes(count % 10) && ![12, 13, 14].includes(count % 100)) return 'пары';
    return 'пар';
  };

  return (
    <div className="flex flex-col h-full w-full overflow-hidden">
      {/* ТУТ ИЗМЕНЕН pt-24 НА safe-pt */}
      <div className="px-6 safe-pt pb-2 shrink-0 flex justify-between items-baseline">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white tracking-tight">
          Расписание
        </h1>
        <span className="text-sm font-medium text-emerald-600 dark:text-emerald-400 opacity-90">
          {currentDate.getFullYear()}
        </span>
      </div>

      <DateStrip selectedDate={currentDate} onSelectDate={onDateChange} />

      <div className="flex-1 overflow-y-auto px-5 pb-20 no-scrollbar">
        {needSetup && !loading && (
          <div className="flex flex-col items-center justify-center h-64 text-center opacity-70">
            <div className="p-5 bg-emerald-50 dark:bg-emerald-900/10 rounded-full mb-4 shadow-sm ring-1 ring-emerald-100 dark:ring-emerald-500/10">
              <SettingsIcon className="w-10 h-10 text-emerald-400 dark:text-emerald-600" />
            </div>
            <p className="text-lg font-medium text-gray-800 dark:text-gray-200">Выберите группу</p>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Перейдите в настройки и укажите свою группу</p>
            <button
              onClick={onGoToSettings}
              className="mt-4 px-5 py-2 rounded-xl bg-emerald-700 text-white text-sm font-semibold shadow-lg shadow-emerald-900/20"
            >
              Настройки
            </button>
          </div>
        )}

        {loading && (
          <div className="flex items-center justify-center h-64">
            <div className="w-8 h-8 border-3 border-emerald-200 border-t-emerald-600 rounded-full animate-spin"></div>
          </div>
        )}

        {!loading && !needSetup && (
          <>
            <div className="mb-3 mt-1 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-200">
                {capitalizedDayName}
              </h2>
              {totalPairs > 0 && (
                <span className="text-[11px] text-emerald-800 dark:text-emerald-300 font-bold bg-emerald-100/50 dark:bg-emerald-900/30 px-2.5 py-1 rounded-lg border border-emerald-200/50 dark:border-emerald-500/20">
                  {totalPairs} {getPairsLabel(totalPairs)}
                </span>
              )}
            </div>

            {totalPairs > 0 ? (
              <div className="space-y-3.5">
                {sortedLessonNums.map((num) => (
                  <LessonCard
                    key={`${dateKey}-${num}`}
                    lessons={groupedLessons[num]}
                    isNext={num === nextLessonNum}
                  />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-64 text-center opacity-70">
                <div className="p-5 bg-emerald-50 dark:bg-emerald-900/10 rounded-full mb-4 shadow-sm ring-1 ring-emerald-100 dark:ring-emerald-500/10">
                  <CalendarIcon className="w-10 h-10 text-emerald-400 dark:text-emerald-600" />
                </div>
                <p className="text-lg font-medium text-gray-800 dark:text-gray-200">Пар нет</p>
                <p className="text-sm text-gray-500">Можно отдыхать!</p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default ScheduleScreen;