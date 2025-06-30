<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import {
      CancelDownload,
      ClearHistory,
      GetActiveTasks,
      GetHistoryTasks
  } from '../../../wailsjs/go/downloader/DownloaderAPI';
  import type { downloader } from '../../../wailsjs/go/models';
  import {
      GetOutputDir
  } from '../../../wailsjs/go/storage/StorageAPI';
  import Header from '../../components/Header.svelte';
  import { createDownloadHandler } from '../../utils/downloadUtils';

  let url = '';
  let loading = false;
  let outputDir = '';
  let error = '';
  
  // 创建下载处理器
  const downloadHandler = createDownloadHandler({
    onStart: () => {
      loading = true;
      error = '';
    },
    onSuccess: async (taskId, downloadUrl) => {
      // 重置表单
      url = '';
      
      // 切换到下载中标签并刷新任务列表
      activeTab = 'downloading';
      await pollTasks();
    },
    onError: (errorMsg) => {
      error = errorMsg;
    },
    onFinally: () => {
      loading = false;
    }
  });
  
  // 下载任务列表
  let activeTasks: downloader.DownloadTask[] = [];
  let historyTasks: downloader.DownloadTask[] = [];
  let activeTab = 'downloading'; // 'downloading' or 'history'
  
  // 轮询间隔（毫秒）
  const POLL_INTERVAL = 1000;
  let pollTimer: ReturnType<typeof setInterval> | null = null;

  onMount(async () => {
    // 获取当前输出目录
    outputDir = await GetOutputDir();
    
    // 开始轮询
    startPolling();
    
    console.log("组件已挂载，轮询已开始");
  });
  
  onDestroy(() => {
    // 停止轮询
    stopPolling();
    console.log("组件已销毁，轮询已停止");
  });
  
  // 开始轮询
  function startPolling() {
    // 清除现有轮询
    stopPolling();
    
    // 立即执行一次轮询
    pollTasks();
    
    // 设置定时轮询
    pollTimer = setInterval(pollTasks, POLL_INTERVAL);
  }
  
  // 停止轮询
  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer);
      pollTimer = null;
    }
  }
  
  // 轮询任务状态
  async function pollTasks() {
    try {
      // 获取活跃任务
      activeTasks = await GetActiveTasks();
      
      // 如果当前显示历史标签，也获取历史任务
      if (activeTab === 'history') {
        historyTasks = await GetHistoryTasks();
      }
    } catch (err) {
      console.error("轮询任务状态出错:", err);
    }
  }

  async function handleSubmit() {
    if (!url.trim()) {
      error = '请输入网址';
      return;
    }

    await downloadHandler(url.trim());
  }
  
  // 取消下载任务
  async function cancelTask(taskId: string) {
    try {
      await CancelDownload(taskId);
      // 立即刷新任务列表
      await pollTasks();
    } catch (err) {
      console.error("取消任务出错:", err);
    }
  }
  
  // 切换标签
  async function switchTab(tab: string) {
    activeTab = tab;
    
    // 如果切换到历史标签，获取历史任务
    if (tab === 'history') {
      historyTasks = await GetHistoryTasks();
    }
  }
  
  // 清除历史记录
  async function handleClearHistory() {
    try {
      await ClearHistory();
      historyTasks = [];
    } catch (err) {
      console.error("清除历史出错:", err);
    }
  }
  
  // 格式化任务状态
  function formatStatus(status: string): string {
    const statusMap : Record<string,string> = {
      'pending': '等待中',
      'downloading': '下载中',
      'completed': '已完成',
      'failed': '失败',
      'cancelled': '已取消'
    };
    return statusMap[status] || status;
  }
  
  // 格式化时间
  function formatTime(timeStr: string): string {
    if (!timeStr) return '';
    const date = new Date(timeStr);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
  }
</script>

<div class="p-5 max-w-[95%] w-full mx-auto font-sans text-gray-800 h-screen flex flex-col gap-5 box-border">
  <Header title="网页图片下载器" />

  <!-- 顶部面板：输出目录 + 下载表单 -->
  <div class="bg-white rounded-md p-5 border border-gray-300 shadow-sm flex-shrink-0">
    <!-- 输出目录 -->
    <div class="flex items-center gap-4 mb-5 pb-4 border-b border-gray-300 md:flex-col md:items-stretch md:gap-2.5 md:mb-4 md:pb-3">
      <label class="font-medium text-gray-600 whitespace-nowrap min-w-[80px] md:min-w-auto">输出目录:</label>
      <div class="flex items-center gap-3 flex-1 bg-gray-50 px-4 py-3 rounded-md border border-gray-300">
        <span class="flex-1 text-sm text-gray-600 overflow-hidden text-ellipsis whitespace-nowrap">{outputDir}</span>
        <button on:click={()=>push('/config')} class="px-4 py-2 bg-blue-600 text-white border-none rounded-md cursor-pointer text-xs font-medium whitespace-nowrap hover:bg-purple-700">设置</button>
      </div>
    </div>
    
    <!-- 下载表单 -->
    <div class="flex flex-col gap-3">
      <div class="flex gap-4 items-center md:flex-col md:gap-3">
        <input 
          type="text" 
          bind:value={url} 
          placeholder="输入网页完整地址，例如: https://example.com/gallery"
          disabled={loading}
          class="flex-1 px-4 py-3 border border-gray-300 rounded-md text-base focus:outline-none focus:border-blue-600 focus:ring-2 focus:ring-blue-200 md:w-full"
          on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
        />
        <button on:click={handleSubmit} disabled={loading} class="px-6 py-3 bg-green-500 text-white border-none rounded-md cursor-pointer text-base font-medium whitespace-nowrap hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed md:w-full">
          {loading ? '添加中...' : '添加任务'}
        </button>
      </div>
      
      {#if error}
        <div class="bg-red-50 text-red-700 px-4 py-2.5 rounded-md border-l-4 border-red-700 text-sm">
          {error}
        </div>
      {/if}
    </div>
  </div>
  
  <!-- 任务管理面板 -->
  <div class="bg-white rounded-md border border-gray-300 shadow-sm flex-1 flex flex-col min-h-0">
    <!-- Tab导航 -->
    <div class="flex items-center gap-1.5 px-5 py-4 border-b border-gray-300 bg-gray-50 rounded-t-md flex-shrink-0">
      <button 
        class="px-5 py-2.5 border-none cursor-pointer rounded-md text-base font-medium text-gray-600 transition-all duration-200 {activeTab === 'downloading' ? 'bg-gray-200 ' : 'hover:bg-gray-200'}"
        on:click={() => switchTab('downloading')}
      >
        下载中 ({activeTasks.length})
      </button>
      <button 
        class="px-5 py-2.5 border-none cursor-pointer rounded-md text-base font-medium text-gray-600 transition-all duration-200 {activeTab === 'history' ? 'bg-gray-200 ' : 'hover:bg-gray-200'}"
        on:click={() => switchTab('history')}
      >
        历史记录
      </button>
      {#if activeTab === 'history' && historyTasks.length > 0}
        <button on:click={handleClearHistory} class="px-4 py-2 bg-red-600 text-white border-none rounded-md cursor-pointer text-xs font-medium ml-auto hover:bg-red-700">
          清除历史
        </button>
      {/if}
    </div>
    
    <!-- Tab内容 -->
    <div class="flex-1 overflow-auto p-5">
      {#if activeTab === 'downloading'}
        <!-- 下载中任务 -->
        {#if activeTasks.length > 0}
          <div class="flex flex-col gap-4">
            {#each activeTasks as task (task.id)}
              <div class="bg-gray-50 rounded-md p-4 border border-gray-300 transition-shadow duration-200 hover:shadow-md">
                <div class="grid grid-cols-[1fr_auto_auto] gap-4 items-center mb-2.5 md:grid-cols-1 md:gap-2">
                  <div class="text-sm text-gray-600 overflow-hidden text-ellipsis whitespace-nowrap">{task.url}</div>
                  <div class="text-xs font-medium px-3 py-1.5 rounded-full whitespace-nowrap {task.status === 'downloading' || task.status === 'pending' ? 'bg-blue-600 text-white' : task.status === 'completed' ? 'bg-green-500 text-white' : 'bg-red-600 text-white'}">{formatStatus(task.status)}</div>
                  <div class="text-xs text-gray-600 whitespace-nowrap md:justify-self-start">{formatTime(task.startTime)}</div>
                </div>
                
                {#if task.status === 'downloading' || task.status === 'pending'}
                  <!-- 进度条 -->
                  <div class="mt-3">
                    <div class="w-full h-2 bg-gray-300 rounded overflow-hidden mb-2">
                      <div class="h-full bg-blue-600 transition-all duration-300" style="width: {task.progress.total > 0 ? Math.round((task.progress.current / task.progress.total) * 100) : 0}%"></div>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-xs text-gray-600">
                        {#if task.progress.total > 0}
                          {task.progress.current}/{task.progress.total} 张图片 ({Math.round((task.progress.current / task.progress.total) * 100)}%)
                        {:else}
                          准备下载中...
                        {/if}
                      </span>
                      <button class="px-3 py-1.5 bg-red-600 text-white border-none rounded-md cursor-pointer text-xs font-medium hover:bg-red-700" on:click={() => cancelTask(task.id)}>
                        取消
                      </button>
                    </div>
                  </div>
                {/if}
                
                {#if task.status === 'failed'}
                  <div class="mt-2.5 px-3 py-2.5 rounded-md text-xs bg-red-50 text-red-700 border-l-4 border-red-700">
                    错误: {task.error || '未知错误'}
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <div class="text-center text-gray-600 italic py-15 px-5 text-base">
            当前没有进行中的下载任务
          </div>
        {/if}
      {:else if activeTab === 'history'}
        <!-- 历史任务 -->
        {#if historyTasks.length > 0}
          <div class="flex flex-col gap-4">
            {#each historyTasks as task (task.id)}
              <div class="bg-gray-50 rounded-md p-4 border border-gray-300 transition-shadow duration-200 hover:shadow-md opacity-90">
                <div class="grid grid-cols-[1fr_auto_auto] gap-4 items-center mb-2.5 md:grid-cols-1 md:gap-2">
                  <div class="text-sm text-gray-600 overflow-hidden text-ellipsis whitespace-nowrap">{task.url}</div>
                  <div class="text-xs font-medium px-3 py-1.5 rounded-full whitespace-nowrap {task.status === 'downloading' || task.status === 'pending' ? 'bg-blue-600 text-white' : task.status === 'completed' ? 'bg-green-500 text-white' : 'bg-red-600 text-white'}">{formatStatus(task.status)}</div>
                  <div class="text-xs text-gray-600 whitespace-nowrap md:justify-self-start">
                    {formatTime(task.startTime)} - {formatTime(task.completeTime)}
                  </div>
                </div>
                
                {#if task.status === 'completed' && task.savePath}
                  <div class="mt-2.5 px-3 py-2.5 rounded-md text-xs bg-green-50 text-green-700 border-l-4 border-green-700">
                    保存路径: {task.savePath}
                  </div>
                {/if}
                
                {#if task.status === 'failed' && task.error}
                  <div class="mt-2.5 px-3 py-2.5 rounded-md text-xs bg-red-50 text-red-700 border-l-4 border-red-700">
                    错误: {task.error}
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <div class="text-center text-gray-600 italic py-15 px-5 text-base">
            暂无历史记录
          </div>
        {/if}
      {/if}
    </div>
  </div>
</div>