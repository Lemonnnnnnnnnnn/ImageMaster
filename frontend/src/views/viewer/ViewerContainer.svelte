<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { GetAllMangas } from '../../../wailsjs/go/main/Viewer';
  import Home from './Home.svelte';

  // 用于在 Viewer 和 Downloader 之间切换
  function goToDownloader() {
    push('/downloader');
  }
</script>

<div class="viewer-container">
  <div class="nav-header">
    <button on:click={goToDownloader} class="switch-btn">切换到下载器</button>
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
  
  .switch-btn {
    padding: 8px 16px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .switch-btn:hover {
    background-color: #3a5f95;
  }
</style> 