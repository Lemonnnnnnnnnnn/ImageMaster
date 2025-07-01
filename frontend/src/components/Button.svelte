<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { twMerge } from 'tailwind-merge';
  
  const dispatch = createEventDispatcher();
  
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
  
  // 基础样式
  $: baseClasses = 'inline-flex items-center justify-center font-medium rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2';
  
  // 构建变体和颜色类名
  $: variantClasses = (() => {
    const colorMap = {
      primary: {
        filled: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
        outlined: 'border border-blue-600 text-blue-600 hover:bg-blue-50 focus:ring-blue-500',
        ghost: 'text-blue-600 hover:bg-blue-50 focus:ring-blue-500'
      },
      secondary: {
        filled: 'bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500',
        outlined: 'border border-gray-600 text-gray-600 hover:bg-gray-50 focus:ring-gray-500',
        ghost: 'text-gray-600 hover:bg-gray-50 focus:ring-gray-500'
      },
      success: {
        filled: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500',
        outlined: 'border border-green-600 text-green-600 hover:bg-green-50 focus:ring-green-500',
        ghost: 'text-green-600 hover:bg-green-50 focus:ring-green-500'
      },
      warning: {
        filled: 'bg-yellow-600 text-white hover:bg-yellow-700 focus:ring-yellow-500',
        outlined: 'border border-yellow-600 text-yellow-600 hover:bg-yellow-50 focus:ring-yellow-500',
        ghost: 'text-yellow-600 hover:bg-yellow-50 focus:ring-yellow-500'
      },
      error: {
        filled: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
        outlined: 'border border-red-600 text-red-600 hover:bg-red-50 focus:ring-red-500',
        ghost: 'text-red-600 hover:bg-red-50 focus:ring-red-500'
      },
      gray: {
        filled: 'bg-gray-500 text-white hover:bg-gray-600 focus:ring-gray-400',
        outlined: 'border border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-400',
        ghost: 'text-gray-700 hover:bg-gray-50 focus:ring-gray-400'
      }
    };
    
    return colorMap[color][variant] || colorMap.primary.filled;
  })();
  
  // 构建大小类名
  $: sizeClasses = (() => {
    switch (size) {
      case 'sm': return 'px-3 py-1.5 text-sm';
      case 'md': return 'px-4 py-2 text-sm';
      case 'lg': return 'px-6 py-3 text-base';
      default: return 'px-4 py-2 text-sm';
    }
  })();
  
  // 禁用状态样式
  $: disabledClasses = disabled ? 'opacity-50 cursor-not-allowed pointer-events-none' : '';
  
  // 组合所有类名，使用 twMerge 确保外部传入的 classes 具有更高优先级
  $: allClasses = twMerge(baseClasses, variantClasses, sizeClasses, disabledClasses, classes);
  
  // 处理点击事件
  function handleClick(event: MouseEvent) {
    if (disabled || loading) {
      event.preventDefault();
      return;
    }
    dispatch('click', event);
  }
  
  // 处理键盘事件
  function handleKeydown(event: KeyboardEvent) {
    if (disabled || loading) return;
    dispatch('keydown', event);
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
    on:click={handleClick}
    on:keydown={handleKeydown}
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
    on:click={handleClick}
    on:keydown={handleKeydown}
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