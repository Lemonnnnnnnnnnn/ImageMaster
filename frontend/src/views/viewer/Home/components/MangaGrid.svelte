<script lang="ts">
  import { mangas } from '../stores/homeStore';
  import MangaCard from './MangaCard.svelte';
</script>

<div class="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-5">
  {#each $mangas as manga}
    <MangaCard {manga} />
  {/each}
</div>