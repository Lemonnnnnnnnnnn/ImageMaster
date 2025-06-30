<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';

  $: ({ mangaName, selectedImages, showNavigation } = $mangaStore);
</script>

<div class="px-5 py-2 bg-gray-700 flex justify-between items-center">
  <div class="flex gap-2 min-w-60">
    <button on:click={MangaService.backToHome} class="px-3 py-2 bg-blue-600 text-white border-none rounded cursor-pointer disabled:bg-gray-600 disabled:cursor-not-allowed">返回</button>
  </div>
  <div class="text-center flex-grow">
    <h2 class="m-0 text-lg overflow-hidden text-ellipsis">{mangaName} ({selectedImages.length})</h2>
  </div>
  <div class="flex gap-2 min-w-60">
    <button on:click={MangaService.toggleNavigation} class="px-3 py-2 bg-blue-600 text-white border-none rounded cursor-pointer disabled:bg-gray-600 disabled:cursor-not-allowed">
      {showNavigation ? '隐藏导航' : '显示导航'}
    </button>
    <button on:click={MangaService.deleteAndViewNextManga} class="px-3 py-2 bg-red-700 text-white border-none rounded cursor-pointer disabled:bg-gray-600 disabled:cursor-not-allowed">删除并看下一部</button>
  </div>
</div>