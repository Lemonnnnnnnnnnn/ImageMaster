<script lang="ts">
  import type { Manga } from '../stores/homeStore';
  import { MangaService } from '../services/mangaService';
  import { NavigationService } from '../services/navigationService';

  export let manga: Manga;

  async function handleDelete(event: Event) {
    event.stopPropagation();
    await MangaService.deleteManga(manga);
  }

  function handleView() {
    NavigationService.viewManga(manga);
  }

  $: mangaImageSrc = MangaService.getMangaImage(manga.previewImg);
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  class="border border-gray-300 rounded-lg overflow-hidden transition-all duration-200 cursor-pointer hover:-translate-y-1 hover:shadow-lg"
  on:click={handleView}
>
  <div class="h-48 overflow-hidden">
    <img
      src={mangaImageSrc}
      alt={manga.name}
      class="w-full h-full object-cover"
    />
  </div>
  <div class="p-2">
    <h3
      class="m-0 mb-1 whitespace-nowrap overflow-hidden text-ellipsis"
    >
      {manga.name}
    </h3>
    <p class="m-0 mb-2 text-gray-600">{manga.imagesCount} 张图片</p>
    <button
      on:click={handleDelete}
      class="px-2 py-1 bg-red-600 text-white border-none rounded cursor-pointer"
    >
      删除
    </button>
  </div>
</div>