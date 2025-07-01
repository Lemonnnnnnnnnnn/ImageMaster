# ImageMaster 色值系统使用指南

## 概述

本项目采用统一的色值系统，确保所有UI组件的视觉效果保持一致。色值系统基于Tailwind CSS扩展配置，并提供CSS变量和TypeScript设计令牌支持。

## 色值体系

### 主色调 (Primary)
- **用途**: 主要操作按钮、链接、重要信息强调
- **色值**: 蓝色系 (#3b82f6)
- **类名**: `primary-50` 到 `primary-900`
- **CSS变量**: `--color-primary`

### 次要色调 (Secondary)
- **用途**: 次要操作、辅助信息
- **色值**: 绿色系 (#22c55e)
- **类名**: `secondary-50` 到 `secondary-900`
- **CSS变量**: `--color-secondary`

### 语义化颜色

#### 成功色 (Success)
- **用途**: 成功状态、确认操作
- **色值**: 绿色 (#22c55e)
- **类名**: `success-50`, `success-100`, `success-500`, `success-600`, `success-700`

#### 警告色 (Warning)
- **用途**: 警告信息、需要注意的操作
- **色值**: 黄色 (#f59e0b)
- **类名**: `warning-50`, `warning-100`, `warning-500`, `warning-600`, `warning-700`

#### 错误色 (Error)
- **用途**: 错误状态、危险操作
- **色值**: 红色 (#ef4444)
- **类名**: `error-50`, `error-100`, `error-500`, `error-600`, `error-700`

#### 中性色 (Neutral)
- **用途**: 文本、边框、背景
- **色值**: 灰色系
- **类名**: `neutral-50` 到 `neutral-900`

## 使用方式

### 1. Tailwind CSS 类名

```html
<!-- 主色按钮 -->
<button class="bg-primary-600 text-white hover:bg-primary-700">
  主要操作
</button>

<!-- 成功状态 -->
<div class="bg-success-50 text-success-600 border border-success-200">
  操作成功
</div>

<!-- 中性色文本 -->
<p class="text-neutral-600">这是一段说明文字</p>
```

### 2. CSS 变量

```css
.custom-element {
  background-color: var(--color-primary);
  border-color: var(--color-border);
  color: var(--color-text-primary);
}
```

### 3. TypeScript 设计令牌

```typescript
import { colors, componentStyles } from '../utils/designTokens';

// 在组件中使用
const buttonStyle = {
  backgroundColor: colors.primary[600],
  color: 'white'
};

// 使用预设样式
const cardClasses = componentStyles.card.base;
```

## 组件样式规范

### 按钮组件

```typescript
// 推荐的按钮颜色映射
const buttonColors = {
  primary: 'bg-primary-600 hover:bg-primary-700',
  secondary: 'bg-secondary-600 hover:bg-secondary-700',
  success: 'bg-success-600 hover:bg-success-700',
  warning: 'bg-warning-600 hover:bg-warning-700',
  error: 'bg-error-600 hover:bg-error-700',
  neutral: 'bg-neutral-500 hover:bg-neutral-600'
};
```

### 卡片组件

```html
<!-- 标准卡片样式 -->
<div class="bg-white border border-neutral-200 rounded-lg shadow-sm">
  <!-- 卡片内容 -->
</div>
```

### 表单元素

```html
<!-- 输入框 -->
<input class="border-neutral-300 focus:border-primary-500 focus:ring-primary-100" />

<!-- 错误状态 -->
<input class="border-error-500 focus:border-error-500 focus:ring-error-100" />
```

### 状态指示器

```html
<!-- 成功状态 -->
<span class="bg-success-50 text-success-600 border border-success-200 px-2 py-1 rounded">
  已完成
</span>

<!-- 警告状态 -->
<span class="bg-warning-50 text-warning-600 border border-warning-200 px-2 py-1 rounded">
  待处理
</span>

<!-- 错误状态 -->
<span class="bg-error-50 text-error-600 border border-error-200 px-2 py-1 rounded">
  失败
</span>
```

## 最佳实践

### 1. 优先使用语义化颜色
- ✅ 使用 `success-600` 表示成功状态
- ❌ 使用 `green-600` 表示成功状态

### 2. 保持颜色层次一致
- 主要内容使用 `600` 色阶
- 悬停状态使用 `700` 色阶
- 背景使用 `50` 色阶
- 边框使用 `200` 色阶

### 3. 避免硬编码颜色值
- ✅ 使用 `bg-primary-600`
- ❌ 使用 `bg-[#2563eb]`

### 4. 合理使用对比度
- 确保文本与背景有足够的对比度
- 深色背景配浅色文字
- 浅色背景配深色文字

## 扩展指南

### 添加新颜色

1. 在 `tailwind.config.js` 中添加新的颜色定义
2. 在 `style.css` 中添加对应的CSS变量
3. 在 `designTokens.ts` 中更新颜色对象
4. 更新相关组件的颜色映射

### 主题切换支持

色值系统已预留深色主题支持，可通过CSS变量实现主题切换：

```css
[data-theme="dark"] {
  --color-background: #1f2937;
  --color-text-primary: #f9fafb;
  /* 其他深色主题变量 */
}
```

## 工具支持

### VS Code 扩展推荐
- Tailwind CSS IntelliSense
- CSS Variable Autocomplete

### 开发工具
- 使用 `twMerge` 合并Tailwind类名
- 使用 ESLint 规则检查硬编码颜色

## 常见问题

### Q: 如何选择合适的色阶？
A: 遵循以下规则：
- 50-100: 背景色
- 200-300: 边框色
- 400-500: 图标色
- 600-700: 文本色和按钮色
- 800-900: 强调色

### Q: 什么时候使用CSS变量而不是Tailwind类名？
A: 在以下情况使用CSS变量：
- 动态计算颜色值
- 需要在JavaScript中访问颜色值
- 实现主题切换功能

### Q: 如何确保颜色的可访问性？
A: 
- 使用在线对比度检查工具
- 确保文本对比度至少达到 4.5:1
- 不仅依赖颜色传达信息，还要配合图标或文字