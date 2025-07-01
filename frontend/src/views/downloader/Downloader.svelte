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

<div class="p-5 max-w-[95%] w-full mx-auto font-sans text-gray-800 h-screen flex flex-col gap-5 box-border">
  <Header title="网页图片下载器" />

  <!-- 顶部面板：输出目录 + 下载表单 -->
  <Card classes="p-5 border border-gray-300 shadow-sm flex-shrink-0">
    <OutputDirectory />
    <DownloadForm />
  </Card>
  
  <!-- 任务管理面板 -->
  <Card classes="border border-gray-300 shadow-sm flex-1 flex flex-col min-h-0">
    <TaskTabs />
    
    <!-- Tab内容 -->
    <div class="flex-1 overflow-auto p-5">
      {#if $activeTab === 'downloading'}
        <ActiveTaskList />
      {:else if $activeTab === 'history'}
        <HistoryTaskList />
      {/if}
    </div>
  </Card>
</div>