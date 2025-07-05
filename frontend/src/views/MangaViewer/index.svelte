<script lang="ts">
  import { onMount } from 'svelte';
  import { mangaStore } from './stores/mangaStore';
  import { MangaService } from './services/mangaService';
  import ViewerHeader from './components/ViewerHeader.svelte';
  import ImageViewer from './components/ImageViewer.svelte';
  import NavigationPanel from './components/NavigationPanel.svelte';
  import QuickDownloader from '../../components/QuickDownloader.svelte';

  // 从路由参数获取漫画路径
  export let params: {
    path?: string
  } = {};
  
  $: ({ showNavigation } = $mangaStore);
  
  // 响应式声明，监听 params 变化
  $: if (params.path) {
    MangaService.loadManga(params.path);
  }

  onMount(() => {
    // 添加键盘事件监听
    window.addEventListener('keydown', MangaService.handleKeyDown);
    return () => window.removeEventListener('keydown', MangaService.handleKeyDown);
  });

  // 处理快速下载成功事件
  function handleDownloadStarted(event: any) {
    const { taskId, url } = event.detail;
    console.log(`下载任务已创建: ${taskId}, URL: ${url}`);
    // 可以在这里添加通知或其他处理逻辑
  }
</script>

<!-- Fluent Design 黑色主题漫画查看器 -->
<div class="h-screen flex flex-col bg-black-primary text-white-primary overflow-hidden">
  <!-- 查看器头部 -->
  <ViewerHeader />
  
  <!-- 主要图片查看区域 -->
  <div class="flex-1 relative">
    <ImageViewer />
    
    <!-- 导航面板叠加层 -->
    {#if showNavigation}
      <div class="absolute inset-0 z-40 bg-black/60 backdrop-blur-sm">
        <NavigationPanel />
      </div>
    {/if}
  </div>
  
  <!-- 快速下载组件 -->
  <QuickDownloader ondownloadstarted={handleDownloadStarted} />
</div>