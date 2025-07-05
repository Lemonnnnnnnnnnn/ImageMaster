<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { loading, libraries, mangas } from './stores/homeStore';
  import { MangaService } from './services/mangaService';
  import { ScrollService } from './services/scrollService';
  import { Library, ImageIcon } from 'lucide-svelte';
  
  import Header from '../../components/Header.svelte';
  import LoadingSpinner from './components/LoadingSpinner.svelte';
  import EmptyState from './components/EmptyState.svelte';
  import MangaGrid from './components/MangaGrid.svelte';
  import ScrollToTopButton from './components/ScrollToTopButton.svelte';

  let cleanupScrollListener: (() => void) | null = null;

  onMount(async () => {
    // 初始化数据加载
    await MangaService.initialize();
    
    // 初始化滚动监听
    cleanupScrollListener = ScrollService.initScrollListener();
  });

  onDestroy(() => {
    // 清理滚动监听器
    if (cleanupScrollListener) {
      cleanupScrollListener();
    }
  });
</script>

<!-- Fluent Design 黑色主题漫画查看器主页 -->
<div class="min-h-screen bg-black-secondary">
  <Header 
    title="漫画查看器" 
    subtitle="浏览和查看您的漫画收藏"
    showActions={true}
  >
    <svelte:fragment slot="actions">
      <!-- 库统计信息 -->
      {#if $libraries.length > 0}
        <div class="flex items-center gap-2 px-3 py-1.5 bg-glass-card rounded-fluent-md border border-white-tertiary/20">
          <Library size={16} class="text-fluent-blue" />
          <span class="text-sm text-white-secondary">
            {$libraries.length} 个库 · {$mangas.length} 部漫画
          </span>
        </div>
      {/if}
    </svelte:fragment>
  </Header>

  <!-- 主要内容区域 -->
  <div class="container mx-auto max-w-7xl px-6 py-8">
    {#if $libraries.length === 0 && !$loading}
      <!-- 未设置漫画库状态 -->
      <EmptyState type="no-libraries" />
    {:else if $loading}
      <!-- 加载状态 -->
      <LoadingSpinner />
    {:else if $mangas.length > 0}
      <!-- 漫画网格 -->
      <MangaGrid />
    {:else if $libraries.length > 0 && !$loading}
      <!-- 无漫画状态 -->
      <EmptyState type="no-mangas" />
    {/if}
  </div>

  <!-- 回到顶部按钮 -->
  <ScrollToTopButton />
</div>