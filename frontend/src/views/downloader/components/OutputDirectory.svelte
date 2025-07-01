<script lang="ts">
  import { push } from 'svelte-spa-router';
  import { outputDir } from '../stores/downloadStore';
</script>

<div class="flex items-center gap-4 mb-5 pb-4 border-b border-gray-300 md:flex-col md:items-stretch md:gap-2.5 md:mb-4 md:pb-3">
  <!-- <label class="font-medium text-gray-600 whitespace-nowrap min-w-[80px] md:min-w-auto">输出目录:</label> -->
  <div class="flex items-center gap-3 flex-1 bg-gray-50 px-4 py-3 rounded-md border border-gray-300">
    <span class="flex-1 text-sm text-gray-600 overflow-hidden text-ellipsis whitespace-nowrap">{$outputDir}</span>
    <button 
      on:click={() => push('/config')} 
      class="px-4 py-2 bg-blue-600 text-white border-none rounded-md cursor-pointer text-xs font-medium whitespace-nowrap hover:bg-purple-700"
    >
      设置
    </button>
  </div>
</div>