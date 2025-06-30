<script lang="ts">
  import { SetOutputDir, GetOutputDir } from '../../../../wailsjs/go/viewer/Viewer';
  
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

<div class="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
  <div class="flex items-center p-3 bg-gray-50 border-b border-gray-200">
    <span class="text-xl mr-2">📂</span>
    <h2 class="m-0 text-base font-semibold text-gray-800">输出目录</h2>
  </div>
  <div class="p-4 flex-1 flex flex-col">
    <div class="w-full">
      <div class="bg-gray-50 px-3 py-2 rounded border border-gray-200 break-all flex items-center mb-4">
        <span class="mr-2 flex-shrink-0">📂</span>
        <span>{outputDir || '未设置'}</span>
      </div>
      <button on:click={changeOutputDir} disabled={loading} class="flex items-center justify-center px-4 py-2 bg-blue-600 text-white border-none rounded cursor-pointer text-sm font-medium transition-all duration-200 whitespace-nowrap hover:bg-blue-700 active:scale-95 disabled:bg-gray-400 disabled:cursor-not-allowed disabled:transform-none">
        <span class="mr-1">📂</span>
        <span>更改输出目录</span>
      </button>
    </div>
  </div>
</div>