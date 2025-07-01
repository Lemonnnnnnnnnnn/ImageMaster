<script lang="ts">
  import { showScrollTop } from '../stores/homeStore';
  import { ScrollService } from '../services/scrollService';
  import Button from '../../../../components/Button.svelte';

  function handleScrollToTop() {
    ScrollService.scrollToTop();
  }
</script>

{#if $showScrollTop}
  <Button
    on:click={handleScrollToTop}
    classes="fixed bottom-8 right-8 w-12 h-12 rounded-full text-2xl shadow-lg"
    color="primary"
  >
    ↑
  </Button>
{/if}