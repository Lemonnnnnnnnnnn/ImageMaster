<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';
  import { MangaService } from '../services/mangaService';
  import LoadingSpinner from './LoadingSpinner.svelte';
  import ScrollView from './ScrollView.svelte';

  $: ({ loading, selectedImages } = $mangaStore);
</script>

{#if loading}
  <LoadingSpinner />
{:else if selectedImages.length === 0}
  <div class="flex-grow flex flex-col items-center justify-center h-full">
    <p class="text-gray-100 mb-5">未找到图片</p>
    <button on:click={MangaService.backToHome} class="px-3 py-2 bg-blue-600 text-white border-none rounded cursor-pointer">返回</button>
  </div>
{:else}
  <ScrollView />
{/if}