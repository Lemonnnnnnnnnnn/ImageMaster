<script lang="ts">
  import { onMount } from 'svelte';
  import { mangaStore } from './stores/mangaStore';
  import { MangaService } from './services/mangaService';
  import ViewerHeader from './components/ViewerHeader.svelte';
  import ImageViewer from './components/ImageViewer.svelte';
  import NavigationPanel from './components/NavigationPanel.svelte';
  import QuickDownloader from '../../../components/QuickDownloader.svelte';

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
  function handleDownloadStarted(event) {
    const { taskId, url } = event.detail;
    console.log(`下载任务已创建: ${taskId}, URL: ${url}`);
    // 可以在这里添加通知或其他处理逻辑
  }
</script>

<div class="manga-viewer">
  <ViewerHeader />
  <ImageViewer />
  
  {#if showNavigation}
    <NavigationPanel />
  {/if}
  
  <!-- 快速下载组件 -->
  <QuickDownloader on:downloadStarted={handleDownloadStarted} />
</div>

<style>
  .manga-viewer {
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: #1a1a1a;
    color: #f0f0f0;
  }
</style>