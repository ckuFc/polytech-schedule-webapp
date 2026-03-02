/** @type {import('tailwindcss').Config} */
export default {
    content: [
      "./index.html",
      "./*.{js,ts,jsx,tsx}",              // Для App.tsx и index.tsx
      "./components/**/*.{js,ts,jsx,tsx}",// Для всех компонентов
      "./screens/**/*.{js,ts,jsx,tsx}"    // Для всех экранов
    ],
    darkMode: 'class',
    theme: {
      extend: {
        fontFamily: {
          sans: ['Inter', 'sans-serif'],
        },
        colors: {
          glass: {
            100: 'rgba(255, 255, 255, 0.1)',
            200: 'rgba(255, 255, 255, 0.2)',
            300: 'rgba(255, 255, 255, 0.3)',
            dark: 'rgba(0, 0, 0, 0.3)',
          }
        },
        animation: {
          'blob': 'blob 7s infinite',
          'fade-in': 'fadeIn 0.3s ease-out forwards',
        },
        keyframes: {
          blob: {
            '0%': { transform: 'translate(0px, 0px) scale(1)' },
            '33%': { transform: 'translate(30px, -50px) scale(1.1)' },
            '66%': { transform: 'translate(-20px, 20px) scale(0.9)' },
            '100%': { transform: 'translate(0px, 0px) scale(1)' },
          },
          fadeIn: {
            '0%': { opacity: '0', transform: 'translateY(5px)' },
            '100%': { opacity: '1', transform: 'translateY(0)' },
          }
        }
      },
    },
    plugins: [],
  }