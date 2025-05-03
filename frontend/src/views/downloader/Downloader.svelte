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
  import Header from '../../components/Header.svelte';

  let url = '';
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
      const taskId = await StartDownload(url);
      
      if (taskId) {
        // 重置表单
        url = '';
        
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
</script>

<div class="downloader-container">
  <Header title="网页图片下载器" />

  <div class="settings-panel">
    <div class="output-dir">
      <h3>当前输出目录:</h3>
      <div class="dir-display">
        <span>{outputDir}</span>
        <button on:click={()=>push('/config')}>前往设置更改</button>
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
  :root {
    --primary-color: #4361ee;
    --secondary-color: #3a0ca3;
    --success-color: #4CAF50;
    --danger-color: #e63946;
    --warning-color: #ff9f1c;
    --light-bg: #f8f9fa;
    --dark-bg: #212529;
    --border-color: #dee2e6;
    --box-shadow: 0 4px 6px rgba(0,0,0,0.1);
    --border-radius: 8px;
  }

  .downloader-container {
    padding: 20px;
    max-width: 1000px;
    margin: 0 auto;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    color: #333;
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    border-bottom: 2px solid var(--border-color);
    padding-bottom: 15px;
  }
  
  .header h1 {
    color: var(--primary-color);
    margin: 0;
    font-size: 28px;
  }
  
  .nav-buttons {
    display: flex;
    gap: 12px;
  }
  
  .nav-btn {
    padding: 10px 18px;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s ease;
    box-shadow: var(--box-shadow);
  }
  
  .nav-btn:hover {
    background-color: var(--secondary-color);
    transform: translateY(-2px);
  }
  
  .config-btn {
    background-color: var(--secondary-color);
  }
  
  .config-btn:hover {
    background-color: #2a0a73;
  }
  
  .settings-panel, .download-panel, .tasks-panel, .instructions {
    background-color: var(--light-bg);
    border-radius: var(--border-radius);
    padding: 25px;
    margin-bottom: 25px;
    box-shadow: var(--box-shadow);
    border: 1px solid var(--border-color);
  }
  
  .settings-panel h3, .download-panel h3, .tasks-panel h3, .instructions h3 {
    margin-top: 0;
    color: var(--primary-color);
    font-size: 20px;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 10px;
    margin-bottom: 20px;
  }
  
  .output-dir {
    margin-bottom: 10px;
  }
  
  .dir-display {
    display: flex;
    align-items: center;
    gap: 15px;
    background-color: white;
    padding: 15px;
    border-radius: var(--border-radius);
    border: 1px solid var(--border-color);
  }
  
  .dir-display span {
    flex: 1;
    word-break: break-all;
    color: #555;
  }
  
  .dir-display button {
    padding: 8px 15px;
    background-color: var(--success-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 500;
  }
  
  .dir-display button:hover {
    background-color: #3d8b40;
    transform: translateY(-2px);
  }
  
  .form {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .form-group label {
    font-weight: 600;
    color: #555;
  }
  
  .form-group input {
    padding: 12px 15px;
    border: 1px solid var(--border-color);
    border-radius: var(--border-radius);
    font-size: 16px;
    transition: border-color 0.2s;
  }
  
  .form-group input:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.2);
  }
  
  .download-btn {
    padding: 14px;
    background-color: var(--success-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 16px;
    font-weight: 600;
    transition: all 0.2s ease;
    box-shadow: var(--box-shadow);
  }
  
  .download-btn:hover:not(:disabled) {
    background-color: #3d8b40;
    transform: translateY(-2px);
  }
  
  .download-btn:disabled {
    background-color: #a5a5a5;
    cursor: not-allowed;
    box-shadow: none;
  }
  
  .error {
    background-color: #ffebee;
    color: #c62828;
    padding: 12px 15px;
    border-radius: var(--border-radius);
    border-left: 4px solid #c62828;
  }
  
  /* 任务列表样式 */
  .tasks-header {
    display: flex;
    align-items: center;
    gap: 15px;
    margin-bottom: 20px;
  }
  
  .tasks-header h3 {
    margin: 0;
    padding: 0;
    border: none;
  }
  
  .history-btn, .clear-btn {
    padding: 8px 15px;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.2s ease;
  }
  
  .history-btn {
    background-color: var(--primary-color);
    color: white;
  }
  
  .history-btn:hover {
    background-color: var(--secondary-color);
  }
  
  .clear-btn {
    background-color: var(--danger-color);
    color: white;
  }
  
  .clear-btn:hover {
    background-color: #c62828;
  }
  
  .no-tasks {
    padding: 30px;
    text-align: center;
    background-color: white;
    border-radius: var(--border-radius);
    border: 1px dashed var(--border-color);
    color: #666;
    font-style: italic;
  }
  
  .active-tasks, .history-tasks {
    margin-bottom: 30px;
  }
  
  .active-tasks h4, .history-tasks h4 {
    margin-top: 0;
    margin-bottom: 15px;
    color: var(--primary-color);
    font-size: 18px;
  }
  
  .tasks-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  
  .task-item {
    background-color: white;
    border-radius: var(--border-radius);
    padding: 20px;
    border: 1px solid var(--border-color);
    box-shadow: var(--box-shadow);
    transition: transform 0.2s ease;
  }
  
  .task-item:hover {
    transform: translateY(-2px);
  }
  
  .task-item.history {
    opacity: 0.85;
  }
  
  .task-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
    flex-wrap: wrap;
    gap: 15px;
  }
  
  .task-name-url {
    flex: 2;
  }
  
  .task-name {
    font-weight: bold;
    margin-bottom: 8px;
    color: #333;
  }
  
  .task-url {
    font-size: 13px;
    color: #666;
    word-break: break-all;
    padding: 5px 0;
  }
  
  .task-status {
    flex: 1;
    text-align: center;
    font-weight: bold;
    padding: 5px 10px;
    border-radius: 20px;
    align-self: flex-start;
  }
  
  .task-status.completed {
    color: white;
    background-color: var(--success-color);
  }
  
  .task-status.downloading, .task-status.pending {
    color: white;
    background-color: var(--primary-color);
  }
  
  .task-status.failed, .task-status.cancelled {
    color: white;
    background-color: var(--danger-color);
  }
  
  .task-time {
    flex: 1;
    font-size: 13px;
    color: #666;
    line-height: 1.6;
  }
  
  .task-path, .task-error {
    margin-top: 15px;
    padding: 12px 15px;
    border-radius: var(--border-radius);
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
    margin: 15px 0;
  }
  
  .progress-bar {
    width: 100%;
    height: 20px;
    background-color: #e0e0e0;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: inset 0 1px 3px rgba(0,0,0,0.1);
  }
  
  .progress-fill {
    height: 100%;
    background-color: var(--primary-color);
    transition: width 0.3s ease;
    background-image: linear-gradient(45deg, 
                      rgba(255,255,255,.15) 25%, transparent 25%, 
                      transparent 50%, rgba(255,255,255,.15) 50%, 
                      rgba(255,255,255,.15) 75%, transparent 75%, 
                      transparent);
    background-size: 40px 40px;
    animation: progress-animation 1s linear infinite;
  }
  
  @keyframes progress-animation {
    0% {
      background-position: 0 0;
    }
    100% {
      background-position: 40px 0;
    }
  }
  
  .progress-text {
    text-align: center;
    margin-top: 8px;
    font-weight: 600;
    color: #555;
  }
  
  .cancel-btn {
    padding: 10px 18px;
    background-color: var(--danger-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    margin-top: 15px;
    font-weight: 500;
    transition: all 0.2s ease;
  }
  
  .cancel-btn:hover {
    background-color: #c62828;
    transform: translateY(-2px);
  }
  
  .instructions ul {
    padding-left: 25px;
    line-height: 1.8;
    color: #555;
  }
  
  .instructions li {
    margin-bottom: 8px;
  }
  
  /* 响应式调整 */
  @media (max-width: 768px) {
    .header {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .nav-buttons {
      margin-top: 15px;
      width: 100%;
    }
    
    .nav-btn {
      flex: 1;
      text-align: center;
    }
    
    .task-info {
      flex-direction: column;
    }
    
    .task-status {
      align-self: flex-start;
    }
  }
</style> 