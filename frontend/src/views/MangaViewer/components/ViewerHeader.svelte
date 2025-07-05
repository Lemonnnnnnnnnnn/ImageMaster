<script lang="ts">
  import { ArrowLeft, Navigation, Trash2, Eye, EyeOff, Image } from 'lucide-svelte';
  import { MangaService } from '../services/mangaService';
  import { mangaStore } from '../stores/mangaStore';
  import Button from '../../../components/Button.svelte';

  $: ({ mangaName, selectedImages, showNavigation } = $mangaStore);
</script>

<!-- Fluent Design 查看器头部 -->
<header class="sticky top-0 z-50 backdrop-fluent border-b border-white-tertiary/20 px-6 py-4">
  <div class="flex items-center justify-between">
    <!-- 左侧控制 -->
    <div class="flex items-center gap-3 flex-shrink-0 min-w-0">
      <Button 
        onclick={MangaService.backToHome} 
        variant="outlined" 
        color="gray"
        size="sm"
      >
        <ArrowLeft size={16} class="mr-2" />
        返回
      </Button>
    </div>
    
    <!-- 中间标题区域 -->
    <div class="flex-1 text-center px-6 min-w-0">
      <div class="flex items-center justify-center gap-3">
        <!-- 漫画图标 -->
        <div class="w-8 h-8 bg-fluent-blue/20 rounded-fluent-md flex items-center justify-center flex-shrink-0">
          <Image size={16} class="text-fluent-blue" />
        </div>
        
        <!-- 标题和统计 -->
        <div class="min-w-0">
          <h2 class="text-lg font-semibold text-white-primary m-0 truncate">
            {mangaName || '加载中...'}
          </h2>
          {#if selectedImages.length > 0}
            <p class="text-sm text-white-secondary m-0">
              共 {selectedImages.length} 张图片
            </p>
          {/if}
        </div>
      </div>
    </div>
    
    <!-- 右侧操作 -->
    <div class="flex items-center gap-3 flex-shrink-0">
      <!-- 导航切换按钮 -->
      <Button 
        onclick={MangaService.toggleNavigation} 
        variant={showNavigation ? "filled" : "outlined"}
        color="primary"
        size="sm"
      >
        {#if showNavigation}
          <EyeOff size={16} class="mr-2" />
          隐藏导航
        {:else}
          <Eye size={16} class="mr-2" />
          显示导航
        {/if}
      </Button>
      
      <!-- 删除按钮 -->
      <Button 
        onclick={MangaService.deleteAndViewNextManga} 
        variant="outlined" 
        color="error"
        size="sm"
      >
        <Trash2 size={16} class="mr-2" />
        删除并看下一部
      </Button>
    </div>
  </div>
  
  <!-- 进度指示器（可选） -->
  {#if selectedImages.length > 0}
    <div class="mt-3 w-full h-1 bg-white-tertiary/20 rounded-full overflow-hidden">
      <div 
        class="h-full bg-gradient-to-r from-fluent-blue to-fluent-blue/80 transition-all duration-300"
        style="width: 0%"
      ></div>
    </div>
  {/if}
</header>