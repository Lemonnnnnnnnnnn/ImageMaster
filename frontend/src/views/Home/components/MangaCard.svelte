<script lang="ts">
    import { MangaService } from "../services";
    import type { Manga } from "../stores/homeStore";
    import { push } from "svelte-spa-router";

    export let manga: Manga;

    const mangaImageSrc = MangaService.getMangaImage(manga.previewImg);
</script>

<button
    class="rounded-lg overflow-hidden xl:w-64 w-32 bg-neutral-400 cursor-pointer text-left hover:translate-y-[-4px] transition-transform duration-300 relative"
    onclick={() => {
        push(`/manga/${manga.path}`);
    }}
>
    <div class="h-48 overflow-hidden">
        <img
            src={mangaImageSrc}
            alt={manga.name}
            class="w-full h-full object-cover"
        />
    </div>
    <div class="p-2 bg-neutral-600 flex flex-col gap-2 xl:h-24 h-20">
        <h3 class="xl:text-sm text-xs font-bold text-white line-clamp-2">
            {manga.name}
        </h3>
        <h3 class="xl:text-xs text-[10px] text-neutral-400">
            total {manga.imagesCount} images
        </h3>
    </div>
</button>
