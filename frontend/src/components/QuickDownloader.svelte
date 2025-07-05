<script lang="ts">
  import { createDownloadHandler, formatDownloadError } from '../utils/downloadUtils';
  import Button from './Button.svelte';
  import { Download, X } from 'lucide-svelte';
  
  // 回调props - 替代createEventDispatcher
  export let ondownloadstarted: ((data: { taskId: string, url: string }) => void) | undefined = undefined;
  
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
      ondownloadstarted?.({ taskId, url: downloadUrl });
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
    // 延迟聚焦输入框
    setTimeout(() => {
      const input = document.getElementById('download-url');
      input?.focus();
    }, 100);
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
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleDownload();
    } else if (event.key === 'Escape') {
      closeModal();
    }
  }
  
  // 处理模态框点击外部关闭
  function handleModalClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      closeModal();
    }
  }
</script>

<!-- 悬浮图标容器 -->
<div 
  class="fixed bottom-6 right-6 z-[1000]"
  on:mouseenter={handleMouseEnter}
  on:mouseleave={handleMouseLeave}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && openModal()}
>
  <!-- 快捷下载图标 - Fluent Design 风格 -->
  <button
    class="w-14 h-14 bg-fluent-blue rounded-fluent-xl flex items-center justify-center text-white cursor-pointer shadow-fluent-lg transition-fluent backdrop-fluent border border-white-tertiary/20 {showIcon ? 'opacity-100 translate-y-0 pointer-events-auto' : 'opacity-0 translate-y-2 pointer-events-none'} hover:bg-fluent-blue/90 hover:scale-110 hover:shadow-fluent-xl focus:outline-none focus:ring-2 focus:ring-fluent-blue/50 focus:ring-offset-2 focus:ring-offset-black-secondary"
    on:click={openModal}
    aria-label="快速下载"
    title="快速下载"
  >
    <Download size={24} />
  </button>
</div>

<!-- 下载弹窗 - Fluent Design 风格 -->
{#if showModal}
  <div 
    class="fixed inset-0 bg-black/60 flex items-center justify-center z-[2000] backdrop-blur-fluent-md"
    on:click={handleModalClick}
    on:keydown={(e) => e.key === 'Escape' && closeModal()}
    role="dialog"
    aria-modal="true"
    aria-labelledby="modal-title"
  >
    <div 
      class="fluent-card w-[90%] max-w-lg shadow-fluent-xl animate-in slide-in-from-bottom-4 duration-fluent-medium"
      on:click|stopPropagation
    >
      <!-- 模态框头部 -->
      <div class="flex justify-between items-center px-6 py-5 border-b border-white-tertiary/20">
        <h3 id="modal-title" class="m-0 text-lg font-semibold text-white-primary">
          快速下载
        </h3>
        <button
          class="p-2 rounded-fluent-md text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent focus:outline-none focus:ring-2 focus:ring-fluent-blue/50"
          on:click={closeModal}
          aria-label="关闭"
        >
          <X size={20} />
        </button>
      </div>
      
      <!-- 模态框内容 -->
      <div class="px-6 py-5">
        <div class="flex flex-col gap-3">
          <label for="download-url" class="text-sm font-medium text-white-primary">
            网页地址:
          </label>
          <input 
            id="download-url"
            type="text" 
            bind:value={url}
            placeholder="输入网页完整地址，例如: https://example.com/gallery"
            disabled={loading}
            on:keydown={handleKeydown}
            class="fluent-input w-full px-4 py-3 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            autocomplete="url"
          />
        </div>
        
        <!-- 错误提示 -->
        {#if error}
          <div class="mt-4 px-4 py-3 bg-fluent-red/20 border border-fluent-red/40 rounded-fluent-md text-fluent-red text-sm backdrop-blur-fluent-sm">
            <div class="flex items-start gap-2">
              <svg class="w-4 h-4 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
              </svg>
              <span>{error}</span>
            </div>
          </div>
        {/if}
        
        <!-- 提示信息 -->
        <div class="mt-4 px-4 py-3 bg-fluent-blue/10 border border-fluent-blue/30 rounded-fluent-md text-fluent-blue text-sm backdrop-blur-fluent-sm">
          <div class="flex items-start gap-2">
            <svg class="w-4 h-4 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
            </svg>
            <span>支持 EHentai、Telegraph 等图片站点</span>
          </div>
        </div>
      </div>
      
      <!-- 模态框底部 -->
      <div class="flex justify-end gap-3 px-6 py-4 border-t border-white-tertiary/20">
        <Button 
          variant="outlined" 
          color="gray" 
          disabled={loading} 
          onclick={closeModal}
        >
          取消
        </Button>
        <Button 
          variant="filled" 
          color="primary" 
          disabled={loading || !url.trim()} 
          loading={loading} 
          onclick={handleDownload}
        >
          {loading ? '添加中...' : '开始下载'}
        </Button>
      </div>
    </div>
  </div>
{/if}