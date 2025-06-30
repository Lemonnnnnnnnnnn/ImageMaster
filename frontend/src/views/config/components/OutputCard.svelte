<script lang="ts">
  import { SetOutputDir, GetOutputDir } from '../../../../wailsjs/go/viewer/Viewer';
  
  export let outputDir = '';
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadOutputDir() {
    outputDir = await GetOutputDir();
  }

  async function changeOutputDir() {
    loading = true;
    
    try {
      const newDir = await SetOutputDir();
      if (newDir) {
        outputDir = newDir;
        onSuccess('成功更改输出目录');
      }
    } catch (err) {
      onError(`更改输出目录失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadOutputDir };
</script>

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

<style>
  .output-card {
    grid-area: output;
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