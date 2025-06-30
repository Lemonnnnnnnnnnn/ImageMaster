<script lang="ts">
  import { activeTasks } from '../stores/downloadStore';
  import TaskItem from './TaskItem.svelte';
</script>

{#if $activeTasks.length > 0}
  <div class="flex flex-col gap-4">
    {#each $activeTasks as task (task.id)}
      <TaskItem {task} showProgress={true} />
    {/each}
  </div>
{:else}
  <div class="text-center text-gray-600 italic py-15 px-5 text-base">
    当前没有进行中的下载任务
  </div>
{/if}