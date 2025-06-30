<script lang="ts">
  import { SelectLibrary, GetLibraries } from '../../../../wailsjs/go/viewer/Viewer';
  
  export let libraries : string[] = [];
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadLibraries() {
    loading = true;
    libraries = await GetLibraries();
    loading = false;
  }

  async function addLibrary() {
    loading = true;
    
    try {
      const newLib = await SelectLibrary();
      if (newLib) {
        await loadLibraries();
        onSuccess('成功添加新的漫画库');
      }
    } catch (err : any) {
      onError(`添加漫画库失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadLibraries };
</script>

<div class="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
  <div class="flex items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
    <span class="text-xl mr-2.5">📚</span>
    <h2 class="m-0 text-base font-semibold text-gray-800">漫画库设置</h2>
  </div>
  <div class="p-4 flex-1 flex flex-col">
    <div class="flex-1 min-h-[180px] overflow-y-auto">
      {#if libraries.length === 0}
        <div class="flex flex-col items-center py-5 text-gray-500">
          <span class="text-3xl mb-2.5 opacity-70">📁</span>
          <p class="m-0">当前未添加任何漫画库</p>
        </div>
      {:else}
        <div class="h-full">
          <h3 class="m-0 mb-2 text-sm font-medium text-gray-700">已添加的漫画库：</h3>
          <ul class="mt-2 mb-0 p-0 list-none max-h-[200px] overflow-y-auto">
            {#each libraries as lib}
              <li class="flex items-center px-2.5 py-2 bg-gray-50 rounded mb-1.5 text-sm">
                <span class="mr-2 flex-shrink-0">📂</span>
                <span class="break-all">{lib}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="mt-2.5 flex justify-end">
      <button on:click={addLibrary} disabled={loading} class="flex items-center justify-center px-4 py-2.5 bg-blue-600 text-white border-none rounded-md cursor-pointer text-sm font-medium transition-all duration-200 whitespace-nowrap hover:bg-blue-700 active:scale-[0.98] disabled:bg-gray-400 disabled:cursor-not-allowed disabled:transform-none">
        <span class="mr-1.5">+</span>
        <span>添加漫画库</span>
      </button>
    </div>
  </div>
</div>