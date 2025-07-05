<script lang="ts">
  import { twMerge } from 'tailwind-merge';

  // Props
  export let variant: 'default' | 'elevated' | 'outlined' = 'default';
  export let padding: 'none' | 'sm' | 'md' | 'lg' = 'md';
  export let hoverable: boolean = false;
  export let clickable: boolean = false;
  export let classes: string = '';

  // 回调事件
  export let onclick: ((event: MouseEvent) => void) | undefined = undefined;

  // Fluent Design 基础样式
  $: baseClasses = 'fluent-card transition-fluent';

  // 变体样式
  $: variantClasses = (() => {
    const variants = {
      default: 'bg-glass-card backdrop-blur-fluent border border-white-tertiary/20',
      elevated: 'bg-glass-card backdrop-blur-fluent border border-white-tertiary/20 shadow-fluent-lg',
      outlined: 'bg-transparent border-2 border-white-tertiary/30'
    };
    return variants[variant] || variants.default;
  })();

  // 内边距样式
  $: paddingClasses = (() => {
    const paddings = {
      none: '',
      sm: 'p-3',
      md: 'p-4',
      lg: 'p-6'
    };
    return paddings[padding] || paddings.md;
  })();

  // 交互样式
  $: interactiveClasses = (() => {
    if (clickable) {
      return 'cursor-pointer hover-lift focus:outline-none focus:ring-2 focus:ring-fluent-blue/50 focus:ring-offset-2 focus:ring-offset-black-secondary';
    }
    if (hoverable) {
      return 'hover-lift';
    }
    return '';
  })();

  // 组合所有类名
  $: allClasses = twMerge(baseClasses, variantClasses, paddingClasses, interactiveClasses, classes);

  // 处理点击事件
  function handleClick(event: MouseEvent) {
    if (clickable && onclick) {
      onclick(event);
    }
  }

  // 处理键盘事件
  function handleKeyDown(event: KeyboardEvent) {
    if (clickable && (event.key === 'Enter' || event.key === ' ')) {
      event.preventDefault();
      if (onclick) {
        onclick(event as unknown as MouseEvent);
      }
    }
  }
</script>

{#if clickable}
  <!-- 可点击卡片 -->
  <div
    class={allClasses}
    onclick={handleClick}
    onkeydown={handleKeyDown}
    role="button"
    tabindex="0"
    {...$$restProps}
  >
    <slot />
  </div>
{:else}
  <!-- 普通卡片 -->
  <div class={allClasses} {...$$restProps}>
    <slot />
  </div>
{/if}