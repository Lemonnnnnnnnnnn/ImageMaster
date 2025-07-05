<script lang="ts">
  import { push } from 'svelte-spa-router';
  import { Book, Calendar, Image, Play } from 'lucide-svelte';
  import type { Manga } from '../stores/homeStore';
  
  export let manga: Manga;
  
  // 导航到漫画查看页面
  function viewManga() {
    const encodedPath = encodeURIComponent(manga.path);
    push(`/manga/${encodedPath}`);
  }
  
  // 处理图片加载错误
  function handleImageError(event: Event) {
    const img = event.target as HTMLImageElement;
    img.style.display = 'none';
  }
  
  // 格式化更新时间
  function formatDate(dateString: string) {
    if (!dateString) return '未知';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch {
      return '未知';
    }
  }
</script>

<!-- Fluent Design 漫画卡片 -->
<div class="group fluent-card hover-lift cursor-pointer transition-fluent overflow-hidden" on:click={viewManga} on:keydown={(e) => e.key === 'Enter' && viewManga()} role="button" tabindex="0">
  <!-- 封面图片容器 -->
  <div class="relative aspect-[3/4] bg-gradient-to-br from-black-tertiary to-black-quaternary overflow-hidden">
    <!-- 封面图片 -->
    {#if manga.previewImg}
      <img
        src={manga.previewImg}
        alt={manga.name}
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        on:error={handleImageError}
        loading="lazy"
      />
    {:else}
      <!-- 默认封面 -->
      <div class="w-full h-full flex items-center justify-center">
        <Book size={48} class="text-white-tertiary" />
      </div>
    {/if}
    
    <!-- 悬停遮罩 -->
    <div class="absolute inset-0 bg-black/60 backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
      <div class="bg-fluent-blue rounded-full p-3 shadow-fluent-lg">
        <Play size={24} class="text-white ml-1" />
      </div>
    </div>
    
    <!-- 章节数量徽章 -->
    {#if manga.imagesCount && manga.imagesCount > 0}
      <div class="absolute top-3 right-3 bg-glass-overlay backdrop-blur-fluent-sm rounded-fluent-md px-2 py-1 border border-white-tertiary/20">
        <span class="text-xs font-medium text-white-primary flex items-center gap-1">
          <Image size={12} />
          {manga.imagesCount}
        </span>
      </div>
    {/if}
  </div>
  
  <!-- 卡片内容 -->
  <div class="p-4 space-y-3">
    <!-- 标题 -->
    <h3 class="font-medium text-white-primary line-clamp-2 leading-tight group-hover:text-fluent-blue transition-colors">
      {manga.name || '未知标题'}
    </h3>
    
    <!-- 元信息 -->
    <div class="space-y-2">
      <!-- 路径信息（仅显示文件夹名） -->
      {#if manga.path}
        <div class="text-xs text-white-tertiary truncate font-mono">
          {manga.path.split('/').pop() || manga.path.split('\\').pop() || manga.path}
        </div>
      {/if}
    </div>
    
    <!-- 文件数量信息 -->
    {#if manga.imagesCount}
      <div class="text-xs text-white-secondary">
        共 {manga.imagesCount} 张图片
      </div>
    {/if}
  </div>
</div>