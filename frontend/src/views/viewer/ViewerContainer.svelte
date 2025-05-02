<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { GetAllMangas } from '../../../wailsjs/go/viewer/Viewer';
  import Home from './Home.svelte';

  // 用于在 Viewer 和 Downloader 之间切换
  function goToDownloader() {
    push('/downloader');
  }
  
  // 用于前往配置页面
  function goToConfig() {
    push('/config');
  }
</script>

<div class="viewer-container">
  <div class="nav-header">
    <div class="nav-buttons">
      <button on:click={goToDownloader} class="nav-btn downloader-btn">前往下载器</button>
      <button on:click={goToConfig} class="nav-btn config-btn">应用设置</button>
    </div>
  </div>
  
  <slot>
    <!-- 默认显示 Home 组件 -->
    <Home />
  </slot>
</div>

<style>
  .viewer-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  
  .nav-header {
    display: flex;
    justify-content: flex-end;
    padding: 10px 20px;
    background-color: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
  }
  
  .nav-buttons {
    display: flex;
    gap: 10px;
  }
  
  .nav-btn {
    padding: 8px 16px;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .downloader-btn {
    background-color: #4CAF50;
  }
  
  .config-btn {
    background-color: #3498db;
  }
  
  .nav-btn:hover {
    opacity: 0.9;
  }
</style> 