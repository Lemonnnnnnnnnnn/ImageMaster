<script lang="ts">
  import { ChevronLeft, ChevronRight, BookOpen, X } from 'lucide-svelte';
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';
  import Button from '../../../components/Button.svelte';

  $: ({ currentMangaIndex, mangas } = $mangaStore);
  
  // 关闭导航面板
  function closeNavigation() {
    MangaService.toggleNavigation();
  }
  
  // 处理键盘事件
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      closeNavigation();
    }
  }
</script>

<!-- Fluent Design 导航面板 -->
<div 
  class="flex items-center justify-center h-full"
  on:keydown={handleKeyDown}
  role="dialog"
  aria-modal="true"
  aria-label="漫画导航"
  tabindex="-1"
>
  <!-- 导航卡片 -->
  <div class="fluent-card p-8 max-w-lg w-full mx-6 relative">
    <!-- 关闭按钮 -->
    <button
      class="absolute top-4 right-4 p-2 rounded-fluent-md text-white-secondary hover:text-white-primary hover:bg-white-tertiary/10 transition-fluent"
      on:click={closeNavigation}
      aria-label="关闭导航"
    >
      <X size={20} />
    </button>
    
    <!-- 标题 -->
    <div class="text-center mb-8">
      <div class="w-16 h-16 bg-fluent-blue/20 rounded-fluent-xl flex items-center justify-center mx-auto mb-4">
        <BookOpen size={32} class="text-fluent-blue" />
      </div>
      <h2 class="text-xl font-semibold text-white-primary mb-2">
        漫画导航
      </h2>
      <p class="text-sm text-white-secondary">
        {currentMangaIndex + 1} / {mangas.length} 部漫画
      </p>
    </div>
    
    <!-- 导航按钮 -->
    <div class="space-y-4">
      <!-- 上一部漫画 -->
      <Button 
        onclick={MangaService.navigateToPrevManga} 
        disabled={currentMangaIndex <= 0}
        variant="outlined" 
        color="primary"
        classes="w-full flex items-center justify-center gap-3 py-4"
      >
        <ChevronLeft size={20} />
        <span class="flex-1 text-center">上一部漫画</span>
        <div class="w-5"></div> <!-- 占位符保持居中 -->
      </Button>
      
      <!-- 下一部漫画 -->
      <Button 
        onclick={MangaService.navigateToNextManga} 
        disabled={currentMangaIndex >= mangas.length - 1}
        variant="filled" 
        color="primary"
        classes="w-full flex items-center justify-center gap-3 py-4"
      >
        <div class="w-5"></div> <!-- 占位符保持居中 -->
        <span class="flex-1 text-center">下一部漫画</span>
        <ChevronRight size={20} />
      </Button>
    </div>
    
    <!-- 进度指示器 -->
    <div class="mt-8">
      <div class="flex justify-between text-xs text-white-secondary mb-2">
        <span>进度</span>
        <span>{Math.round(((currentMangaIndex + 1) / mangas.length) * 100)}%</span>
      </div>
      <div class="w-full h-2 bg-white-tertiary/20 rounded-full overflow-hidden">
        <div 
          class="h-full bg-gradient-to-r from-fluent-blue to-fluent-blue/80 transition-all duration-500"
          style="width: {((currentMangaIndex + 1) / mangas.length) * 100}%"
        ></div>
      </div>
    </div>
    
    <!-- 提示信息 -->
    <div class="mt-6 p-4 bg-fluent-blue/10 border border-fluent-blue/30 rounded-fluent-md">
      <p class="text-sm text-fluent-blue text-center">
        使用 ← → 方向键或点击按钮来导航漫画
      </p>
    </div>
  </div>
</div>