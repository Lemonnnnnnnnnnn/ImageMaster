<script lang="ts">
  import { SelectLibrary, GetLibraries } from '../../../../wailsjs/go/viewer/Viewer';
  
  export let libraries = [];
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadLibraries() {
    loading = true;
    libraries = await GetLibraries();
    loading = false;
  }

  async function addLibrary() {
    loading = true;
    
    try {
      const newLib = await SelectLibrary();
      if (newLib) {
        await loadLibraries();
        onSuccess('成功添加新的漫画库');
      }
    } catch (err) {
      onError(`添加漫画库失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadLibraries };
</script>

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

<style>
  .libraries-card {
    grid-area: libraries;
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
</style>