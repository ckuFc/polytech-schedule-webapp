export interface Lesson {
  id: number;
  external_id: string;
  subject: string;
  teacher: string;
  group: string;
  date: string;
  lesson_num: number;
  room: string;
  sub_group: number;
  zam: number;
}

export interface DaySchedule {
  date: string;
  lessons: Lesson[];
}

export enum Tab {
  SCHEDULE = 'schedule',
  SETTINGS = 'settings',
}