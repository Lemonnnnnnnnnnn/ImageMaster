<script lang="ts">
  import { twMerge } from 'tailwind-merge';
  
  // Props - 保持向后兼容
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
  
  // Fluent Design 基础样式
  $: baseClasses = 'fluent-button inline-flex items-center justify-center font-medium transition-fluent focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-black-secondary active:scale-98';
  
  // 构建变体和颜色类名 - 基于 Fluent Design
  $: variantClasses = (() => {
    const colorMap = {
      primary: {
        filled: 'fluent-button-primary focus:ring-fluent-blue/50',
        outlined: 'fluent-button-secondary border-fluent-blue text-fluent-blue hover:bg-fluent-blue/10 focus:ring-fluent-blue/50',
        ghost: 'text-fluent-blue hover:bg-fluent-blue/10 focus:ring-fluent-blue/50'
      },
      secondary: {
        filled: 'bg-fluent-green text-white hover:bg-fluent-green/90 focus:ring-fluent-green/50',
        outlined: 'fluent-button-secondary border-fluent-green text-fluent-green hover:bg-fluent-green/10 focus:ring-fluent-green/50',
        ghost: 'text-fluent-green hover:bg-fluent-green/10 focus:ring-fluent-green/50'
      },
      success: {
        filled: 'bg-fluent-green text-white hover:bg-fluent-green/90 focus:ring-fluent-green/50',
        outlined: 'fluent-button-secondary border-fluent-green text-fluent-green hover:bg-fluent-green/10 focus:ring-fluent-green/50',
        ghost: 'text-fluent-green hover:bg-fluent-green/10 focus:ring-fluent-green/50'
      },
      warning: {
        filled: 'bg-fluent-orange text-white hover:bg-fluent-orange/90 focus:ring-fluent-orange/50',
        outlined: 'fluent-button-secondary border-fluent-orange text-fluent-orange hover:bg-fluent-orange/10 focus:ring-fluent-orange/50',
        ghost: 'text-fluent-orange hover:bg-fluent-orange/10 focus:ring-fluent-orange/50'
      },
      error: {
        filled: 'bg-fluent-red text-white hover:bg-fluent-red/90 focus:ring-fluent-red/50',
        outlined: 'fluent-button-secondary border-fluent-red text-fluent-red hover:bg-fluent-red/10 focus:ring-fluent-red/50',
        ghost: 'text-fluent-red hover:bg-fluent-red/10 focus:ring-fluent-red/50'
      },
      gray: {
        filled: 'bg-white-tertiary text-white-primary hover:bg-white-tertiary/80 focus:ring-white-tertiary/50',
        outlined: 'fluent-button-secondary border-white-tertiary text-white-secondary hover:bg-white-tertiary/10 focus:ring-white-tertiary/50',
        ghost: 'text-white-secondary hover:bg-white-tertiary/10 focus:ring-white-tertiary/50'
      }
    };
    
    return colorMap[color][variant] || colorMap.primary.filled;
  })();
  
  // Fluent Design 尺寸样式
  $: sizeClasses = (() => {
    const sizes = {
      sm: 'px-3 py-1.5 text-sm rounded-fluent-sm min-h-[32px]',
      md: 'px-4 py-2 text-sm rounded-fluent-md min-h-[40px]',
      lg: 'px-6 py-3 text-base rounded-fluent-lg min-h-[48px]'
    };
    return sizes[size] || sizes.md;
  })();
  
  // 禁用状态样式
  $: disabledClasses = disabled ? 'opacity-50 cursor-not-allowed pointer-events-none' : '';
  
  // 加载状态样式
  $: loadingClasses = loading ? 'cursor-wait' : '';
  
  // 组合所有类名
  $: allClasses = twMerge(baseClasses, variantClasses, sizeClasses, disabledClasses, loadingClasses, classes);
  
  // 回调props
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
    onclick={handleClick}
    onkeydown={handleKeydown}
    role="button"
    tabindex={disabled ? -1 : 0}
    aria-disabled={disabled}
    {...$$restProps}
  >
    {#if loading}
      <div class="animate-spin -ml-1 mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full" role="status" aria-label="加载中"></div>
    {/if}
    <slot />
  </a>
{:else}
  <!-- 普通按钮 -->
  <button
    {type}
    {disabled}
    class={allClasses}
    onclick={handleClick}
    onkeydown={handleKeydown}
    aria-disabled={disabled}
    {...$$restProps}
  >
    {#if loading}
      <div class="animate-spin -ml-1 mr-2 h-4 w-4 border-2 border-current border-t-transparent rounded-full" role="status" aria-label="加载中"></div>
    {/if}
    <slot />
  </button>
{/if}