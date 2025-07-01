/**
 * 设计令牌 - 统一管理项目的设计系统
 * 包含颜色、间距、字体、阴影等设计元素
 */

// 颜色系统
export const colors = {
  // 主色调 - 蓝色系
  primary: {
    50: '#eff6ff',
    100: '#dbeafe',
    200: '#bfdbfe',
    300: '#93c5fd',
    400: '#60a5fa',
    500: '#3b82f6', // 主色
    600: '#2563eb',
    700: '#1d4ed8',
    800: '#1e40af',
    900: '#1e3a8a'
  },
  
  // 次要色调 - 绿色系
  secondary: {
    50: '#f0fdf4',
    100: '#dcfce7',
    200: '#bbf7d0',
    300: '#86efac',
    400: '#4ade80',
    500: '#22c55e', // 次要色
    600: '#16a34a',
    700: '#15803d',
    800: '#166534',
    900: '#14532d'
  },
  
  // 成功色
  success: {
    50: '#f0fdf4',
    100: '#dcfce7',
    500: '#22c55e',
    600: '#16a34a',
    700: '#15803d'
  },
  
  // 警告色
  warning: {
    50: '#fffbeb',
    100: '#fef3c7',
    500: '#f59e0b',
    600: '#d97706',
    700: '#b45309'
  },
  
  // 错误色
  error: {
    50: '#fef2f2',
    100: '#fee2e2',
    500: '#ef4444',
    600: '#dc2626',
    700: '#b91c1c'
  },
  
  // 中性色
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
} as const;

// 间距系统
export const spacing = {
  xs: '0.25rem',   // 4px
  sm: '0.5rem',    // 8px
  md: '1rem',      // 16px
  lg: '1.5rem',    // 24px
  xl: '2rem',      // 32px
  '2xl': '3rem',   // 48px
  '3xl': '4rem',   // 64px
} as const;

// 字体系统
export const typography = {
  fontFamily: {
    sans: ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Open Sans', 'Helvetica Neue', 'sans-serif'],
    mono: ['SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', 'monospace']
  },
  fontSize: {
    xs: '0.75rem',   // 12px
    sm: '0.875rem',  // 14px
    base: '1rem',    // 16px
    lg: '1.125rem',  // 18px
    xl: '1.25rem',   // 20px
    '2xl': '1.5rem', // 24px
    '3xl': '1.875rem', // 30px
    '4xl': '2.25rem'   // 36px
  },
  fontWeight: {
    normal: '400',
    medium: '500',
    semibold: '600',
    bold: '700'
  }
} as const;

// 阴影系统
export const shadows = {
  sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
  md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
  lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
  xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
} as const;

// 圆角系统
export const borderRadius = {
  sm: '0.375rem',  // 6px
  md: '0.5rem',    // 8px
  lg: '0.75rem',   // 12px
  xl: '1rem'       // 16px
} as const;

// 过渡动画
export const transitions = {
  fast: '150ms ease-in-out',
  normal: '200ms ease-in-out',
  slow: '300ms ease-in-out'
} as const;

// 组件样式预设
export const componentStyles = {
  // 按钮样式
  button: {
    base: 'inline-flex items-center justify-center font-medium rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2',
    sizes: {
      sm: 'px-3 py-1.5 text-sm',
      md: 'px-4 py-2 text-sm',
      lg: 'px-6 py-3 text-base'
    }
  },
  
  // 卡片样式
  card: {
    base: 'bg-white border border-neutral-200 rounded-lg shadow-sm',
    hover: 'hover:shadow-md transition-shadow duration-200'
  },
  
  // 表单元素样式
  input: {
    base: 'w-full px-3 py-2 border border-neutral-300 rounded-md text-sm transition-colors focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-100',
    error: 'border-error-500 focus:border-error-500 focus:ring-error-100'
  },
  
  // 状态样式
  status: {
    success: 'text-success-600 bg-success-50 border-success-200',
    warning: 'text-warning-600 bg-warning-50 border-warning-200',
    error: 'text-error-600 bg-error-50 border-error-200',
    info: 'text-primary-600 bg-primary-50 border-primary-200'
  }
} as const;

// 导出所有设计令牌
export const designTokens = {
  colors,
  spacing,
  typography,
  shadows,
  borderRadius,
  transitions,
  componentStyles
} as const;

export default designTokens;