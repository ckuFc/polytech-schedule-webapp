import React, { useState, useEffect, useCallback, useRef } from 'react';
import Layout from './components/Layout';
import BottomMenu from './components/BottomMenu';
import ScheduleScreen from './screens/ScheduleScreen';
import SettingsScreen from './screens/SettingsScreen';
import { Tab, Lesson } from './types';
import { api, ApiError } from './api';

function detectDesktop() {
  const tgPlatform = window.Telegram?.WebApp?.platform;
  if (tgPlatform) {
    return !['android', 'ios'].includes(tgPlatform);
  }
  return !/Android|iPhone|iPad|iPod/i.test(navigator.userAgent);
}

const App = () => {
  const [activeTab, setActiveTab] = useState<Tab>(Tab.SCHEDULE);
  const [currentDate, setCurrentDate] = useState<Date>(new Date());
  const [theme, setTheme] = useState<'light' | 'dark'>('light');

  const [group, setGroup] = useState<string>('');
  const [groups, setGroups] = useState<string[]>([]);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [loading, setLoading] = useState(true);
  const [needSetup, setNeedSetup] = useState(false);
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);

  const dataLoadedRef = useRef(false);

  const fetchSchedule = useCallback(async (grp?: string) => {
    setLoading(true);
    try {
      const res = await api.getSchedule(grp);
      setLessons(res.items || []);

      const returnedGroup = (res.group || '').trim();
      if (returnedGroup.length > 0) {
        setGroup(returnedGroup);
        setNeedSetup(false);
        dataLoadedRef.current = true;
      } else {
        setGroup('');
        setNeedSetup(true);
      }
    } catch (err) {
      const apiErr = err as ApiError;
      if (apiErr.status === 401 || apiErr.need_setup) {
        setNeedSetup(true);
      }
      setLessons([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const loadAllData = useCallback(() => {
    api.getGroups()
      .then(setGroups)
      .catch(() => {});
    
    fetchSchedule();

    api.getMe()
      .then((me) => {
        setNotificationsEnabled(me.notifications_enabled);
      })
      .catch(() => {});
  }, [fetchSchedule]);

  useEffect(() => {
    if (detectDesktop()) {
      document.documentElement.classList.add('desktop-app');
    }
  }, []);

  useEffect(() => {
    const initTelegram = () => {
      const tg = window.Telegram?.WebApp;
      if (!tg) return false;

      tg.ready();
      
      const isMobile = ['android', 'ios'].includes(tg.platform);
      
      if (isMobile) {
        tg.expand();
        try {
          if (tg.requestFullscreen) tg.requestFullscreen();
        } catch (e) {
          console.warn("Fullscreen not supported");
        }
        document.documentElement.classList.remove('desktop-app');
      } else {
        document.documentElement.classList.add('desktop-app');
      }

      const savedTheme = localStorage.getItem('app-theme');
      if (savedTheme === 'light' || savedTheme === 'dark') {
        setTheme(savedTheme);
      } else {
        if (tg.colorScheme === 'dark') {
          setTheme('dark');
        } else if (window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
          setTheme('dark');
        }
      }

      return true;
    };

    if (initTelegram()) {
      loadAllData();
      return;
    }

    let resolved = false;
    
    const interval = setInterval(() => {
      if (window.Telegram?.WebApp) {
        clearInterval(interval);
        initTelegram();
        
        if (!resolved) {
          resolved = true;
          loadAllData();
        } else if (!dataLoadedRef.current) {
          loadAllData();
        }
      }
    }, 100);


    const fallbackTimer = setTimeout(() => {
      if (!resolved) {
        resolved = true;
        const savedTheme = localStorage.getItem('app-theme');
        if (savedTheme === 'light' || savedTheme === 'dark') {
          setTheme(savedTheme);
        } else if (window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
          setTheme('dark');
        }
        loadAllData();
      }
    }, 1500);

    const finalTimer = setTimeout(() => {
      clearInterval(interval);
      if (!dataLoadedRef.current) {
        loadAllData();
      }
    }, 8000);

    return () => {
      clearInterval(interval);
      clearTimeout(fallbackTimer);
      clearTimeout(finalTimer);
    };
  }, [loadAllData]);

  useEffect(() => {
    const now = new Date();
    if (now.getHours() >= 18) {
      const tomorrow = new Date(now);
      tomorrow.setDate(now.getDate() + 1);
      setCurrentDate(tomorrow);
    }
  }, []);

  useEffect(() => {
    const root = window.document.documentElement;
    const tg = window.Telegram?.WebApp;
    const bgColor = theme === 'dark' ? '#022c22' : '#e0eddf';

    if (theme === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }

    if (tg && tg.setHeaderColor && tg.setBackgroundColor) {
      try {
        tg.setHeaderColor(bgColor);
        tg.setBackgroundColor(bgColor);
      } catch (e) {
        console.warn("Color API not supported");
      }
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => {
      const newTheme = prev === 'light' ? 'dark' : 'light';
      localStorage.setItem('app-theme', newTheme);
      return newTheme;
    });
  };

  const handleSetGroup = async (newGroup: string) => {
    setGroup(newGroup);
    setNeedSetup(false);
    try {
      await api.setGroup(newGroup);
    } catch { /* ignore */ }
    fetchSchedule(newGroup);
  };

  const handleToggleNotifications = async (val: boolean) => {
    setNotificationsEnabled(val);
    try {
      await api.toggleNotifications(val);
    } catch {
      setNotificationsEnabled(!val);
    }
  };

  return (
    <Layout>
      <main className="flex-1 w-full relative flex flex-col overflow-hidden">
        {activeTab === Tab.SCHEDULE && (
          <ScheduleScreen
            currentDate={currentDate}
            onDateChange={setCurrentDate}
            lessons={lessons}
            loading={loading}
            needSetup={needSetup}
            onGoToSettings={() => setActiveTab(Tab.SETTINGS)}
          />
        )}
        {activeTab === Tab.SETTINGS && (
          <SettingsScreen
            group={group}
            setGroup={handleSetGroup}
            groups={groups}
            theme={theme}
            toggleTheme={toggleTheme}
            notificationsEnabled={notificationsEnabled}
            onToggleNotifications={handleToggleNotifications}
          />
        )}
      </main>

      <BottomMenu
        activeTab={activeTab}
        onTabChange={setActiveTab}
      />
    </Layout>
  );
};

export default App;