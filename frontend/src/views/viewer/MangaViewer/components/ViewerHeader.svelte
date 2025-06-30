<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';

  $: ({ mangaName, selectedImages, viewMode, showNavigation } = $mangaStore);
</script>

<div class="viewer-header">
  <div class="left-controls">
    <button on:click={MangaService.backToHome}>返回</button>
    <button on:click={MangaService.toggleViewMode}>
      {viewMode === 'single' ? '切换为滚动模式' : '切换为单图模式'}
    </button>
  </div>
  <div class="center-title">
    <h2>{mangaName} ({selectedImages.length})</h2>
  </div>
  <div class="right-controls">
    <button on:click={MangaService.toggleNavigation}>
      {showNavigation ? '隐藏导航' : '显示导航'}
    </button>
    <button class="delete-btn" on:click={MangaService.deleteAndViewNextManga}>删除并看下一部</button>
  </div>
</div>

<style>
  .viewer-header {
    padding: 10px 20px;
    background-color: #2a2a2a;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .center-title {
    text-align: center;
    flex-grow: 1;
  }
  
  .center-title h2 {
    margin: 0;
    font-size: 18px;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .left-controls, .right-controls {
    display: flex;
    gap: 10px;
    min-width: 240px;
  }
  
  button {
    padding: 8px 12px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:disabled {
    background-color: #555;
    cursor: not-allowed;
  }
  
  .delete-btn {
    background-color: #c62828;
  }
</style>