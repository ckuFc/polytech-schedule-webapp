export const LESSON_TIMES: Record<number, { start: string; end: string }> = {
  1: { start: '09:00', end: '10:30' },
  2: { start: '10:45', end: '12:15' },
  3: { start: '13:05', end: '14:35' },
  4: { start: '14:50', end: '16:20' },
  5: { start: '16:30', end: '18:00' },
};

const SUBJECT_COLORS = [
  {
    badge: 'bg-blue-100/60 text-blue-600 border-blue-200/40 dark:bg-blue-500/15 dark:text-blue-300 dark:border-blue-500/20',
    room: 'text-blue-500 dark:text-blue-300',
  },
  {
    badge: 'bg-emerald-100/60 text-emerald-600 border-emerald-200/40 dark:bg-emerald-500/15 dark:text-emerald-300 dark:border-emerald-500/20',
    room: 'text-emerald-500 dark:text-emerald-300',
  },
  {
    badge: 'bg-purple-100/60 text-purple-600 border-purple-200/40 dark:bg-purple-500/15 dark:text-purple-300 dark:border-purple-500/20',
    room: 'text-purple-500 dark:text-purple-300',
  },
  {
    badge: 'bg-rose-100/60 text-rose-600 border-rose-200/40 dark:bg-rose-500/15 dark:text-rose-300 dark:border-rose-500/20',
    room: 'text-rose-500 dark:text-rose-300',
  },
  {
    badge: 'bg-amber-100/60 text-amber-600 border-amber-200/40 dark:bg-amber-500/15 dark:text-amber-300 dark:border-amber-500/20',
    room: 'text-amber-500 dark:text-amber-300',
  },
  {
    badge: 'bg-cyan-100/60 text-cyan-600 border-cyan-200/40 dark:bg-cyan-500/15 dark:text-cyan-300 dark:border-cyan-500/20',
    room: 'text-cyan-500 dark:text-cyan-300',
  },
];

function hashString(str: string): number {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) - hash) + str.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

export function getSubjectStyles(subject: string) {
  const idx = hashString(subject) % SUBJECT_COLORS.length;
  return SUBJECT_COLORS[idx];
}