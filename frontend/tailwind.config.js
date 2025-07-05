import { defineConfig } from '@tailwindcss/vite'

export default defineConfig({
  theme: {
    extend: {
      colors: {
        // 黑色系主色调
        black: {
          primary: '#000000',    // 深黑
          secondary: '#1a1a1a',  // 中黑  
          tertiary: '#2a2a2a',   // 浅黑
          quaternary: '#3a3a3a', // 暗灰
          quinary: '#4a4a4a',    // 中灰
        },
        // 白色文本系统
        white: {
          primary: '#ffffff',    // 主要文本
          secondary: '#c8c8c8',  // 次要文本
          tertiary: '#6e6e6e',   // 禁用文本
        },
        // Fluent 强调色系统
        fluent: {
          blue: '#0078d4',      // 主要强调色
          green: '#107c10',     // 成功色
          orange: '#ff8c00',    // 警告色
          red: '#d13438',       // 错误色
        },
        // 半透明背景色
        glass: {
          card: 'rgba(42, 42, 42, 0.8)',      // 卡片背景
          overlay: 'rgba(58, 58, 58, 0.9)',   // 浮层背景
          input: 'rgba(74, 74, 74, 0.6)',     // 输入框背景
        },
        // 保持向后兼容的原有颜色
        primary: {
          50: '#eff6ff',
          100: '#dbeafe', 
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a'
        },
        secondary: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0', 
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d'
        },
        success: {
          50: '#f0fdf4',
          100: '#dcfce7',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d'
        },
        warning: {
          50: '#fffbeb',
          100: '#fef3c7',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309'
        },
        error: {
          50: '#fef2f2',
          100: '#fee2e2',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c'
        },
        neutral: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827'
        }
      },
      // Fluent Design 阴影系统
      boxShadow: {
        'fluent-sm': '0 1px 3px rgba(0, 0, 0, 0.4)',
        'fluent-md': '0 4px 12px rgba(0, 0, 0, 0.3)',
        'fluent-lg': '0 8px 24px rgba(0, 0, 0, 0.25)',
        'fluent-xl': '0 12px 40px rgba(0, 0, 0, 0.2)',
      },
      // Fluent Design 圆角
      borderRadius: {
        'fluent-sm': '4px',
        'fluent-md': '8px',
        'fluent-lg': '12px',
        'fluent-xl': '16px',
      },
      // Fluent Design 动画
      transitionDuration: {
        'fluent-fast': '150ms',
        'fluent-normal': '200ms',
        'fluent-medium': '300ms',
        'fluent-slow': '500ms',
      },
      transitionTimingFunction: {
        'fluent-ease-out': 'cubic-bezier(0.25, 0.46, 0.45, 0.94)',
        'fluent-ease-in': 'cubic-bezier(0.55, 0.06, 0.68, 0.19)',
        'fluent-ease-in-out': 'cubic-bezier(0.645, 0.045, 0.355, 1)',
      },
      // 毛玻璃效果
      backdropBlur: {
        'fluent': '20px',
        'fluent-sm': '8px',
        'fluent-md': '16px',
        'fluent-lg': '24px',
      },
      // 导航栏宽度
      width: {
        'nav-collapsed': '64px',
        'nav-expanded': '280px',
      },
      // 间距系统 (8px 基数)
      spacing: {
        '18': '4.5rem', // 72px
        '22': '5.5rem', // 88px
      }
    }
  },
  // 支持暗色主题
  darkMode: 'class'
})