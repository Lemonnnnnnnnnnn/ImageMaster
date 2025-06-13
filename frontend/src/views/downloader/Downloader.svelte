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
    const statusMap = {
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

<div class="downloader-container">
  <Header title="网页图片下载器" />

  <!-- 顶部面板：输出目录 + 下载表单 -->
  <div class="top-panel">
    <!-- 输出目录 -->
    <div class="output-section">
      <label>输出目录:</label>
      <div class="dir-display">
        <span>{outputDir}</span>
        <button on:click={()=>push('/config')} class="config-btn">设置</button>
      </div>
    </div>
    
    <!-- 下载表单 -->
    <div class="download-section">
      <div class="form-row">
        <input 
          type="text" 
          bind:value={url} 
          placeholder="输入网页完整地址，例如: https://example.com/gallery"
          disabled={loading}
          class="url-input"
          on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
        />
        <button on:click={handleSubmit} disabled={loading} class="download-btn">
          {loading ? '添加中...' : '添加任务'}
        </button>
      </div>
      
      {#if error}
        <div class="error">
          {error}
        </div>
      {/if}
    </div>
  </div>
  
  <!-- 任务管理面板 -->
  <div class="tasks-panel">
    <!-- Tab导航 -->
    <div class="tab-nav">
      <button 
        class="tab-btn {activeTab === 'downloading' ? 'active' : ''}"
        on:click={() => switchTab('downloading')}
      >
        下载中 ({activeTasks.length})
      </button>
      <button 
        class="tab-btn {activeTab === 'history' ? 'active' : ''}"
        on:click={() => switchTab('history')}
      >
        历史记录
      </button>
      {#if activeTab === 'history' && historyTasks.length > 0}
        <button on:click={handleClearHistory} class="clear-btn">
          清除历史
        </button>
      {/if}
    </div>
    
    <!-- Tab内容 -->
    <div class="tab-content">
      {#if activeTab === 'downloading'}
        <!-- 下载中任务 -->
        {#if activeTasks.length > 0}
          <div class="tasks-list">
            {#each activeTasks as task (task.id)}
              <div class="task-item">
                <div class="task-header">
                  <div class="task-url">{task.url}</div>
                  <div class="task-status downloading">{formatStatus(task.status)}</div>
                  <div class="task-time">{formatTime(task.startTime)}</div>
                </div>
                
                {#if task.status === 'downloading' || task.status === 'pending'}
                  <!-- 进度条 -->
                  <div class="progress-container">
                    <div class="progress-bar">
                      <div class="progress-fill" style="width: {task.progress.total > 0 ? Math.round((task.progress.current / task.progress.total) * 100) : 0}%"></div>
                    </div>
                    <div class="progress-info">
                      <span class="progress-text">
                        {#if task.progress.total > 0}
                          {task.progress.current}/{task.progress.total} 张图片 ({Math.round((task.progress.current / task.progress.total) * 100)}%)
                        {:else}
                          准备下载中...
                        {/if}
                      </span>
                      <button class="cancel-btn" on:click={() => cancelTask(task.id)}>
                        取消
                      </button>
                    </div>
                  </div>
                {/if}
                
                {#if task.status === 'failed'}
                  <div class="task-error">
                    错误: {task.error || '未知错误'}
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <div class="no-tasks">
            当前没有进行中的下载任务
          </div>
        {/if}
      {:else if activeTab === 'history'}
        <!-- 历史任务 -->
        {#if historyTasks.length > 0}
          <div class="tasks-list">
            {#each historyTasks as task (task.id)}
              <div class="task-item history">
                <div class="task-header">
                  <div class="task-url">{task.url}</div>
                  <div class="task-status {task.status}">{formatStatus(task.status)}</div>
                  <div class="task-time">
                    {formatTime(task.startTime)} - {formatTime(task.completeTime)}
                  </div>
                </div>
                
                {#if task.status === 'completed' && task.savePath}
                  <div class="task-path">
                    保存路径: {task.savePath}
                  </div>
                {/if}
                
                {#if task.status === 'failed' && task.error}
                  <div class="task-error">
                    错误: {task.error}
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {:else}
          <div class="no-tasks">
            暂无历史记录
          </div>
        {/if}
      {/if}
    </div>
  </div>
</div>

<style>
  :root {
    --primary-color: #4361ee;
    --secondary-color: #3a0ca3;
    --success-color: #4CAF50;
    --danger-color: #e63946;
    --warning-color: #ff9f1c;
    --light-bg: #f8f9fa;
    --border-color: #dee2e6;
    --box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    --border-radius: 6px;
  }

  .downloader-container {
    padding: 20px;
    max-width: 95%;
    width: 100%;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    color: #333;
    height: 100vh;
    display: flex;
    flex-direction: column;
    gap: 20px;
    box-sizing: border-box;
  }
  
  /* 顶部面板 */
  .top-panel {
    background: white;
    border-radius: var(--border-radius);
    padding: 20px;
    border: 1px solid var(--border-color);
    box-shadow: var(--box-shadow);
    flex-shrink: 0;
  }
  
  .output-section {
    display: flex;
    align-items: center;
    gap: 15px;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 1px solid var(--border-color);
  }
  
  .output-section label {
    font-weight: 500;
    color: #555;
    white-space: nowrap;
    min-width: 80px;
  }
  
  .dir-display {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    background: var(--light-bg);
    padding: 12px 15px;
    border-radius: var(--border-radius);
    border: 1px solid var(--border-color);
  }
  
  .dir-display span {
    flex: 1;
    font-size: 14px;
    color: #666;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .config-btn {
    padding: 8px 16px;
    background: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
  }
  
  .config-btn:hover {
    background: var(--secondary-color);
  }
  
  .download-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .form-row {
    display: flex;
    gap: 15px;
    align-items: center;
  }
  
  .url-input {
    flex: 1;
    padding: 12px 15px;
    border: 1px solid var(--border-color);
    border-radius: var(--border-radius);
    font-size: 15px;
  }
  
  .url-input:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(67, 97, 238, 0.2);
  }
  
  .download-btn {
    padding: 12px 24px;
    background: var(--success-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 15px;
    font-weight: 500;
    white-space: nowrap;
  }
  
  .download-btn:hover:not(:disabled) {
    background: #3d8b40;
  }
  
  .download-btn:disabled {
    background: #a5a5a5;
    cursor: not-allowed;
  }
  
  .error {
    background: #ffebee;
    color: #c62828;
    padding: 10px 15px;
    border-radius: var(--border-radius);
    border-left: 3px solid #c62828;
    font-size: 14px;
  }
  
  /* 任务面板 */
  .tasks-panel {
    background: white;
    border-radius: var(--border-radius);
    border: 1px solid var(--border-color);
    box-shadow: var(--box-shadow);
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }
  
  .tab-nav {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 15px 20px;
    border-bottom: 1px solid var(--border-color);
    background: var(--light-bg);
    border-radius: var(--border-radius) var(--border-radius) 0 0;
    flex-shrink: 0;
  }
  
  .tab-btn {
    padding: 10px 20px;
    border: none;
    background: none;
    cursor: pointer;
    border-radius: var(--border-radius);
    font-size: 15px;
    font-weight: 500;
    color: #666;
    transition: all 0.2s;
  }
  
  .tab-btn.active {
    background: var(--primary-color);
    color: white;
  }
  
  .tab-btn:not(.active):hover {
    background: #e9ecef;
  }
  
  .clear-btn {
    padding: 8px 16px;
    background: var(--danger-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
    margin-left: auto;
  }
  
  .clear-btn:hover {
    background: #c62828;
  }
  
  .tab-content {
    flex: 1;
    overflow: auto;
    padding: 20px;
  }
  
  .no-tasks {
    text-align: center;
    color: #666;
    font-style: italic;
    padding: 60px 20px;
    font-size: 16px;
  }
  
  .tasks-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  .task-item {
    background: var(--light-bg);
    border-radius: var(--border-radius);
    padding: 16px;
    border: 1px solid var(--border-color);
    transition: box-shadow 0.2s ease;
  }
  
  .task-item:hover {
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  }
  
  .task-item.history {
    opacity: 0.9;
  }
  
  .task-header {
    display: grid;
    grid-template-columns: 1fr auto auto;
    gap: 15px;
    align-items: center;
    margin-bottom: 10px;
  }
  
  .task-url {
    font-size: 14px;
    color: #666;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .task-status {
    font-size: 13px;
    font-weight: 500;
    padding: 5px 12px;
    border-radius: 15px;
    white-space: nowrap;
  }
  
  .task-status.downloading, .task-status.pending {
    background: var(--primary-color);
    color: white;
  }
  
  .task-status.completed {
    background: var(--success-color);
    color: white;
  }
  
  .task-status.failed, .task-status.cancelled {
    background: var(--danger-color);
    color: white;
  }
  
  .task-time {
    font-size: 13px;
    color: #666;
    white-space: nowrap;
  }
  
  .progress-container {
    margin-top: 12px;
  }
  
  .progress-bar {
    width: 100%;
    height: 8px;
    background: #e0e0e0;
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 8px;
  }
  
  .progress-fill {
    height: 100%;
    background: var(--primary-color);
    transition: width 0.3s ease;
  }
  
  .progress-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .progress-text {
    font-size: 13px;
    color: #666;
  }
  
  .cancel-btn {
    padding: 6px 12px;
    background: var(--danger-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
  }
  
  .cancel-btn:hover {
    background: #c62828;
  }
  
  .task-path, .task-error {
    margin-top: 10px;
    padding: 10px 12px;
    border-radius: var(--border-radius);
    font-size: 13px;
  }
  
  .task-path {
    background: #e8f5e9;
    color: #2e7d32;
    border-left: 3px solid #2e7d32;
  }
  
  .task-error {
    background: #ffebee;
    color: #c62828;
    border-left: 3px solid #c62828;
  }
  
  /* 响应式 */
  @media (max-width: 1200px) {
    .downloader-container {
      max-width: 98%;
      padding: 15px;
    }
  }
  
  @media (max-width: 768px) {
    .downloader-container {
      padding: 12px;
      gap: 15px;
    }
    
    .top-panel {
      padding: 15px;
    }
    
    .output-section {
      flex-direction: column;
      align-items: stretch;
      gap: 10px;
      margin-bottom: 15px;
      padding-bottom: 12px;
    }
    
    .output-section label {
      min-width: auto;
    }
    
    .task-header {
      grid-template-columns: 1fr;
      gap: 8px;
    }
    
    .task-status, .task-time {
      justify-self: start;
    }
    
    .form-row {
      flex-direction: column;
      gap: 12px;
    }
    
    .url-input, .download-btn {
      width: 100%;
    }
    
    .tab-nav {
      padding: 12px 15px;
    }
    
    .tab-content {
      padding: 15px;
    }
  }
</style>