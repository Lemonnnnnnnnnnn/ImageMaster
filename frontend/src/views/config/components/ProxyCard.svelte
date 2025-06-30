<script lang="ts">
  import { SetProxy, GetProxy } from '../../../../wailsjs/go/viewer/Viewer';
  
  export let proxyURL = '';
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadProxySettings() {
    try {
      proxyURL = await GetProxy();
    } catch (err) {
      console.error('无法加载代理设置:', err);
      proxyURL = '';
    }
  }

  async function saveProxySettings() {
    loading = true;
    
    try {
      await SetProxy(proxyURL);
      onSuccess('成功保存代理设置');
    } catch (err) {
      onError(`保存代理设置失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadProxySettings };
</script>

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

<style>
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
  
  .input-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
  }
  
  .input-wrapper {
    flex: 1;
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
</style>