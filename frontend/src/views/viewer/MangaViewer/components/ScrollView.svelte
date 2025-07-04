<script lang="ts">
  import { mangaStore } from '../stores/mangaStore';
  import { ProgressService } from '../services/progressService';
  import { onMount, onDestroy, beforeUpdate } from 'svelte';

  $: ({ selectedImages, mangaPath } = $mangaStore);
  
  let scrollContainer: HTMLElement;
  let saveTimeout: number;
  let isRestoringProgress = false;
  let lastMangaPath = '';
  let hasRestoredForCurrentManga = false;

  // 防抖保存进度
  function debounceSaveProgress() {
    if (saveTimeout) {
      clearTimeout(saveTimeout);
    }
    
    saveTimeout = setTimeout(() => {
      if (scrollContainer && mangaPath && !isRestoringProgress) {
        const scrollPosition = scrollContainer.scrollTop;
        ProgressService.saveProgress(mangaPath, scrollPosition, selectedImages.length);
      }
    }, 1000); // 1秒防抖
  }

  // 处理滚动事件
  function handleScroll() {
    debounceSaveProgress();
  }

  // 恢复滚动位置
  function restoreScrollPosition() {
    if (!scrollContainer || !mangaPath || hasRestoredForCurrentManga) return;
    
    const progress = ProgressService.getProgress(mangaPath);
    if (progress && progress.scrollPosition > 0) {
      isRestoringProgress = true;
      hasRestoredForCurrentManga = true;
      
      // 延迟恢复滚动位置，确保图片已加载
      setTimeout(() => {
        if (scrollContainer) {
          scrollContainer.scrollTop = progress.scrollPosition;
          console.log(`已恢复到上次阅读位置：${progress.scrollPosition}px`);
        }
        isRestoringProgress = false;
      }, 100);
    } else {
      hasRestoredForCurrentManga = true;
    }
  }

  // 监听漫画路径变化，只在新漫画时恢复位置
  $: if (mangaPath && mangaPath !== lastMangaPath) {
    lastMangaPath = mangaPath;
    hasRestoredForCurrentManga = false;
    
    // 延迟恢复，等待图片加载
    if (selectedImages.length > 0 && scrollContainer) {
      setTimeout(restoreScrollPosition, 200);
    }
  }

  // 监听图片加载完成，如果还没恢复位置则尝试恢复
  $: if (selectedImages.length > 0 && scrollContainer && !hasRestoredForCurrentManga) {
    setTimeout(restoreScrollPosition, 200);
  }

  onMount(() => {
    // 定期清理过期的进度记录
    ProgressService.cleanupOldProgress();
  });


  onDestroy(() => {
    // 组件销毁时清理定时器
    if (saveTimeout) {
      clearTimeout(saveTimeout);
    }
  });
</script>

<div 
  bind:this={scrollContainer}
  on:scroll={handleScroll}
  class="flex-grow overflow-y-auto p-5 flex flex-col items-center gap-5 scroll-smooth"
>
  {#each selectedImages as image, i}
    <div class="max-w-[800px] w-full">
      <img src={image} alt="Manga page {i + 1}" class="w-full h-auto block rounded" />
    </div>
  {/each}
</div>