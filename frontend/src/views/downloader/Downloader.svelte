<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { 
    CrawlFromWeb,
    SetOutputDir,
    GetOutputDir
  } from '../../../wailsjs/go/viewer/Viewer';

  let url = '';
  let saveName = '';
  let loading = false;
  let outputDir = '';
  let result = '';
  let error = '';

  onMount(async () => {
    // 获取当前输出目录
    outputDir = await GetOutputDir();
  });

  async function handleSubmit() {
    if (!url) {
      error = '请输入网址';
      return;
    }

    try {
      loading = true;
      error = '';
      result = '';
      
      // 开始下载
      const savedDir = await CrawlFromWeb(url, saveName);
      
      if (savedDir) {
        result = `已成功下载图片到: ${savedDir}`;
      } else {
        error = '下载失败，请检查网址是否正确';
      }
    } catch (err) {
      error = `下载出错: ${err.message || '未知错误'}`;
    } finally {
      loading = false;
    }
  }

  async function changeOutputDir() {
    const newDir = await SetOutputDir();
    if (newDir) {
      outputDir = newDir;
    }
  }

  function goToViewer() {
    push('/');
  }
</script>

<div class="downloader-container">
  <div class="header">
    <h1>网页图片下载器</h1>
    <button on:click={goToViewer} class="nav-btn">返回漫画查看器</button>
  </div>

  <div class="settings-panel">
    <div class="output-dir">
      <h3>当前输出目录:</h3>
      <div class="dir-display">
        <span>{outputDir}</span>
        <button on:click={changeOutputDir}>更改</button>
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
        {loading ? '下载中...' : '开始下载'}
      </button>
      
      {#if loading}
        <div class="loading">
          <div class="spinner"></div>
          <p>正在下载图片，请稍候...</p>
        </div>
      {/if}
      
      {#if error}
        <div class="error">
          <p>{error}</p>
        </div>
      {/if}
      
      {#if result}
        <div class="success">
          <p>{result}</p>
        </div>
      {/if}
    </div>
  </div>
  
  <div class="instructions">
    <h3>使用说明</h3>
    <ul>
      <li>输入包含图片的网页地址</li>
      <li>可以指定自定义保存文件夹名称</li>
      <li>点击"开始下载"按钮开始抓取并下载网页中的图片</li>
      <li>下载完成后，可在漫画查看器中浏览已下载的图片</li>
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
  
  .nav-btn {
    padding: 8px 16px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .settings-panel, .download-panel, .instructions {
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
  
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 15px;
    gap: 10px;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(0, 0, 0, 0.1);
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
  }
  
  .success {
    background-color: #e8f5e9;
    color: #2e7d32;
    padding: 10px;
    border-radius: 4px;
    border-left: 4px solid #2e7d32;
  }
  
  .instructions ul {
    padding-left: 20px;
    line-height: 1.6;
  }
</style> 