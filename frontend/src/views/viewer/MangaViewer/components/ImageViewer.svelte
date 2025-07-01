<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';
  import { MangaService } from '../services/mangaService';
  import LoadingSpinner from './LoadingSpinner.svelte';
  import ScrollView from './ScrollView.svelte';
  import Button from '../../../../components/Button.svelte';

  $: ({ loading, selectedImages } = $mangaStore);
</script>

{#if loading}
  <LoadingSpinner />
{:else if selectedImages.length === 0}
  <div class="flex-grow flex flex-col items-center justify-center h-full">
    <p class="text-gray-100 mb-5">未找到图片</p>
    <Button on:click={MangaService.backToHome} variant="filled" color="primary">返回</Button>
  </div>
{:else}
  <ScrollView />
{/if}