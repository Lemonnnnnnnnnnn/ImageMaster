# ImageMaster Fluent Design System

## 设计理念

本项目采用 Microsoft Fluent Design System 作为核心设计语言，结合黑色主色调，打造现代化、优雅且具有深度感的用户界面。

## 核心原则

### 1. 光线与深度 (Light & Depth)
- 使用阴影和光效营造层次感
- 通过亚克力材料效果增强视觉深度
- 利用Z轴层级管理界面元素

### 2. 运动与动画 (Motion)
- 平滑、有意义的过渡动画
- 统一的缓动函数和时长
- 响应式的交互反馈

### 3. 材质与质感 (Material)
- 毛玻璃效果 (Acrylic)
- 半透明层级
- 柔和的边框和阴影

### 4. 黑色主题 (Dark Theme)
- 深黑色背景营造专业感
- 高对比度文本确保可读性
- 蓝色作为主要强调色

## 色彩系统

### 主色调
```css
--black-primary: #000000     /* 深黑 - 主要背景 */
--black-secondary: #1a1a1a   /* 中黑 - 次要背景 */
--black-tertiary: #2a2a2a    /* 浅黑 - 卡片背景 */
--black-quaternary: #3a3a3a  /* 暗灰 - 边框/分割线 */
--black-quinary: #4a4a4a     /* 中灰 - 禁用状态 */
```

### 文本颜色
```css
--white-primary: #ffffff     /* 主要文本 */
--white-secondary: #c8c8c8   /* 次要文本 */
--white-tertiary: #6e6e6e    /* 禁用/提示文本 */
```

### 强调色
```css
--fluent-blue: #0078d4      /* 主要强调色 - 链接、按钮 */
--fluent-green: #107c10     /* 成功状态 */
--fluent-orange: #ff8c00    /* 警告状态 */
--fluent-red: #d13438       /* 错误状态 */
```

### 半透明背景
```css
--glass-card: rgba(42, 42, 42, 0.8)      /* 卡片背景 */
--glass-overlay: rgba(58, 58, 58, 0.9)   /* 浮层背景 */
--glass-input: rgba(74, 74, 74, 0.6)     /* 输入框背景 */
```

## 布局系统

### 导航栏
- **位置**: 左侧固定
- **宽度**: 
  - 展开状态: 280px (`w-nav-expanded`)
  - 收起状态: 64px (`w-nav-collapsed`)
- **支持折叠**: 通过按钮或键盘快捷键 (Ctrl/Cmd + B)
- **响应式**: 移动端自动收起，显示为覆盖层

### 间距系统
基于 8px 网格系统：
```css
--spacing-xs: 4px
--spacing-sm: 8px
--spacing-md: 16px
--spacing-lg: 24px
--spacing-xl: 32px
--spacing-2xl: 48px
```

### 圆角系统
```css
--radius-sm: 4px      /* 小元素 */
--radius-md: 8px      /* 按钮、输入框 */
--radius-lg: 12px     /* 卡片 */
--radius-xl: 16px     /* 大型容器 */
```

## 组件规范

### 按钮 (Button)
```css
/* 主要按钮 */
.fluent-button-primary {
  background: var(--fluent-blue);
  color: white;
  border-radius: var(--radius-md);
  padding: 8px 16px;
  transition: all 150ms ease-out;
}

/* 次要按钮 */
.fluent-button-secondary {
  background: transparent;
  border: 1px solid var(--white-tertiary);
  color: var(--white-primary);
}
```

### 卡片 (Card)
```css
.fluent-card {
  background: var(--glass-card);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(110, 110, 110, 0.2);
  border-radius: var(--radius-lg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 200ms ease-out;
}

.fluent-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}
```

### 输入框 (Input)
```css
.fluent-input {
  background: var(--glass-input);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(110, 110, 110, 0.3);
  border-radius: var(--radius-md);
  color: var(--white-primary);
  transition: all 200ms ease-out;
}

.fluent-input:focus {
  border-color: var(--fluent-blue);
  box-shadow: 0 0 0 2px rgba(0, 120, 212, 0.2);
  background: rgba(74, 74, 74, 0.8);
}
```

### 导航项 (Navigation Item)
```css
.fluent-nav-item {
  display: flex;
  align-items: center;
  height: 48px;
  padding: 0 16px;
  margin: 0 8px;
  border-radius: var(--radius-md);
  color: var(--white-secondary);
  transition: all 150ms ease-out;
  cursor: pointer;
}

.fluent-nav-item:hover {
  color: var(--white-primary);
  background: rgba(110, 110, 110, 0.1);
}

.fluent-nav-item.active {
  color: var(--white-primary);
  background: rgba(0, 120, 212, 0.2);
  border-left: 2px solid var(--fluent-blue);
}
```

## 动画系统

### 持续时间
```css
--duration-fast: 150ms      /* 按钮悬停 */
--duration-normal: 200ms    /* 普通过渡 */
--duration-medium: 300ms    /* 布局变化 */
--duration-slow: 500ms      /* 页面切换 */
```

### 缓动函数
```css
--ease-out: cubic-bezier(0.25, 0.46, 0.45, 0.94)
--ease-in: cubic-bezier(0.55, 0.06, 0.68, 0.19)
--ease-in-out: cubic-bezier(0.645, 0.045, 0.355, 1)
```

### 常用动画
- **悬停上升**: `transform: translateY(-2px)`
- **按钮按下**: `transform: scale(0.98)`
- **卡片出现**: `opacity + transform`
- **侧边栏收起**: `width + margin-left`

## 阴影系统

```css
--shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.4)
--shadow-md: 0 4px 12px rgba(0, 0, 0, 0.3)
--shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.25)
--shadow-xl: 0 12px 40px rgba(0, 0, 0, 0.2)
```

## 毛玻璃效果

### 背景模糊
```css
--blur-sm: 8px
--blur-md: 16px
--blur-lg: 20px
--blur-xl: 24px
```

### 应用场景
- **卡片背景**: `backdrop-filter: blur(20px)`
- **导航栏**: `backdrop-filter: blur(20px)`
- **模态框**: `backdrop-filter: blur(16px)`
- **输入框**: `backdrop-filter: blur(8px)`

## 响应式设计

### 断点
```css
--breakpoint-sm: 640px
--breakpoint-md: 768px
--breakpoint-lg: 1024px
--breakpoint-xl: 1280px
```

### 移动端适配
- 导航栏自动收起
- 触摸友好的交互尺寸 (最小 44px)
- 覆盖层导航菜单
- 优化的间距和字体大小

## 可访问性

### 键盘导航
- Tab 键顺序逻辑
- Enter/Space 激活元素
- Esc 关闭模态框/菜单
- Ctrl/Cmd + B 切换侧边栏

### 颜色对比
- 文本对比度 ≥ 4.5:1
- 重要元素对比度 ≥ 7:1
- 状态指示不仅依赖颜色

### ARIA 支持
- 适当的 role 属性
- aria-label 描述
- aria-current 状态指示
- 语义化 HTML 结构

## 组件命名规范

### CSS 类命名
- **前缀**: `fluent-` 用于系统组件
- **状态**: `active`, `disabled`, `loading`
- **尺寸**: `sm`, `md`, `lg`, `xl`
- **变体**: `primary`, `secondary`, `ghost`

### 组件文件命名
- **Pascal Case**: `ButtonComponent.svelte`
- **文件夹结构**: `components/common/`, `components/layout/`
- **Store 文件**: `stores/componentStore.ts`

## 使用指南

### 1. 引入样式
```svelte
<div class="fluent-card">
  <button class="fluent-button-primary">
    点击按钮
  </button>
</div>
```

### 2. 自定义变体
```css
.my-custom-card {
  @apply fluent-card;
  /* 添加自定义样式 */
}
```

### 3. 响应式组件
```svelte
<div class="fluent-card lg:w-1/2 md:w-full">
  <!-- 内容 -->
</div>
```

### 4. 动画集成
```svelte
<div class="transition-fluent hover-lift">
  <!-- 自动悬停上升效果 -->
</div>
```

## 最佳实践

1. **保持一致性**: 使用预定义的设计令牌和组件
2. **性能优化**: 合理使用毛玻璃效果，避免过度使用
3. **可访问性优先**: 确保键盘导航和屏幕阅读器支持
4. **移动优先**: 先设计移动端，再扩展到桌面端
5. **语义化**: 使用正确的 HTML 元素和 ARIA 属性

## 更新日志

- **v1.0.0**: 初始版本，建立基础设计系统
- 支持黑色主题
- 左侧导航栏设计
- 基础组件库
- 响应式布局系统 