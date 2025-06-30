<script lang="ts">
  export let loading = false;
  export let error = '';
  export let success = '';
</script>

<div class="mb-4">
  {#if loading}
    <div class="flex items-center gap-2.5 px-3 py-2.5 bg-blue-50 rounded-md mb-3 animate-in fade-in duration-300">
      <div class="w-4 h-4 border-2 border-blue-200 border-l-blue-600 rounded-full animate-spin"></div>
      <p class="m-0">加载中...</p>
    </div>
  {/if}

  {#if error}
    <div class="flex items-center px-3 py-2.5 rounded-md mb-3 bg-red-50 text-red-700 border-l-4 border-red-700 animate-in slide-in-from-top-2 duration-300">
      <div class="w-5 h-5 flex items-center justify-center rounded-full mr-2.5 font-bold text-xs bg-red-700 text-white">✕</div>
      <p class="m-0">{error}</p>
    </div>
  {/if}

  {#if success}
    <div class="flex items-center px-3 py-2.5 rounded-md mb-3 bg-green-50 text-green-700 border-l-4 border-green-700 animate-in slide-in-from-top-2 duration-300">
      <div class="w-5 h-5 flex items-center justify-center rounded-full mr-2.5 font-bold text-xs bg-green-700 text-white">✓</div>
      <p class="m-0">{success}</p>
    </div>
  {/if}
</div>