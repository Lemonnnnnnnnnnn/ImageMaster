<script lang="ts">
  import type { downloader } from '../../../../wailsjs/go/models';
  import { formatStatus, formatTime, formatProgress, calculateProgressPercentage, getStatusClass } from '../utils/taskFormatter';
  import { cancelTask } from '../stores/downloadStore';

  export let task: downloader.DownloadTask;
  export let showProgress = false;
  export let showTimeRange = false;

  async function handleCancel() {
    try {
      await cancelTask(task.id);
    } catch (err) {
      console.error('取消任务失败:', err);
    }
  }
</script>

<div class="bg-gray-50 rounded-md p-4 border border-gray-300 transition-shadow duration-200 hover:shadow-md {showTimeRange ? 'opacity-90' : ''}">
  <div class="grid grid-cols-[1fr_auto_auto] gap-4 items-center mb-2.5 md:grid-cols-1 md:gap-2">
    <div class="text-sm text-gray-600 overflow-hidden text-ellipsis whitespace-nowrap">{task.url}</div>
    <div class="text-xs font-medium px-3 py-1.5 rounded-full whitespace-nowrap {getStatusClass(task.status)}">
      {formatStatus(task.status)}
    </div>
    <div class="text-xs text-gray-600 whitespace-nowrap md:justify-self-start">
      {#if showTimeRange}
        {formatTime(task.startTime)} - {formatTime(task.completeTime)}
      {:else}
        {formatTime(task.startTime)}
      {/if}
    </div>
  </div>
  
  {#if showProgress && (task.status === 'downloading' || task.status === 'pending')}
    <!-- 进度条 -->
    <div class="mt-3">
      <div class="w-full h-2 bg-gray-300 rounded overflow-hidden mb-2">
        <div 
          class="h-full bg-blue-600 transition-all duration-300" 
          style="width: {calculateProgressPercentage(task.progress.current, task.progress.total)}%"
        ></div>
      </div>
      <div class="flex justify-between items-center">
        <span class="text-xs text-gray-600">
          {formatProgress(task.progress.current, task.progress.total)}
        </span>
        <button 
          class="px-3 py-1.5 bg-red-600 text-white border-none rounded-md cursor-pointer text-xs font-medium hover:bg-red-700" 
          on:click={handleCancel}
        >
          取消
        </button>
      </div>
    </div>
  {/if}
  
  {#if task.status === 'failed'}
    <div class="mt-2.5 px-3 py-2.5 rounded-md text-xs bg-red-50 text-red-700 border-l-4 border-red-700">
      错误: {task.error || '未知错误'}
    </div>
  {/if}
  
  {#if task.status === 'completed' && task.savePath}
    <div class="mt-2.5 px-3 py-2.5 rounded-md text-xs bg-green-50 text-green-700 border-l-4 border-green-700">
      保存路径: {task.savePath}
    </div>
  {/if}
</div>