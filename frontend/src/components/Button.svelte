<script lang="ts">
  import { twMerge } from 'tailwind-merge';
  import { componentStyles } from '../utils/designTokens';
  
  // Props
  export let variant: 'filled' | 'outlined' | 'ghost' = 'filled';
  export let color: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'gray' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let disabled: boolean = false;
  export let loading: boolean = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let href: string | undefined = undefined;
  export let target: string | undefined = undefined;
  export let rel: string | undefined = undefined;
  export let download: string | boolean | undefined = undefined;
  
  // Additional classes
  export let classes: string = '';
  
  // 使用设计令牌中的按钮基础样式
  $: baseClasses = componentStyles.button.base;
  
  // 构建变体和颜色类名
  $: variantClasses = (() => {
    const colorMap = {
      primary: {
        filled: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500',
        outlined: 'border border-primary-600 text-primary-600 hover:bg-primary-50 focus:ring-primary-500',
        ghost: 'text-primary-600 hover:bg-primary-50 focus:ring-primary-500'
      },
      secondary: {
        filled: 'bg-secondary-600 text-white hover:bg-secondary-700 focus:ring-secondary-500',
        outlined: 'border border-secondary-600 text-secondary-600 hover:bg-secondary-50 focus:ring-secondary-500',
        ghost: 'text-secondary-600 hover:bg-secondary-50 focus:ring-secondary-500'
      },
      success: {
        filled: 'bg-success-600 text-white hover:bg-success-700 focus:ring-success-500',
        outlined: 'border border-success-600 text-success-600 hover:bg-success-50 focus:ring-success-500',
        ghost: 'text-success-600 hover:bg-success-50 focus:ring-success-500'
      },
      warning: {
        filled: 'bg-warning-600 text-white hover:bg-warning-700 focus:ring-warning-500',
        outlined: 'border border-warning-600 text-warning-600 hover:bg-warning-50 focus:ring-warning-500',
        ghost: 'text-warning-600 hover:bg-warning-50 focus:ring-warning-500'
      },
      error: {
        filled: 'bg-error-600 text-white hover:bg-error-700 focus:ring-error-500',
        outlined: 'border border-error-600 text-error-600 hover:bg-error-50 focus:ring-error-500',
        ghost: 'text-error-600 hover:bg-error-50 focus:ring-error-500'
      },
      gray: {
        filled: 'bg-neutral-500 text-white hover:bg-neutral-600 focus:ring-neutral-400',
        outlined: 'border border-neutral-300 text-neutral-700 hover:bg-neutral-50 focus:ring-neutral-400',
        ghost: 'text-neutral-700 hover:bg-neutral-50 focus:ring-neutral-400'
      }
    };
    
    return colorMap[color][variant] || colorMap.primary.filled;
  })();
  
  // 使用设计令牌中的按钮大小样式
  $: sizeClasses = componentStyles.button.sizes[size] || componentStyles.button.sizes.md;
  
  // 禁用状态样式
  $: disabledClasses = disabled ? 'opacity-50 cursor-not-allowed pointer-events-none' : '';
  
  // 组合所有类名，使用 twMerge 确保外部传入的 classes 具有更高优先级
  $: allClasses = twMerge(baseClasses, variantClasses, sizeClasses, disabledClasses, classes);
  
  // 回调props - 替代createEventDispatcher
  export let onclick: ((event: MouseEvent) => void) | undefined = undefined;
  export let onkeydown: ((event: KeyboardEvent) => void) | undefined = undefined;
  
  // 处理点击事件
  function handleClick(event: MouseEvent) {
    if (disabled || loading) {
      event.preventDefault();
      return;
    }
    onclick?.(event);
  }
  
  // 处理键盘事件
  function handleKeydown(event: KeyboardEvent) {
    if (disabled || loading) return;
    onkeydown?.(event);
  }
</script>

{#if href}
  <!-- 链接按钮 -->
  <a
    {href}
    {target}
    {rel}
    {download}
    class={allClasses}
    class:pointer-events-none={loading}
    onclick={handleClick}
    onkeydown={handleKeydown}
    role="button"
    tabindex={disabled ? -1 : 0}
    aria-disabled={disabled}
    {...$$restProps}
  >
    {#if loading}
      <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    {/if}
    <slot />
  </a>
{:else}
  <!-- 普通按钮 -->
  <button
    {type}
    {disabled}
    class={allClasses}
    class:pointer-events-none={loading}
    onclick={handleClick}
    onkeydown={handleKeydown}
    aria-disabled={disabled}
    {...$$restProps}
  >
    {#if loading}
      <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    {/if}
    <slot />
  </button>
{/if}