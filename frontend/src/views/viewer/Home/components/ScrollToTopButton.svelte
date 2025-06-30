<script lang="ts">
  import { showScrollTop } from '../stores/homeStore';
  import { ScrollService } from '../services/scrollService';

  function handleScrollToTop() {
    ScrollService.scrollToTop();
  }
</script>

{#if $showScrollTop}
  <button
    on:click={handleScrollToTop}
    class="fixed bottom-8 right-8 w-12 h-12 rounded-full bg-blue-500 text-white text-2xl border-none cursor-pointer flex items-center justify-center shadow-lg"
  >
    ↑
  </button>
{/if}