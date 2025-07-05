<script lang="ts">
  import { onMount } from 'svelte';
  import { Menu } from 'lucide-svelte';
  import Sidebar from './Sidebar.svelte';

  // 响应式状态管理
  let innerWidth = 0;
  let sidebarCollapsed = false;
  let isMobile = false;

  // 响应式断点
  $: isMobile = innerWidth < 768;
  $: if (isMobile) {
    sidebarCollapsed = true;
  }

  // 切换侧边栏状态
  function toggleSidebar() {
    sidebarCollapsed = !sidebarCollapsed;
  }

  // 处理键盘快捷键
  function handleKeyDown(event: KeyboardEvent) {
    // Ctrl/Cmd + B 切换侧边栏
    if ((event.ctrlKey || event.metaKey) && event.key === 'b') {
      event.preventDefault();
      toggleSidebar();
    }
    // ESC 在移动端关闭侧边栏
    if (event.key === 'Escape' && isMobile && !sidebarCollapsed) {
      sidebarCollapsed = true;
    }
  }

  onMount(() => {
    // 监听窗口大小变化
    const handleResize = () => {
      innerWidth = window.innerWidth;
    };
    
    window.addEventListener('resize', handleResize);
    handleResize(); // 初始化

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  });
</script>

<svelte:window bind:innerWidth on:keydown={handleKeyDown} />

<!-- 主布局容器 -->
<div class="h-screen flex bg-black-secondary text-white-primary">
  <!-- 侧边栏 -->
  <Sidebar 
    bind:isCollapsed={sidebarCollapsed}
    {isMobile}
    on:toggle={toggleSidebar}
  />

  <!-- 主内容区域 -->
  <main 
    class="transition-all duration-fluent-medium ease-fluent-ease-out min-h-screen"
    class:ml-nav-collapsed={sidebarCollapsed && !isMobile}
    class:ml-nav-expanded={!sidebarCollapsed && !isMobile}
    class:ml-0={isMobile}
  >
    <!-- 移动端顶部栏 -->
    {#if isMobile}
      <div class="sticky top-0 z-30 backdrop-fluent border-b border-white-tertiary/20 px-4 py-3">
        <div class="flex items-center justify-between">
          <button
            class="p-2 rounded-fluent-md text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent"
            on:click={toggleSidebar}
            aria-label="打开导航菜单"
          >
            <Menu size={24} />
          </button>
          
          <h1 class="text-white-primary font-semibold text-lg">图像大师</h1>
          
          <div class="w-10 h-10"></div> <!-- 占位符保持居中 -->
        </div>
      </div>
    {/if}

    <!-- 页面内容插槽 -->
    <div class="p-6">
      <slot />
    </div>
  </main>
</div>

<style>
  /* 确保主内容区域不被侧边栏遮挡 */
  main {
    position: relative;
  }

  /* 平滑的布局过渡 */
  @media (prefers-reduced-motion: no-preference) {
    main {
      transition-property: margin-left;
    }
  }

  /* 移动端优化 */
  @media (max-width: 767px) {
    main {
      margin-left: 0 !important;
    }
  }
</style> 