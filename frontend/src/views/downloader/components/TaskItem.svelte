<script lang="ts">
  import type { models } from '../../../../wailsjs/go/models';
  import { formatStatus, formatTime, formatProgress, calculateProgressPercentage, getStatusClass } from '../utils/taskFormatter';
  import { cancelTask } from '../stores/downloadStore';
  import Button from '../../../components/Button.svelte';

  export let task: models.DownloadTask;
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
  <!-- 任务名称和状态行 -->
  <div class="grid grid-cols-[1fr_auto_auto] gap-4 items-center mb-2.5 md:grid-cols-1 md:gap-2">
    <div class="flex flex-col gap-1">
      <!-- 显示任务名称，如果没有名称则显示URL -->
      <div class="text-sm font-medium text-gray-800 overflow-hidden text-ellipsis whitespace-nowrap">
        {task.name || '未命名任务'}
      </div>
      <!-- URL作为副标题显示 -->
      <div class="text-xs text-gray-500 overflow-hidden text-ellipsis whitespace-nowrap">
        {task.url}
      </div>
    </div>
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
        <Button 
          variant="filled"
          color="error"
          size="sm"
          onclick={handleCancel}
        >
          取消
        </Button>
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