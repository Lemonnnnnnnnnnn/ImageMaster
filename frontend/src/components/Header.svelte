<script lang="ts">
  // 页面标题
  export let title: string = '图像大师';
  
  // 可选的子标题或描述
  export let subtitle: string = '';
  
  // 可选的动作按钮插槽内容
  export let showActions: boolean = false;
</script>

<!-- Fluent Design 页面标题栏 -->
<header class="sticky top-0 z-40 backdrop-fluent border-b border-white-tertiary/20 px-6 py-4">
  <div class="flex items-center justify-between max-w-none">
    <!-- 标题区域 -->
    <div class="flex flex-col gap-1">
      <h1 class="text-2xl font-semibold text-white-primary m-0 tracking-tight">
        {title}
      </h1>
      {#if subtitle}
        <p class="text-sm text-white-secondary m-0">
          {subtitle}
        </p>
      {/if}
    </div>
    
    <!-- 动作按钮区域 -->
    {#if showActions}
      <div class="flex items-center gap-3">
        <slot name="actions" />
      </div>
    {/if}
  </div>
</header>