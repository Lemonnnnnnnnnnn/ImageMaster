<script lang="ts">
  import type { Manga } from '../stores/homeStore';
  import { MangaService } from '../services/mangaService';
  import { NavigationService } from '../services/navigationService';
  import Button from '../../../../components/Button.svelte';
  import Card from '../../../../components/Card.svelte';

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
<Card
  classes="overflow-hidden transition-all duration-200 cursor-pointer hover:-translate-y-1 hover:shadow-lg"
  onclick={handleView}
>
  <div class="h-48 overflow-hidden">
    <img
      src={mangaImageSrc}
      alt={manga.name}
      class="w-full h-full object-cover"
    />
  </div>
  <div class="p-4">
    <h3
      class="card-header mb-2 whitespace-nowrap overflow-hidden text-ellipsis"
    >
      {manga.name}
    </h3>
    <p class="mb-4 text-surface-500">{manga.imagesCount} 张图片</p>
    <Button
      variant="filled"
      color="error"
      size="sm"
      onclick={handleDelete}
    >
      删除
    </Button>
  </div>
</Card>