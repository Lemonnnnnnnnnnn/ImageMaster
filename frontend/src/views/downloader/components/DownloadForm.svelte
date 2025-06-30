<script lang="ts">
  import { createDownloadHandler } from '../../../utils/downloadUtils';
  import { loading, activeTab, pollTasks } from '../stores/downloadStore';

  let url = '';
  let error = '';

  // 创建下载处理器
  const downloadHandler = createDownloadHandler({
    onStart: () => {
      $loading = true;
      error = '';
    },
    onSuccess: async (taskId, downloadUrl) => {
      // 重置表单
      url = '';
      
      // 切换到下载中标签并刷新任务列表
      $activeTab = 'downloading';
      await pollTasks();
    },
    onError: (errorMsg) => {
      error = errorMsg;
    },
    onFinally: () => {
      $loading = false;
    }
  });

  async function handleSubmit() {
    if (!url.trim()) {
      error = '请输入网址';
      return;
    }

    await downloadHandler(url.trim());
  }
</script>

<div class="flex flex-col gap-3">
  <div class="flex gap-4 items-center md:flex-col md:gap-3">
    <input 
      type="text" 
      bind:value={url} 
      placeholder="输入网页完整地址，例如: https://example.com/gallery"
      disabled={$loading}
      class="flex-1 px-4 py-3 border border-gray-300 rounded-md text-base focus:outline-none focus:border-blue-600 focus:ring-2 focus:ring-blue-200 md:w-full"
      on:keydown={(e) => e.key === 'Enter' && handleSubmit()}
    />
    <button 
      on:click={handleSubmit} 
      disabled={$loading} 
      class="px-6 py-3 bg-green-500 text-white border-none rounded-md cursor-pointer text-base font-medium whitespace-nowrap hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed md:w-full"
    >
      {$loading ? '添加中...' : '添加任务'}
    </button>
  </div>
  
  {#if error}
    <div class="bg-red-50 text-red-700 px-4 py-2.5 rounded-md border-l-4 border-red-700 text-sm">
      {error}
    </div>
  {/if}
</div>