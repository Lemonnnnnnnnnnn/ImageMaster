<script lang="ts">
  import Router from 'svelte-spa-router';
  import routes from './routes';
  import { Toaster, toast } from 'svelte-sonner';
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';

  onMount(() => {
    // 监听下载完成事件
    EventsOn('download:completed', (data: any) => {
      toast.success(`下载完成！`, {
        description: `已成功下载 ${data.name}`,
        duration: 5000,
      });
    });

    // 监听下载失败事件
    EventsOn('download:failed', (data: any) => {
      toast.error(`下载失败`, {
        description: `下载 ${data.name} 失败：${data.message || '下载过程中发生错误'}`,
        duration: 5000,
      });
    });

    // 监听下载取消事件
    EventsOn('download:cancelled', (data: any) => {
      toast.warning(`下载已取消`, {
        description: `已取消下载任务：${data.name}`,
        duration: 3000,
      });
    });
  });
</script>

<main>
  <Router {routes} />
</main>
<Toaster />

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #f5f5f5;
    color: #333;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  main {
    height: 100vh;
    width: 100%;
    /* overflow: hidden; */
    display: flex;
    flex-direction: column;
  }
</style>
