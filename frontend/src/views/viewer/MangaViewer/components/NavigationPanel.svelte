<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';

  $: ({ currentMangaIndex, mangas } = $mangaStore);
</script>

<div class="fixed bottom-5 right-5 flex flex-col gap-2 bg-gray-700 bg-opacity-80 p-2 rounded-lg">
  <button on:click={MangaService.navigateToPrevManga} disabled={currentMangaIndex <= 0} class="px-3 py-2 bg-blue-600 text-white border-none rounded cursor-pointer disabled:bg-gray-600 disabled:cursor-not-allowed">
    上一部漫画
  </button>
  <button on:click={MangaService.navigateToNextManga} disabled={currentMangaIndex >= mangas.length - 1} class="px-3 py-2 bg-blue-600 text-white border-none rounded cursor-pointer disabled:bg-gray-600 disabled:cursor-not-allowed">
    下一部漫画
  </button>
</div>