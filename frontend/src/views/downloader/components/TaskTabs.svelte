<script lang="ts">
  import { activeTab, activeTasksCount, historyTasks, switchTab, clearHistory } from '../stores/downloadStore';

  async function handleTabSwitch(tab: 'downloading' | 'history') {
    await switchTab(tab);
  }

  async function handleClearHistory() {
    try {
      await clearHistory();
    } catch (err) {
      console.error('清除历史失败:', err);
    }
  }
</script>

<div class="flex items-center gap-1.5 px-5 py-4 border-b border-gray-300 bg-gray-50 rounded-t-md flex-shrink-0">
  <button 
    class="px-5 py-2.5 border-none cursor-pointer rounded-md text-base font-medium text-gray-600 transition-all duration-200 {$activeTab === 'downloading' ? 'bg-gray-200 ' : 'hover:bg-gray-200'}"
    on:click={() => handleTabSwitch('downloading')}
  >
    下载中 ({$activeTasksCount})
  </button>
  <button 
    class="px-5 py-2.5 border-none cursor-pointer rounded-md text-base font-medium text-gray-600 transition-all duration-200 {$activeTab === 'history' ? 'bg-gray-200 ' : 'hover:bg-gray-200'}"
    on:click={() => handleTabSwitch('history')}
  >
    历史记录
  </button>
  {#if $activeTab === 'history' && $historyTasks.length > 0}
    <button 
      on:click={handleClearHistory} 
      class="px-4 py-2 bg-red-600 text-white border-none rounded-md cursor-pointer text-xs font-medium ml-auto hover:bg-red-700"
    >
      清除历史
    </button>
  {/if}
</div>