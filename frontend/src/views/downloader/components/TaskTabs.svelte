<script lang="ts">
  import { activeTab, activeTasksCount, historyTasks, switchTab, clearHistory } from '../stores/downloadStore';
  import Button from '../../../components/Button.svelte';

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
  <Button 
    on:click={() => handleTabSwitch('downloading')}
    variant={$activeTab === 'downloading' ? 'filled' : 'ghost'}
    color="gray"
    size="md"
  >
    下载中 ({$activeTasksCount})
  </Button>
  <Button 
    on:click={() => handleTabSwitch('history')}
    variant={$activeTab === 'history' ? 'filled' : 'ghost'}
    color="secondary"
    size="md"
    classes="text-gray-600"
  >
    历史记录
  </Button>
  {#if $activeTab === 'history' && $historyTasks.length > 0}
    <Button 
      on:click={handleClearHistory}
      variant="filled"
      color="error"
      size="sm"
      classes="ml-auto"
    >
      清除历史
    </Button>
  {/if}
</div>