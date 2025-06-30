<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';

  $: ({ selectedImages } = $mangaStore);
</script>

<div class="flex-grow overflow-y-auto p-5 flex flex-col items-center gap-5">
  {#each selectedImages as image, i}
    <div class="max-w-[800px] w-full">
      <img src={image} alt="Manga page {i + 1}" class="w-full h-auto block rounded" />
    </div>
  {/each}
</div>