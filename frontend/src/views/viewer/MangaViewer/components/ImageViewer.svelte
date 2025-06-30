<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';
  import { MangaService } from '../services/mangaService';
  import LoadingSpinner from './LoadingSpinner.svelte';
  import SingleView from './SingleView.svelte';
  import ScrollView from './ScrollView.svelte';

  $: ({ loading, selectedImages, viewMode } = $mangaStore);
</script>

{#if loading}
  <LoadingSpinner />
{:else if selectedImages.length === 0}
  <div class="no-images">
    <p>未找到图片</p>
    <button on:click={MangaService.backToHome}>返回</button>
  </div>
{:else if viewMode === 'single'}
  <SingleView />
{:else}
  <ScrollView />
{/if}

<style>
  .no-images {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
  }
  
  .no-images button {
    padding: 8px 12px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .no-images p {
    color: #f0f0f0;
    margin-bottom: 20px;
  }
</style>