<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';

  $: ({ selectedImages, currentImageIndex } = $mangaStore);
</script>

<div class="single-view">
  <img src={selectedImages[currentImageIndex]} alt="Manga page {currentImageIndex + 1}" />
  
  <div class="single-view-controls">
    <button on:click={MangaService.prevImage} disabled={currentImageIndex === 0}>上一页</button>
    <span>{currentImageIndex + 1} / {selectedImages.length}</span>
    <button on:click={MangaService.nextImage} disabled={currentImageIndex === selectedImages.length - 1}>下一页</button>
  </div>
</div>

<style>
  .single-view {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    position: relative;
  }
  
  .single-view img {
    max-width: 100%;
    max-height: calc(100vh - 120px);
    object-fit: contain;
  }
  
  .single-view-controls {
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 20px;
    align-items: center;
    padding: 10px;
    background-color: rgba(42, 42, 42, 0.8);
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
  
  span {
    color: #f0f0f0;
  }
</style>