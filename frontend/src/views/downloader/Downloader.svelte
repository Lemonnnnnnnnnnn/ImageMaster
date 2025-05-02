<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import {
      CancelDownload,
      ClearHistory,
      GetActiveTasks,
      GetHistoryTasks,
      StartDownload
  } from '../../../wailsjs/go/downloader/DownloaderAPI';
  import type { downloader } from '../../../wailsjs/go/models';
  import {
      GetOutputDir
  } from '../../../wailsjs/go/storage/StorageAPI';

  let url = '';
  let saveName = '';
  let loading = false;
  let outputDir = '';
  let error = '';
  
  // 下载任务列表
  let activeTasks: downloader.DownloadTask[] = [];
  let historyTasks: downloader.DownloadTask[] = [];
  let showHistory = false;
  
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
      
      // 如果显示历史，也获取历史任务
      if (showHistory) {
        historyTasks = await GetHistoryTasks();
      }
    } catch (err) {
      console.error("轮询任务状态出错:", err);
    }
  }

  async function handleSubmit() {
    if (!url) {
      error = '请输入网址';
      return;
    }

    try {
      loading = true;
      error = '';
      
      // 开始下载任务
      const taskId = await StartDownload(url, saveName);
      
      if (taskId) {
        // 重置表单
        url = '';
        saveName = '';
        
        // 立即刷新任务列表
        await pollTasks();
      } else {
        error = '下载失败，请检查网址是否正确';
      }
    } catch (err) {
      error = `下载出错: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
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
  
  // 切换显示历史记录
  async function toggleHistory() {
    showHistory = !showHistory;
    
    // 如果显示历史，获取历史任务
    if (showHistory) {
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

  function goToViewer() {
    push('/');
  }
  
  function goToConfig() {
    push('/config');
  }
</script>

<div class="downloader-container">
  <div class="header">
    <h1>网页图片下载器</h1>
    <div class="nav-buttons">
      <button on:click={goToViewer} class="nav-btn">返回漫画查看器</button>
      <button on:click={goToConfig} class="nav-btn config-btn">应用设置</button>
    </div>
  </div>

  <div class="settings-panel">
    <div class="output-dir">
      <h3>当前输出目录:</h3>
      <div class="dir-display">
        <span>{outputDir}</span>
        <button on:click={goToConfig}>前往设置更改</button>
      </div>
    </div>
  </div>
  
  <div class="download-panel">
    <h3>从网页下载图片</h3>
    
    <div class="form">
      <div class="form-group">
        <label for="url">网页地址</label>
        <input 
          type="text" 
          id="url" 
          bind:value={url} 
          placeholder="输入网页完整地址，例如: https://example.com/gallery"
          disabled={loading}
        />
      </div>
      
      <div class="form-group">
        <label for="saveName">保存文件夹名称 (可选)</label>
        <input 
          type="text" 
          id="saveName" 
          bind:value={saveName} 
          placeholder="自定义保存的文件夹名称，默认使用网页标题"
          disabled={loading}
        />
      </div>
      
      <button on:click={handleSubmit} disabled={loading} class="download-btn">
        {loading ? '添加中...' : '添加下载任务'}
      </button>
      
      {#if error}
        <div class="error">
          <p>{error}</p>
        </div>
      {/if}
    </div>
  </div>
  
  <!-- 下载任务列表 -->
  <div class="tasks-panel">
    <div class="tasks-header">
      <h3>下载任务</h3>
      <button on:click={toggleHistory} class="history-btn">
        {showHistory ? '隐藏历史记录' : '显示历史记录'}
      </button>
      {#if showHistory && historyTasks.length > 0}
        <button on:click={handleClearHistory} class="clear-btn">
          清除历史记录
        </button>
      {/if}
    </div>
    
    <!-- 活跃任务列表 -->
    {#if activeTasks.length > 0}
      <div class="active-tasks">
        <h4>当前任务 ({activeTasks.length})</h4>
        <div class="tasks-list">
          {#each activeTasks as task (task.id)}
            <div class="task-item">
              <div class="task-info">
                <div class="task-name-url">
                  <div class="task-name">{task.name || '未命名任务'}</div>
                  <div class="task-url">{task.url}</div>
                </div>
                <div class="task-status">{formatStatus(task.status)}</div>
                <div class="task-time">
                  开始于: {formatTime(task.startTime)}
                </div>
              </div>
              
              {#if task.status === 'downloading' || task.status === 'pending'}
                <!-- 进度条 -->
                <div class="progress-container">
                  <div class="progress-bar">
                    <div class="progress-fill" style="width: {task.progress.total > 0 ? Math.round((task.progress.current / task.progress.total) * 100) : 0}%"></div>
                  </div>
                  <div class="progress-text">
                    {#if task.progress.total > 0}
                      已下载: {task.progress.current}/{task.progress.total} 张图片 
                      ({Math.round((task.progress.current / task.progress.total) * 100)}%)
                    {:else}
                      准备下载中...
                    {/if}
                  </div>
                </div>
                
                <!-- 取消按钮 -->
                <button class="cancel-btn" on:click={() => cancelTask(task.id)}>
                  取消下载
                </button>
              {/if}
              
              {#if task.status === 'failed'}
                <div class="task-error">
                  错误信息: {task.error || '未知错误'}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {:else}
      <div class="no-tasks">
        当前没有进行中的下载任务
      </div>
    {/if}
    
    <!-- 历史任务列表 -->
    {#if showHistory && historyTasks.length > 0}
      <div class="history-tasks">
        <h4>历史记录 ({historyTasks.length})</h4>
        <div class="tasks-list">
          {#each historyTasks as task (task.id)}
            <div class="task-item history">
              <div class="task-info">
                <div class="task-name-url">
                  <div class="task-name">{task.name || '未命名任务'}</div>
                  <div class="task-url">{task.url}</div>
                </div>
                <div class="task-status {task.status}">{formatStatus(task.status)}</div>
                <div class="task-time">
                  <div>开始于: {formatTime(task.startTime)}</div>
                  <div>结束于: {formatTime(task.completeTime)}</div>
                </div>
              </div>
              
              {#if task.status === 'completed' && task.savePath}
                <div class="task-path">
                  保存路径: {task.savePath}
                </div>
              {/if}
              
              {#if task.status === 'failed' && task.error}
                <div class="task-error">
                  错误信息: {task.error}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
  
  <div class="instructions">
    <h3>使用说明</h3>
    <ul>
      <li>输入包含图片的网页地址</li>
      <li>可以指定自定义保存文件夹名称</li>
      <li>点击"添加下载任务"按钮开始抓取并下载网页中的图片</li>
      <li>可以同时添加多个下载任务，系统会并行处理</li>
      <li>可以随时取消正在进行的下载任务</li>
      <li>查看历史记录以了解已完成或失败的任务</li>
      <li>下载完成后，可在漫画查看器中浏览已下载的图片</li>
      <li>可以在"应用设置"中配置代理服务器和输出目录</li>
    </ul>
  </div>
</div>

<style>
  .downloader-container {
    padding: 20px;
    max-width: 1000px;
    margin: 0 auto;
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  
  .nav-buttons {
    display: flex;
    gap: 10px;
  }
  
  .nav-btn {
    padding: 8px 16px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .config-btn {
    background-color: #3498db;
  }
  
  .settings-panel, .download-panel, .tasks-panel, .instructions {
    background-color: #f5f5f5;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
  }
  
  .output-dir {
    margin-bottom: 10px;
  }
  
  .dir-display {
    display: flex;
    align-items: center;
    gap: 10px;
    background-color: white;
    padding: 10px;
    border-radius: 4px;
    border: 1px solid #ddd;
  }
  
  .dir-display span {
    flex: 1;
    word-break: break-all;
  }
  
  .dir-display button {
    padding: 5px 10px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .form {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }
  
  .form-group label {
    font-weight: bold;
  }
  
  .form-group input {
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
  }
  
  .download-btn {
    padding: 12px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    font-weight: bold;
  }
  
  .download-btn:disabled {
    background-color: #a5a5a5;
    cursor: not-allowed;
  }
  
  .error {
    background-color: #ffebee;
    color: #c62828;
    padding: 10px;
    border-radius: 4px;
    border-left: 4px solid #c62828;
  }
  
  /* 任务列表样式 */
  .tasks-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 15px;
  }
  
  .history-btn, .clear-btn {
    padding: 6px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .history-btn {
    background-color: #3498db;
    color: white;
  }
  
  .clear-btn {
    background-color: #e74c3c;
    color: white;
  }
  
  .no-tasks {
    padding: 20px;
    text-align: center;
    background-color: #f9f9f9;
    border-radius: 4px;
    border: 1px dashed #ccc;
  }
  
  .active-tasks, .history-tasks {
    margin-bottom: 20px;
  }
  
  .tasks-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  .task-item {
    background-color: white;
    border-radius: 4px;
    padding: 15px;
    border: 1px solid #ddd;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }
  
  .task-item.history {
    opacity: 0.8;
  }
  
  .task-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .task-name-url {
    flex: 2;
  }
  
  .task-name {
    font-weight: bold;
    margin-bottom: 5px;
  }
  
  .task-url {
    font-size: 12px;
    color: #666;
    word-break: break-all;
  }
  
  .task-status {
    flex: 1;
    text-align: center;
    font-weight: bold;
  }
  
  .task-status.completed {
    color: #2e7d32;
  }
  
  .task-status.failed, .task-status.cancelled {
    color: #c62828;
  }
  
  .task-time {
    flex: 1;
    font-size: 12px;
    color: #666;
  }
  
  .task-path, .task-error {
    margin-top: 10px;
    padding: 10px;
    border-radius: 4px;
    font-size: 14px;
  }
  
  .task-path {
    background-color: #e8f5e9;
    color: #2e7d32;
    border-left: 4px solid #2e7d32;
  }
  
  .task-error {
    background-color: #ffebee;
    color: #c62828;
    border-left: 4px solid #c62828;
  }
  
  .progress-container {
    margin: 10px 0;
  }
  
  .progress-bar {
    width: 100%;
    height: 20px;
    background-color: #e0e0e0;
    border-radius: 10px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background-color: #4CAF50;
    transition: width 0.3s ease;
  }
  
  .progress-text {
    text-align: center;
    margin-top: 5px;
    font-weight: bold;
  }
  
  .cancel-btn {
    padding: 8px 16px;
    background-color: #e74c3c;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-top: 10px;
  }
  
  .instructions ul {
    padding-left: 20px;
    line-height: 1.6;
  }
</style> 