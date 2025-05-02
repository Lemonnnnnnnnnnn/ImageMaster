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
  
  let libraries = [];
  let outputDir = '';
  let proxyURL = '';
  let loading = false;
  let error = '';
  let success = '';

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
        success = '成功添加新的漫画库';
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
        success = '成功更改输出目录';
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
      success = '成功保存代理设置';
    } catch (err) {
      error = `保存代理设置失败: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
  }

  function goToHome() {
    push('/');
  }
  
  function goToDownloader() {
    push('/downloader');
  }
</script>

<div class="config-container">
  <div class="header">
    <h1>应用设置</h1>
    <div class="nav-buttons">
      <button on:click={goToHome} class="nav-btn">返回主页</button>
      <button on:click={goToDownloader} class="nav-btn">前往下载器</button>
    </div>
  </div>

  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>
  {/if}

  {#if error}
    <div class="error">
      <p>{error}</p>
    </div>
  {/if}

  {#if success}
    <div class="success">
      <p>{success}</p>
    </div>
  {/if}

  <div class="config-section">
    <h2>漫画库设置</h2>
    <div class="section-content">
      {#if libraries.length === 0}
        <p>当前未添加任何漫画库。</p>
      {:else}
        <h3>已添加的漫画库：</h3>
        <ul class="libraries-list">
          {#each libraries as lib}
            <li>{lib}</li>
          {/each}
        </ul>
      {/if}
      <button on:click={addLibrary} disabled={loading} class="action-btn">添加漫画库</button>
    </div>
  </div>

  <div class="config-section">
    <h2>下载设置</h2>
    <div class="section-content">
      <div class="setting-group">
        <h3>当前输出目录:</h3>
        <div class="setting-value">
          <span>{outputDir || '未设置'}</span>
        </div>
        <button on:click={changeOutputDir} disabled={loading} class="action-btn">更改输出目录</button>
      </div>

      <div class="setting-group">
        <h3>代理设置:</h3>
        <div class="input-group">
          <label for="proxy">代理服务器 URL (例如: http://127.0.0.1:7890)</label>
          <input 
            type="text" 
            id="proxy" 
            bind:value={proxyURL} 
            placeholder="代理服务器地址，留空表示不使用代理"
            disabled={loading}
          />
        </div>
        <p class="hint">支持 HTTP 和 SOCKS 代理，格式为 http://host:port 或 socks5://host:port</p>
        <button on:click={saveProxySettings} disabled={loading} class="action-btn">保存代理设置</button>
      </div>
    </div>
  </div>
</div>

<style>
  .config-container {
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
  
  .config-section {
    background-color: #f5f5f5;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
  }
  
  .section-content {
    margin-top: 15px;
  }
  
  .setting-group {
    margin-bottom: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid #ddd;
  }
  
  .setting-group:last-child {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 0;
  }
  
  .setting-value {
    background-color: white;
    padding: 10px;
    border-radius: 4px;
    border: 1px solid #ddd;
    margin: 10px 0;
    word-break: break-all;
  }
  
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
    margin-bottom: 10px;
  }
  
  .input-group input {
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
  }
  
  .hint {
    font-size: 14px;
    color: #666;
    margin: 5px 0 15px 0;
  }
  
  .action-btn {
    padding: 10px 15px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }
  
  .action-btn:disabled {
    background-color: #a5a5a5;
    cursor: not-allowed;
  }
  
  .libraries-list {
    margin: 10px 0;
    padding-left: 20px;
    line-height: 1.6;
  }
  
  .loading {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    margin-bottom: 15px;
  }
  
  .spinner {
    width: 24px;
    height: 24px;
    border: 3px solid rgba(0, 0, 0, 0.1);
    border-left-color: #4CAF50;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  .error {
    background-color: #ffebee;
    color: #c62828;
    padding: 10px;
    border-radius: 4px;
    border-left: 4px solid #c62828;
    margin-bottom: 15px;
  }
  
  .success {
    background-color: #e8f5e9;
    color: #2e7d32;
    padding: 10px;
    border-radius: 4px;
    border-left: 4px solid #2e7d32;
    margin-bottom: 15px;
  }
</style> 