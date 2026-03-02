declare global {
  interface Window {
    Telegram?: {
      WebApp?: {
        initData: string;
        platform: string;
        colorScheme: string;
        setHeaderColor: (color: string) => void;
        setBackgroundColor: (color: string) => void;
        requestFullscreen?: () => void;
        themeParams: Record<string, string>;
        ready: () => void;
        expand: () => void;
        close: () => void;
      };
    };
  }
}

function getInitData(): string {
  return window.Telegram?.WebApp?.initData || '';
}

export interface ApiError {
  status: number;
  error?: string;
  need_setup?: boolean;
}

export interface UserInfo {
  group: string;
  notifications_enabled: boolean;
}

const API_BASE = '/api/v1';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const initData = getInitData();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (initData) {
    headers['Authorization'] = `tma ${initData}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      ...headers,
      ...(options?.headers as Record<string, string>),
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    const err: ApiError = { status: res.status, ...body };
    throw err;
  }

  const json = await res.json();
  return json.data;
}

export interface ScheduleResponse {
  group: string;
  items: Lesson[];
}

import { Lesson } from './types';

export const api = {
  getSchedule: (group?: string) =>
    request<{ group: string; items: Lesson[] }>(
      `/schedule${group ? `?group=${encodeURIComponent(group)}` : ''}`
    ),

  getGroups: () => request<string[]>('/groups'),

  getMe: () => request<UserInfo>('/user/me'),

  setGroup: (group: string) =>
    request<{ status: string }>('/schedule/group', {
      method: 'POST',
      body: JSON.stringify({ group }),
    }),

  toggleNotifications: (enabled: boolean) =>
    request<{ status: string }>('/user/notifications', {
      method: 'POST',
      body: JSON.stringify({ notifications_enabled: enabled }),
    }),
};