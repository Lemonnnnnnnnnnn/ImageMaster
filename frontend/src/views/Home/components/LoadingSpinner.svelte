<script lang="ts">
  import { Loader2 } from 'lucide-svelte';
  
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let text: string = '加载中...';
  export let showText: boolean = true;
  
  $: sizeClass = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12'
  }[size];
</script>

<!-- Fluent Design 加载动画 -->
<div class="flex flex-col items-center justify-center py-16">
  <!-- 主加载动画 -->
  <div class="relative">
    <!-- 外层光环 -->
    <div class="absolute inset-0 bg-fluent-blue/20 rounded-full animate-ping"></div>
    
    <!-- 旋转图标 -->
    <div class="relative bg-glass-card backdrop-blur-fluent {sizeClass} rounded-full flex items-center justify-center border border-white-tertiary/20 shadow-fluent-md">
      <Loader2 size={size === 'sm' ? 12 : size === 'md' ? 20 : 32} class="text-fluent-blue animate-spin" />
    </div>
  </div>
  
  <!-- 加载文本 -->
  {#if showText}
    <div class="mt-6 text-center">
      <p class="text-white-primary font-medium mb-1">{text}</p>
      <p class="text-sm text-white-secondary">请稍候片刻</p>
    </div>
  {/if}
  
  <!-- 进度指示器 -->
  <div class="mt-4 w-48 h-1 bg-white-tertiary/20 rounded-full overflow-hidden">
    <div class="h-full bg-gradient-to-r from-fluent-blue to-fluent-blue/80 rounded-full animate-pulse"></div>
  </div>
</div>