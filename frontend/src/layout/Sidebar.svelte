<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { push, location } from 'svelte-spa-router';
  import { 
    Home, 
    Download, 
    Settings, 
    Menu,
    X,
    ChevronLeft,
    ChevronRight
  } from 'lucide-svelte';

  // 事件分发器
  const dispatch = createEventDispatcher();

  // 导航栏状态
  export let isCollapsed = false;
  export let isMobile = false;

  // 导航项配置
  const navigationItems = [
    {
      id: 'viewer',
      label: '漫画查看器',
      icon: Home,
      path: '/',
      description: '浏览和查看漫画'
    },
    {
      id: 'downloader', 
      label: '图片下载器',
      icon: Download,
      path: '/downloader',
      description: '下载图片和漫画'
    },
    {
      id: 'config',
      label: '应用设置',
      icon: Settings,
      path: '/config',
      description: '配置应用选项'
    }
  ];

  // 获取当前激活的路由
  $: currentPath = $location || '/';
  
  // 判断导航项是否激活
  function isActive(path: string): boolean {
    if (path === '/') {
      return currentPath === '/' || currentPath.startsWith('/manga/');
    }
    return currentPath.startsWith(path);
  }

  // 导航到指定路径
  function navigateTo(path: string) {
    push(path);
    // 在移动端点击导航后自动收起侧边栏
    if (isMobile) {
      dispatch('toggle');
    }
  }

  // 切换侧边栏状态
  function toggleSidebar() {
    dispatch('toggle');
  }

  // 键盘导航支持
  function handleKeyDown(event: KeyboardEvent, path: string) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      navigateTo(path);
    }
  }

  let 
</script>

<!-- 侧边栏容器 -->
<aside 
  class="h-full z-50 transition-all duration-fluent-medium ease-fluent-ease-out backdrop-fluent border-r border-white-tertiary/20"
  class:w-nav-collapsed={isCollapsed}
  class:w-nav-expanded={!isCollapsed}
  class:-translate-x-full={isMobile && isCollapsed}
  class:translate-x-0={!isMobile || !isCollapsed}
  role="navigation"
  aria-label="主导航"
>
  <!-- 侧边栏头部 -->
  <div class="flex items-center justify-between h-16 px-4 border-b border-white-tertiary/20">
    <!-- 折叠/展开按钮 -->
    <button
      class="p-2 rounded-fluent-md text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent"
      class:hidden={isMobile}
      on:click={toggleSidebar}
      aria-label={isCollapsed ? '展开侧边栏' : '收起侧边栏'}
      title={isCollapsed ? '展开侧边栏' : '收起侧边栏'}
    >
      {#if isCollapsed}
        <ChevronRight size={20} />
      {:else}
        <ChevronLeft size={20} />
      {/if}
    </button>

    <!-- 移动端关闭按钮 -->
    {#if isMobile && !isCollapsed}
      <button
        class="p-2 rounded-fluent-md text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent"
        on:click={toggleSidebar}
        aria-label="关闭侧边栏"
      >
        <X size={20} />
      </button>
    {/if}
  </div>

  <!-- 导航菜单 -->
  <nav class="flex-1 py-4 overflow-y-auto scrollbar-thin scrollbar-track-transparent scrollbar-thumb-white-tertiary/30">
    <ul class="space-y-1" role="menubar">
      {#each navigationItems as item (item.id)}
        <li role="none">
          <button
            class="fluent-nav-item w-full flex items-center p-2 pl-4 text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent"
            class:active={isActive(item.path)}
            on:click={() => navigateTo(item.path)}
            on:keydown={(e) => handleKeyDown(e, item.path)}
            role="menuitem"
            tabindex="0"
            aria-current={isActive(item.path) ? 'page' : undefined}
            title={isCollapsed ? item.label : item.description}
          >
            <!-- 图标 -->
            <div class="flex-shrink-0">
              <svelte:component 
                this={item.icon} 
                size={24} 
                class="transition-colors duration-fluent-fast text-white-secondary"
              />
            </div>
            
            <!-- 标签文本 -->
            {#if !isCollapsed}
              <span class="ml-3 font-medium truncate">{item.label}</span>
            {/if}
            
          </button>
        </li>
      {/each}
    </ul>
  </nav>

  <!-- 侧边栏底部 -->
  <div class="p-4 border-t border-white-tertiary/20">
    {#if !isCollapsed}
      <div class="text-xs text-white-tertiary text-center">
        <p>ImageMaster v1.0</p>
        <p class="mt-1">Fluent Design</p>
      </div>
    {/if}
  </div>
</aside>

<!-- 移动端遮罩层 -->
{#if isMobile && !isCollapsed}
  <div 
    class="fixed inset-0 bg-black/50 z-40 backdrop-blur-sm"
    on:click={toggleSidebar}
    on:keydown={(e) => e.key === 'Escape' && toggleSidebar()}
    role="button"
    tabindex="0"
    aria-label="关闭侧边栏"
  ></div>
{/if}

<style>
  /* 自定义滚动条样式 */
  .scrollbar-thin::-webkit-scrollbar {
    width: 4px;
  }
  
  .scrollbar-track-transparent::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .scrollbar-thumb-white-tertiary\/30::-webkit-scrollbar-thumb {
    background: rgba(110, 110, 110, 0.3);
    border-radius: 2px;
  }
  
  .scrollbar-thumb-white-tertiary\/30::-webkit-scrollbar-thumb:hover {
    background: rgba(110, 110, 110, 0.5);
  }

  /* 激活状态的导航项样式增强 */
  .fluent-nav-item.active {
    position: relative;
  }
  
  .fluent-nav-item.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 24px;
    background: #fff;
    border-radius: 0 2px 2px 0;
  }

  /* 响应式处理 */
  @media (max-width: 768px) {
    aside {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
    }
  }
</style> 