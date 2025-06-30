<script lang="ts">
  import { historyTasks } from '../stores/downloadStore';
  import TaskItem from './TaskItem.svelte';
</script>

{#if $historyTasks.length > 0}
  <div class="flex flex-col gap-4">
    {#each $historyTasks as task (task.id)}
      <TaskItem {task} showTimeRange={true} />
    {/each}
  </div>
{:else}
  <div class="text-center text-gray-600 italic py-15 px-5 text-base">
    暂无历史记录
  </div>
{/if}