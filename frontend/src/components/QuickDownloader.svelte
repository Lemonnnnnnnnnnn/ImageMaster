<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { createDownloadHandler, formatDownloadError } from '../utils/downloadUtils';
  import Button from './Button.svelte';
  
  const dispatch = createEventDispatcher();
  
  let showModal = false;
  let url = '';
  let loading = false;
  let error = '';
  let showIcon = false;
  
  // 创建下载处理器
  const downloadHandler = createDownloadHandler({
    onStart: () => {
      loading = true;
      error = '';
    },
    onSuccess: (taskId, downloadUrl) => {
      // 下载任务创建成功
      dispatch('downloadStarted', { taskId, url: downloadUrl });
      closeModal();
    },
    onError: (errorMsg) => {
      error = errorMsg;
    },
    onFinally: () => {
      loading = false;
    }
  });
  
  // 显示/隐藏快捷图标
  function handleMouseEnter() {
    showIcon = true;
  }
  
  function handleMouseLeave() {
    showIcon = false;
  }
  
  // 打开下载弹窗
  function openModal() {
    showModal = true;
    url = '';
    error = '';
  }
  
  // 关闭下载弹窗
  function closeModal() {
    showModal = false;
    url = '';
    error = '';
    loading = false;
  }
  
  // 处理下载
  async function handleDownload() {
    if (!url.trim()) {
      error = '请输入网址';
      return;
    }
    
    await downloadHandler(url.trim());
  }
  
  // 处理键盘事件
  function handleKeydown(event: any) {
    if (event.key === 'Enter') {
      handleDownload();
    } else if (event.key === 'Escape') {
      closeModal();
    }
  }
</script>

<!-- 悬浮图标容器 -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div 
  class="fixed bottom-5 right-5 z-[1000]"
  on:mouseenter={handleMouseEnter}
  on:mouseleave={handleMouseLeave}
>
  <!-- 快捷下载图标 -->
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div class="w-14 h-14 bg-blue-600 rounded-full flex items-center justify-center text-white cursor-pointer shadow-lg transition-all duration-300 ease-in-out {showIcon ? 'opacity-100 translate-y-0 pointer-events-auto' : 'opacity-0 translate-y-2.5 pointer-events-none'} hover:bg-purple-700 hover:scale-110 hover:shadow-xl" on:click={openModal}>
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2V16M12 16L8 12M12 16L16 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M3 17V19C3 20.1046 3.89543 21 5 21H19C20.1046 21 21 20.1046 21 19V17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </div>
</div>

<!-- 下载弹窗 -->
{#if showModal}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-[2000] backdrop-blur-sm" on:click={closeModal}>
    <div class="bg-white rounded-xl w-[90%] max-w-lg shadow-2xl animate-in slide-in-from-bottom-4 duration-300" on:click|stopPropagation>
      <div class="flex justify-between items-center px-6 py-5 border-b border-gray-200">
        <h3 class="m-0 text-lg font-semibold text-gray-800">快速下载</h3>
        <Button variant="ghost" color="primary" size="sm" classes="p-1 rounded-md text-gray-500 hover:bg-gray-100 hover:text-gray-700" on:click={closeModal}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </Button>
      </div>
      
      <div class="px-6 py-5">
        <div class="flex flex-col gap-2">
          <label for="download-url" class="text-sm font-medium text-gray-700">网页地址:</label>
          <input 
            id="download-url"
            type="text" 
            bind:value={url}
            placeholder="输入网页完整地址，例如: https://example.com/gallery"
            disabled={loading}
            on:keydown={handleKeydown}
            class="w-full px-4 py-3 border text-black border-gray-300 rounded-lg text-sm transition-colors duration-200 box-border focus:outline-none focus:border-blue-600 focus:ring-2 focus:ring-blue-100 disabled:bg-gray-50 disabled:text-gray-500"
          />
        </div>
        
        {#if error}
          <div class="mt-3 px-4 py-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        {/if}
      </div>
      
      <div class="flex justify-end gap-3 px-6 py-4 border-t border-gray-200">
        <Button variant="filled" color="gray" disabled={loading} on:click={closeModal}>
          取消
        </Button>
        <Button variant="filled" color="primary" disabled={loading} loading={loading} on:click={handleDownload}>
          {loading ? '添加中...' : '开始下载'}
        </Button>
      </div>
    </div>
  </div>
{/if}