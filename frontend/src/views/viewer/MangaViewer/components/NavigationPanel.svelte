<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';

  $: ({ currentMangaIndex, mangas } = $mangaStore);
</script>

<div class="navigation-panel">
  <button on:click={MangaService.navigateToPrevManga} disabled={currentMangaIndex <= 0}>
    上一部漫画
  </button>
  <button on:click={MangaService.navigateToNextManga} disabled={currentMangaIndex >= mangas.length - 1}>
    下一部漫画
  </button>
</div>

<style>
  .navigation-panel {
    position: fixed;
    bottom: 20px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    background-color: rgba(42, 42, 42, 0.8);
    padding: 10px;
    border-radius: 8px;
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
</style>