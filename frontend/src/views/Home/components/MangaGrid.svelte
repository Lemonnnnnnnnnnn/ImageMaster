<script lang="ts">
  import { Grid, List } from 'lucide-svelte';
  import { mangas } from '../stores/homeStore';
  import MangaCard from './MangaCard.svelte';
  
  // 视图模式
  let viewMode: 'grid' | 'list' = 'grid';
  
  // 响应式网格配置
  $: gridClass = viewMode === 'grid' 
    ? 'grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8 gap-6'
    : 'flex flex-col gap-4';
</script>

<!-- Fluent Design 漫画网格 -->
<div class="space-y-6">
  <!-- 网格控制栏 -->
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-3">
      <h3 class="text-lg font-medium text-white-primary">
        漫画收藏
      </h3>
      <span class="px-2 py-1 bg-fluent-blue/20 text-fluent-blue text-xs font-medium rounded-fluent-sm">
        {$mangas.length} 部
      </span>
    </div>
    
    <!-- 视图切换按钮 -->
    <div class="flex items-center bg-glass-card rounded-fluent-lg p-1 border border-white-tertiary/20">
      <button
        class="p-2 rounded-fluent-md transition-fluent {viewMode === 'grid' ? 'bg-fluent-blue text-white' : 'text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10'}"
        on:click={() => viewMode = 'grid'}
        aria-label="网格视图"
        title="网格视图"
      >
        <Grid size={16} />
      </button>
      <button
        class="p-2 rounded-fluent-md transition-fluent {viewMode === 'list' ? 'bg-fluent-blue text-white' : 'text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10'}"
        on:click={() => viewMode = 'list'}
        aria-label="列表视图"
        title="列表视图"
      >
        <List size={16} />
      </button>
    </div>
  </div>
  
  <!-- 漫画网格/列表 -->
  <div class={gridClass}>
    {#each $mangas as manga (manga.path)}
      <MangaCard {manga} />
    {/each}
  </div>
  
  <!-- 底部统计信息 -->
  {#if $mangas.length > 0}
    <div class="text-center pt-8 border-t border-white-tertiary/20">
      <p class="text-sm text-white-secondary">
        共计 <span class="font-medium text-white-primary">{$mangas.length}</span> 部漫画
        · 
        <span class="text-fluent-blue">持续更新中</span>
      </p>
    </div>
  {/if}
</div>