<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';
  import Button from '../../../../components/Button.svelte';

  $: ({ mangaName, selectedImages, showNavigation } = $mangaStore);
</script>

<div class="px-5 py-2 bg-gray-700 flex justify-between items-center">
  <div class="flex gap-2 min-w-60">
    <Button onclick={MangaService.backToHome} variant="filled" color="primary">返回</Button>
  </div>
  <div class="text-center flex-grow">
    <h2 class="m-0 text-lg overflow-hidden text-ellipsis">{mangaName} ({selectedImages.length})</h2>
  </div>
  <div class="flex gap-2 min-w-60">
    <Button onclick={MangaService.toggleNavigation} variant="filled" color="primary">
      {showNavigation ? '隐藏导航' : '显示导航'}
    </Button>
    <Button onclick={MangaService.deleteAndViewNextManga} variant="filled" color="error">删除并看下一部</Button>
  </div>
</div>