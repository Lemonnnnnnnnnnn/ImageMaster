<script lang="ts">
  import { SetOutputDir, GetOutputDir } from '../../../../wailsjs/go/viewer/Viewer';
  import Button from '../../../components/Button.svelte';
  import Card from '../../../components/Card.svelte';
  
  export let outputDir = '';
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadOutputDir() {
    outputDir = await GetOutputDir();
  }

  async function changeOutputDir() {
    loading = true;
    
    try {
      const newDir = await SetOutputDir();
      if (newDir) {
        outputDir = newDir;
        onSuccess('成功更改输出目录');
      }
    } catch (err : any) {
      onError(`更改输出目录失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadOutputDir };
</script>

<Card classes="p-4 preset-outlined-surface-500 overflow-hidden flex flex-col">
  <header class="card-header flex items-center">
    <span class="text-xl mr-2">📂</span>
    <h2 class="text-base font-semibold">输出目录</h2>
  </header>
  <section class="p-4 flex-1 flex flex-col">
    <div class="w-full">
      <div class="bg-surface-100 px-3 py-2 rounded border border-surface-300 break-all flex items-center mb-4">
        <span class="mr-2 flex-shrink-0">📂</span>
        <span>{outputDir || '未设置'}</span>
      </div>
      <Button 
        onclick={changeOutputDir} 
        disabled={loading}
        loading={loading}
        variant="filled"
        color="primary"
        classes="flex items-center"
      >
        <span class="mr-1">📂</span>
        <span>更改输出目录</span>
      </Button>
    </div>
  </section>
</Card>