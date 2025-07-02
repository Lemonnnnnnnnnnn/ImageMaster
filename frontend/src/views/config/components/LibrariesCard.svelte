<script lang="ts">
  import { SelectLibrary } from '../../../../wailsjs/go/library/API';
import { GetLibraries } from '../../../../wailsjs/go/storage/API';
  import Button from '../../../components/Button.svelte';
  import Card from '../../../components/Card.svelte';
  
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

<Card classes="p-4 preset-outlined-surface-500 overflow-hidden flex flex-col">
  <header class="card-header flex items-center">
    <span class="text-xl mr-2.5">📚</span>
    <h2 class="text-base font-semibold">漫画库设置</h2>
  </header>
  <section class="p-4 flex-1 flex flex-col">
    <div class="flex-1 min-h-[180px] overflow-y-auto">
      {#if libraries.length === 0}
        <div class="flex flex-col items-center py-5 text-surface-500">
          <span class="text-3xl mb-2.5 opacity-70">📁</span>
          <p class="m-0">当前未添加任何漫画库</p>
        </div>
      {:else}
        <div class="h-full">
          <h3 class="m-0 mb-2 text-sm font-medium">已添加的漫画库：</h3>
          <ul class="mt-2 mb-0 p-0 list-none max-h-[200px] overflow-y-auto">
            {#each libraries as lib}
              <li class="flex items-center px-2.5 py-2 bg-surface-100 rounded mb-1.5 text-sm">
                <span class="mr-2 flex-shrink-0">📂</span>
                <span class="break-all">{lib}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <footer class="card-footer flex justify-end">
      <Button 
        onclick={addLibrary} 
        disabled={loading}
        loading={loading}
        variant="outlined"
        color="secondary"
        classes="flex items-center"
      >
        <span class="mr-1.5">+</span>
        <span>添加漫画库</span>
      </Button>
    </footer>
  </section>
</Card>