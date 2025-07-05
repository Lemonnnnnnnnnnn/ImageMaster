<script lang="ts">
  import { Library } from 'lucide-svelte';
  
  // 可选props
  export let title: string = '漫画查看器';
  export let subtitle: string = '浏览和查看您的漫画收藏';
</script>

<!-- Fluent Design 页面标题 -->
<div class="flex items-center justify-between mb-8 pb-4 border-b border-white-tertiary/20">
  <div class="flex items-center gap-4">
    <!-- 图标 -->
    <div class="w-12 h-12 bg-fluent-blue/20 rounded-fluent-lg flex items-center justify-center">
      <Library size={24} class="text-fluent-blue" />
    </div>
    
    <!-- 标题和副标题 -->
    <div>
      <h1 class="text-2xl font-semibold text-white-primary m-0 tracking-tight">
        {title}
      </h1>
      {#if subtitle}
        <p class="text-sm text-white-secondary m-0 mt-1">
          {subtitle}
        </p>
      {/if}
    </div>
  </div>
  
  <!-- 右侧插槽 -->
  <div class="flex items-center gap-3">
    <slot />
  </div>
</div>