<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { 
    SelectLibrary, 
    GetLibraries, 
    SetOutputDir,
    GetOutputDir,
    SetProxy,
    GetProxy
  } from '../../../wailsjs/go/viewer/Viewer';
  import Header from '../../components/Header.svelte';
  
  let libraries = [];
  let outputDir = '';
  let proxyURL = '';
  let loading = false;
  let error = '';
  let success = '';
  let successTimeout;

  onMount(async () => {
    await loadLibraries();
    await loadOutputDir();
    await loadProxySettings();
  });

  async function loadLibraries() {
    loading = true;
    libraries = await GetLibraries();
    loading = false;
  }

  async function loadOutputDir() {
    outputDir = await GetOutputDir();
  }

  async function loadProxySettings() {
    try {
      proxyURL = await GetProxy();
    } catch (err) {
      console.error('无法加载代理设置:', err);
      proxyURL = '';
    }
  }

  async function addLibrary() {
    loading = true;
    error = '';
    success = '';
    
    try {
      const newLib = await SelectLibrary();
      if (newLib) {
        await loadLibraries();
        showSuccessMessage('成功添加新的漫画库');
      }
    } catch (err) {
      error = `添加漫画库失败: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
  }

  async function changeOutputDir() {
    loading = true;
    error = '';
    success = '';
    
    try {
      const newDir = await SetOutputDir();
      if (newDir) {
        outputDir = newDir;
        showSuccessMessage('成功更改输出目录');
      }
    } catch (err) {
      error = `更改输出目录失败: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
  }

  async function saveProxySettings() {
    loading = true;
    error = '';
    success = '';
    
    try {
      await SetProxy(proxyURL);
      showSuccessMessage('成功保存代理设置');
    } catch (err) {
      error = `保存代理设置失败: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
  }

  function showSuccessMessage(message) {
    success = message;
    if (successTimeout) clearTimeout(successTimeout);
    successTimeout = setTimeout(() => {
      success = '';
    }, 3000);
  }
</script>

<div class="config-container">
  <Header title="应用设置" />

  <div class="notification-area">
    {#if loading}
      <div class="loading-indicator">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    {/if}

    {#if error}
      <div class="notification error">
        <div class="icon">✕</div>
        <p>{error}</p>
      </div>
    {/if}

    {#if success}
      <div class="notification success">
        <div class="icon">✓</div>
        <p>{success}</p>
      </div>
    {/if}
  </div>

  <div class="layout-grid">
    <!-- 漫画库设置 -->
    <div class="card libraries-card">
      <div class="card-header">
        <span class="card-icon">📚</span>
        <h2>漫画库设置</h2>
      </div>
      <div class="card-content">
        <div class="libraries-container">
          {#if libraries.length === 0}
            <div class="empty-state">
              <span class="empty-icon">📁</span>
              <p>当前未添加任何漫画库</p>
            </div>
          {:else}
            <div class="libraries-list-container">
              <h3>已添加的漫画库：</h3>
              <ul class="libraries-list">
                {#each libraries as lib}
                  <li>
                    <span class="folder-icon">📂</span>
                    <span class="lib-path">{lib}</span>
                  </li>
                {/each}
              </ul>
            </div>
          {/if}
        </div>
        <div class="action-row">
          <button on:click={addLibrary} disabled={loading} class="action-btn">
            <span class="btn-icon">+</span>
            <span>添加漫画库</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 输出目录设置 -->
    <div class="card output-card">
      <div class="card-header">
        <span class="card-icon">📂</span>
        <h2>输出目录</h2>
      </div>
      <div class="card-content">
        <div class="output-container">
          <div class="setting-value">
            <span class="folder-icon">📂</span>
            <span>{outputDir || '未设置'}</span>
          </div>
          <button on:click={changeOutputDir} disabled={loading} class="action-btn">
            <span class="btn-icon">📂</span>
            <span>更改输出目录</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 代理设置 -->
    <div class="card proxy-card">
      <div class="card-header">
        <span class="card-icon">⚙️</span>
        <h2>代理设置</h2>
      </div>
      <div class="card-content">
        <div class="proxy-container">
          <div class="input-row">
            <div class="input-wrapper">
              <label for="proxy">代理服务器 URL</label>
              <input 
                type="text" 
                id="proxy" 
                bind:value={proxyURL} 
                placeholder="例如: http://127.0.0.1:7890"
                disabled={loading}
              />
            </div>
            <button on:click={saveProxySettings} disabled={loading} class="action-btn">
              <span class="btn-icon">💾</span>
              <span>保存</span>
            </button>
          </div>
          <p class="hint">支持 HTTP 和 SOCKS 代理，格式为 http://host:port 或 socks5://host:port</p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .config-container {
    padding: 16px;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .notification-area {
    margin-bottom: 16px;
  }
  
  .layout-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: 16px;
    grid-template-areas: 
      "libraries output"
      "libraries proxy";
  }
  
  .libraries-card {
    grid-area: libraries;
  }
  
  .output-card {
    grid-area: output;
  }
  
  .proxy-card {
    grid-area: proxy;
  }
  
  .card {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .card-header {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    background-color: #f8f9fa;
    border-bottom: 1px solid #eee;
  }
  
  .card-icon {
    font-size: 20px;
    margin-right: 10px;
  }
  
  .card-header h2 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #333;
  }
  
  .card-content {
    padding: 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  
  .libraries-container {
    flex: 1;
    min-height: 180px;
    overflow-y: auto;
  }
  
  .libraries-list-container {
    height: 100%;
  }
  
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 0;
    color: #888;
  }
  
  .empty-icon {
    font-size: 32px;
    margin-bottom: 10px;
    opacity: 0.7;
  }
  
  .action-row {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
  }
  
  .input-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
  }
  
  .input-wrapper {
    flex: 1;
  }
  
  .setting-value {
    background-color: #f8f9fa;
    padding: 10px 12px;
    border-radius: 6px;
    border: 1px solid #eaeaea;
    word-break: break-all;
    display: flex;
    align-items: center;
    margin-bottom: 10px;
  }
  
  .folder-icon {
    margin-right: 8px;
    flex-shrink: 0;
  }
  
  .input-wrapper {
    margin-bottom: 8px;
  }
  
  .input-wrapper label {
    display: block;
    margin-bottom: 6px;
    font-size: 13px;
    color: #555;
  }
  
  .input-wrapper input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    transition: border-color 0.2s, box-shadow 0.2s;
  }
  
  .input-wrapper input:focus {
    border-color: #4a6fa5;
    box-shadow: 0 0 0 3px rgba(74, 111, 165, 0.15);
    outline: none;
  }
  
  .input-wrapper input::placeholder {
    color: #aaa;
  }
  
  .hint {
    font-size: 12px;
    color: #888;
    margin: 5px 0;
    line-height: 1.4;
  }
  
  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10px 16px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background-color 0.2s, transform 0.1s;
    white-space: nowrap;
  }
  
  .action-btn:hover {
    background-color: #3e5d8a;
  }
  
  .action-btn:active {
    transform: scale(0.98);
  }
  
  .action-btn:disabled {
    background-color: #b0bec5;
    cursor: not-allowed;
    transform: none;
  }
  
  .btn-icon {
    margin-right: 6px;
  }
  
  .libraries-list {
    margin: 8px 0 0 0;
    padding: 0;
    list-style-type: none;
    max-height: 200px;
    overflow-y: auto;
  }
  
  .libraries-list li {
    display: flex;
    align-items: center;
    padding: 8px 10px;
    background-color: #f8f9fa;
    border-radius: 4px;
    margin-bottom: 6px;
    font-size: 14px;
  }
  
  .lib-path {
    word-break: break-all;
  }
  
  .loading-indicator {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    background-color: #e3f2fd;
    border-radius: 6px;
    margin-bottom: 12px;
    animation: fadeIn 0.3s;
  }
  
  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(74, 111, 165, 0.2);
    border-left-color: #4a6fa5;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  .notification {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    border-radius: 6px;
    margin-bottom: 12px;
    animation: slideIn 0.3s;
  }
  
  .notification .icon {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    margin-right: 10px;
    font-weight: bold;
    font-size: 12px;
  }
  
  .notification.error {
    background-color: #ffebee;
    color: #c62828;
    border-left: 3px solid #c62828;
  }
  
  .notification.error .icon {
    background-color: #c62828;
    color: white;
  }
  
  .notification.success {
    background-color: #e8f5e9;
    color: #2e7d32;
    border-left: 3px solid #2e7d32;
  }
  
  .notification.success .icon {
    background-color: #2e7d32;
    color: white;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  @keyframes slideIn {
    from { transform: translateY(-10px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
  
  /* 响应式布局调整 */
  @media (max-width: 768px) {
    .layout-grid {
      grid-template-columns: 1fr;
      grid-template-areas: 
        "libraries"
        "output"
        "proxy";
    }
  }
  
  @media (min-width: 1200px) {
    .layout-grid {
      grid-template-columns: 1fr 1fr 1fr;
      grid-template-areas: "libraries output proxy";
    }
  }
</style> 