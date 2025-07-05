<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Header from '../../components/Header.svelte';
  import Card from '../../components/Card.svelte';
  import OutputDirectory from './components/OutputDirectory.svelte';
  import DownloadForm from './components/DownloadForm.svelte';
  import TaskTabs from './components/TaskTabs.svelte';
  import ActiveTaskList from './components/ActiveTaskList.svelte';
  import HistoryTaskList from './components/HistoryTaskList.svelte';
  import { activeTab, initializeStore, stopPolling } from './stores/downloadStore';

  onMount(async () => {
    await initializeStore();
    console.log("组件已挂载，轮询已开始");
  });
  
  onDestroy(() => {
    stopPolling();
    console.log("组件已销毁，轮询已停止");
  });
</script>

<!-- Fluent Design 黑色主题布局 -->
<div class="min-h-screen bg-black-secondary">
  <Header 
    title="网页图片下载器" 
    subtitle="从支持的网站下载图片和漫画"
  />

  <!-- 主要内容区域 -->
  <div class="container mx-auto max-w-7xl px-6 py-8 space-y-6">
    <!-- 顶部面板：输出目录 + 下载表单 -->
    <Card variant="elevated" padding="lg" classes="transition-fluent hover-lift">
      <OutputDirectory />
      <div class="mt-6 pt-6 border-t border-white-tertiary/20">
        <DownloadForm />
      </div>
    </Card>
    
    <!-- 任务管理面板 -->
    <Card variant="elevated" padding="none" classes="flex flex-col min-h-[500px]">
      <div class="border-b border-white-tertiary/20">
        <TaskTabs />
      </div>
      
      <!-- Tab内容 -->
      <div class="flex-1 overflow-auto p-6">
        {#if $activeTab === 'downloading'}
          <ActiveTaskList />
        {:else if $activeTab === 'history'}
          <HistoryTaskList />
        {/if}
      </div>
    </Card>
  </div>
</div>