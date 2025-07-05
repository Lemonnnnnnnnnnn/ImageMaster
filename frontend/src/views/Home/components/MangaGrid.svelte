<script lang="ts">
    // import { BookOpen } from "lucide-svelte";
    import { MangaCard } from ".";
    import { mangas } from "../stores/homeStore";
</script>

<div class="h-screen overflow-auto">
    <!-- <div class="flex gap-3 items-center p-8">
        <span class="font-bold text-white text-xl">Comics</span>
        <span class=" bg-blue-500/20 text-blue-500 text-xs p-1 flex items-center gap-1">
            <BookOpen size={12} />
            {$mangas.length}
        </span>
    </div> -->
    <div class="flex flex-wrap gap-4 p-8">
        {#each $mangas as manga}
            <MangaCard {manga} />
        {/each}
    </div>
</div>
