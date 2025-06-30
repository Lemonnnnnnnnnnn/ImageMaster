<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';

  $: ({ selectedImages } = $mangaStore);
</script>

<div class="scroll-view">
  {#each selectedImages as image, i}
    <div class="scroll-image-container">
      <img src={image} alt="Manga page {i + 1}" />
    </div>
  {/each}
</div>

<style>
  .scroll-view {
    flex-grow: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
  
  .scroll-image-container {
    max-width: 800px;
    width: 100%;
  }
  
  .scroll-image-container img {
    width: 100%;
    height: auto;
    display: block;
    border-radius: 4px;
  }
</style>