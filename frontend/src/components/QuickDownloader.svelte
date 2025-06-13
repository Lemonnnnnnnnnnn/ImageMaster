<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { createDownloadHandler, formatDownloadError } from '../utils/downloadUtils';
  
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
  function handleKeydown(event) {
    if (event.key === 'Enter') {
      handleDownload();
    } else if (event.key === 'Escape') {
      closeModal();
    }
  }
</script>

<!-- 悬浮图标容器 -->
<div 
  class="quick-downloader-container"
  on:mouseenter={handleMouseEnter}
  on:mouseleave={handleMouseLeave}
>
  <!-- 快捷下载图标 -->
  <div class="download-icon {showIcon ? 'visible' : ''}" on:click={openModal}>
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2V16M12 16L8 12M12 16L16 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M3 17V19C3 20.1046 3.89543 21 5 21H19C20.1046 21 21 20.1046 21 19V17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </div>
</div>

<!-- 下载弹窗 -->
{#if showModal}
  <div class="modal-overlay" on:click={closeModal}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <h3>快速下载</h3>
        <button class="close-btn" on:click={closeModal}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
      
      <div class="modal-body">
        <div class="input-group">
          <label for="download-url">网页地址:</label>
          <input 
            id="download-url"
            type="text" 
            bind:value={url}
            placeholder="输入网页完整地址，例如: https://example.com/gallery"
            disabled={loading}
            on:keydown={handleKeydown}
            class="url-input"
          />
        </div>
        
        {#if error}
          <div class="error-message">
            {error}
          </div>
        {/if}
      </div>
      
      <div class="modal-footer">
        <button class="cancel-btn" on:click={closeModal} disabled={loading}>
          取消
        </button>
        <button class="download-btn" on:click={handleDownload} disabled={loading}>
          {loading ? '添加中...' : '开始下载'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .quick-downloader-container {
    position: fixed;
    bottom: 20px;
    right: 20px;
    z-index: 1000;
  }
  
  .download-icon {
    width: 56px;
    height: 56px;
    background: #4361ee;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(67, 97, 238, 0.3);
    transition: all 0.3s ease;
    opacity: 0;
    transform: translateY(10px);
    pointer-events: none;
  }
  
  .download-icon.visible {
    opacity: 1;
    transform: translateY(0);
    pointer-events: auto;
  }
  
  .download-icon:hover {
    background: #3a0ca3;
    transform: scale(1.1);
    box-shadow: 0 6px 16px rgba(67, 97, 238, 0.4);
  }
  
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    backdrop-filter: blur(4px);
  }
  
  .modal-content {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
    animation: modalSlideIn 0.3s ease;
  }
  
  @keyframes modalSlideIn {
    from {
      opacity: 0;
      transform: translateY(-20px) scale(0.95);
    }
    to {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px 16px;
    border-bottom: 1px solid #e5e7eb;
  }
  
  .modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
  }
  
  .close-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    border-radius: 6px;
    color: #6b7280;
    transition: all 0.2s;
  }
  
  .close-btn:hover {
    background: #f3f4f6;
    color: #374151;
  }
  
  .modal-body {
    padding: 20px 24px;
  }
  
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .input-group label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }
  
  .url-input {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #d1d5db;
    border-radius: 8px;
    font-size: 14px;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }
  
  .url-input:focus {
    outline: none;
    border-color: #4361ee;
    box-shadow: 0 0 0 3px rgba(67, 97, 238, 0.1);
  }
  
  .url-input:disabled {
    background: #f9fafb;
    color: #6b7280;
  }
  
  .error-message {
    margin-top: 12px;
    padding: 12px 16px;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 8px;
    color: #dc2626;
    font-size: 14px;
  }
  
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 24px 20px;
    border-top: 1px solid #e5e7eb;
  }
  
  .cancel-btn, .download-btn {
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }
  
  .cancel-btn {
    background: #f3f4f6;
    color: #374151;
  }
  
  .cancel-btn:hover:not(:disabled) {
    background: #e5e7eb;
  }
  
  .download-btn {
    background: #4361ee;
    color: white;
  }
  
  .download-btn:hover:not(:disabled) {
    background: #3a0ca3;
  }
  
  .cancel-btn:disabled, .download-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>