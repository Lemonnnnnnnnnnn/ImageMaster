<script lang="ts">
  import { componentStyles } from '../utils/designTokens';
  import { twMerge } from 'tailwind-merge';
  
  // Additional classes
  export let classes: string = '';
  export let hover: boolean = false;
  
  // 回调props - 替代createEventDispatcher
  export let onclick: ((event: MouseEvent) => void) | undefined = undefined;
  
  // 使用设计令牌中的卡片样式
  $: baseClasses = componentStyles.card.base;
  $: hoverClasses = hover ? componentStyles.card.hover : '';
  
  // 组合所有类名，使用 twMerge 确保外部传入的 classes 具有更高优先级
  $: allClasses = twMerge(baseClasses, hoverClasses, classes);
  
  // 处理点击事件，确保事件能够正确传递
  function handleClick(event: MouseEvent) {
    onclick?.(event);
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class={allClasses} on:click={handleClick} {...$$restProps}>
  <slot />
</div>