<script lang="ts">
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';
  import Button from '../../../../components/Button.svelte';

  $: ({ currentMangaIndex, mangas } = $mangaStore);
</script>

<div class="fixed bottom-5 right-5 flex flex-col gap-2 bg-gray-700 bg-opacity-80 p-2 rounded-lg">
  <Button on:click={MangaService.navigateToPrevManga} disabled={currentMangaIndex <= 0} variant="filled" color="primary">
    上一部漫画
  </Button>
  <Button on:click={MangaService.navigateToNextManga} disabled={currentMangaIndex >= mangas.length - 1} variant="filled" color="primary">
    下一部漫画
  </Button>
</div>