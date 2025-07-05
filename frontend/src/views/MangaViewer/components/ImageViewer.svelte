<script lang="ts">
  import { ImageOff, ArrowLeft } from 'lucide-svelte';
  import { mangaStore } from '../stores/mangaStore';
  import { MangaService } from '../services/mangaService';
  import LoadingSpinner from './LoadingSpinner.svelte';
  import ScrollView from './ScrollView.svelte';
  import Button from '../../../components/Button.svelte';

  $: ({ loading, selectedImages } = $mangaStore);
</script>

<!-- Fluent Design 图片查看器 -->
<div class="flex-1 relative bg-black-primary">
  {#if loading}
    <!-- 加载状态 -->
    <LoadingSpinner />
  {:else if selectedImages.length === 0}
    <!-- 空状态 -->
    <div class="flex items-center justify-center h-full">
      <div class="text-center max-w-md px-6">
        <!-- 图标 -->
        <div class="w-24 h-24 bg-white-tertiary/10 rounded-fluent-xl flex items-center justify-center mb-8 mx-auto relative">
          <div class="absolute inset-0 bg-gradient-to-br from-fluent-blue/20 to-transparent rounded-fluent-xl"></div>
          <ImageOff size={48} class="text-white-tertiary relative z-10" />
        </div>
        
        <!-- 标题和描述 -->
        <div class="mb-8">
          <h2 class="text-2xl font-semibold text-white-primary mb-3">
            未找到图片
          </h2>
          <p class="text-white-secondary leading-relaxed">
            在当前漫画目录中没有找到可显示的图片文件
          </p>
        </div>
        
        <!-- 操作按钮 -->
        <Button 
          onclick={MangaService.backToHome} 
          variant="filled" 
          color="primary"
          classes="px-6 py-3"
        >
          <ArrowLeft size={18} class="mr-2" />
          返回首页
        </Button>
        
        <!-- 提示信息 -->
        <div class="mt-8">
          <div class="bg-glass-card/50 backdrop-blur-fluent-sm rounded-fluent-lg p-6 border border-white-tertiary/10">
            <p class="text-sm text-white-secondary">
              支持的图片格式：JPG、PNG、WEBP、GIF、BMP
            </p>
          </div>
        </div>
      </div>
    </div>
  {:else}
    <!-- 图片滚动视图 -->
    <ScrollView />
  {/if}
</div>