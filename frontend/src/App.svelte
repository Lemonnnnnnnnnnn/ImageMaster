<script lang="ts">
  import Router from 'svelte-spa-router';
  import routes from './routes';
  import { Toaster, toast } from 'svelte-sonner';
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import MainLayout from './layout/MainLayout.svelte';

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

<MainLayout>
  <Router {routes} />
</MainLayout>
<Toaster 
  theme="dark" 
  position="top-right"
  richColors
  closeButton
/>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background-color: #1a1a1a;
    color: #ffffff;
    overflow-x: hidden;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  :global(html) {
    height: 100%;
    background-color: #1a1a1a;
  }

  /* 全局暗色主题 Toaster 样式 */
  :global([data-sonner-toaster]) {
    --normal-bg: rgba(42, 42, 42, 0.95) !important;
    --normal-border: rgba(110, 110, 110, 0.3) !important;
    --normal-text: #ffffff !important;
  }

  :global([data-sonner-toast]) {
    background: rgba(42, 42, 42, 0.95) !important;
    border: 1px solid rgba(110, 110, 110, 0.3) !important;
    color: #ffffff !important;
  }

  :global([data-sonner-toast][data-type="success"]) {
    --success-bg: rgba(16, 124, 16, 0.2) !important;
    --success-border: #107c10 !important;
  }

  :global([data-sonner-toast][data-type="error"]) {
    --error-bg: rgba(209, 52, 56, 0.2) !important;
    --error-border: #d13438 !important;
  }

  :global([data-sonner-toast][data-type="warning"]) {
    --warning-bg: rgba(255, 140, 0, 0.2) !important;
    --warning-border: #ff8c00 !important;
  }
</style>
