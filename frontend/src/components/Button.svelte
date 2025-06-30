<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  // Props
  export let variant: 'filled' | 'tonal' | 'outlined' = 'filled';
  export let color: 'primary' | 'secondary' | 'tertiary' | 'success' | 'warning' | 'error' | 'surface' = 'primary';
  export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
  export let disabled: boolean = false;
  export let loading: boolean = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let href: string | undefined = undefined;
  export let target: string | undefined = undefined;
  export let rel: string | undefined = undefined;
  export let download: string | boolean | undefined = undefined;
  
  // Additional classes
  export let classes: string = '';
  
  // 构建基础类名
  $: baseClasses = 'btn';
  
  // 构建变体类名 - 使用静态类名映射避免动态生成
  $: variantClasses = (() => {
    const classMap = {
      // Filled variants
      'filled-primary': 'preset-filled-primary-500',
      'filled-secondary': 'preset-filled-secondary-500',
      'filled-tertiary': 'preset-filled-tertiary-500',
      'filled-success': 'preset-filled-success-500',
      'filled-warning': 'preset-filled-warning-500',
      'filled-error': 'preset-filled-error-500',
      'filled-surface': 'preset-filled-surface-500',
      
      // Tonal variants
      'tonal-primary': 'preset-tonal-primary',
      'tonal-secondary': 'preset-tonal-secondary',
      'tonal-tertiary': 'preset-tonal-tertiary',
      'tonal-success': 'preset-tonal-success',
      'tonal-warning': 'preset-tonal-warning',
      'tonal-error': 'preset-tonal-error',
      'tonal-surface': 'preset-tonal-surface',
      
      // Outlined variants
      'outlined-primary': 'preset-outlined-primary-500',
      'outlined-secondary': 'preset-outlined-secondary-500',
      'outlined-tertiary': 'preset-outlined-tertiary-500',
      'outlined-success': 'preset-outlined-success-500',
      'outlined-warning': 'preset-outlined-warning-500',
      'outlined-error': 'preset-outlined-error-500',
      'outlined-surface': 'preset-outlined-surface-500'
    };
    
    const key = `${variant}-${color}` as keyof typeof classMap;
    return classMap[key] || 'preset-filled-primary-500'; // 默认值
  })();
  
  // 构建大小类名
  $: sizeClasses = (() => {
    switch (size) {
      case 'sm': return 'btn-sm';
      case 'md': return 'btn-md';
      case 'lg': return 'btn-lg';
      case 'xl': return 'btn-xl';
      default: return 'btn-md';
    }
  })();
  
  // 组合所有类名
  $: allClasses = [baseClasses, variantClasses, sizeClasses, classes].filter(Boolean).join(' ');
  
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
    class:opacity-50={disabled}
    class:cursor-not-allowed={disabled}
    class:pointer-events-none={disabled || loading}
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
    class:opacity-50={disabled}
    class:cursor-not-allowed={disabled}
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