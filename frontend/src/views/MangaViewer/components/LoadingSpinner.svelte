<script lang="ts">
  import { Loader2, Image } from 'lucide-svelte';
</script>

<!-- Fluent Design 漫画查看器加载动画 -->
<div class="flex flex-col items-center justify-center h-full py-20">
  <!-- 主加载动画 -->
  <div class="relative mb-8">
    <!-- 外层光环 -->
    <div class="absolute inset-0 w-16 h-16 bg-fluent-blue/20 rounded-full animate-ping"></div>
    
    <!-- 旋转容器 -->
    <div class="relative w-16 h-16 bg-glass-card backdrop-blur-fluent rounded-full flex items-center justify-center border border-white-tertiary/20 shadow-fluent-md">
      <Loader2 size={32} class="text-fluent-blue animate-spin" />
    </div>
  </div>
  
  <!-- 加载文本 -->
  <div class="text-center max-w-sm">
    <h3 class="text-white-primary font-medium mb-2 flex items-center justify-center gap-2">
      <Image size={18} class="text-fluent-blue" />
      正在加载漫画
    </h3>
    <p class="text-sm text-white-secondary">
      正在读取图片文件，请稍候...
    </p>
  </div>
  
  <!-- 进度指示器 -->
  <div class="mt-6 w-64 h-1 bg-white-tertiary/20 rounded-full overflow-hidden">
    <div class="h-full bg-gradient-to-r from-fluent-blue to-fluent-blue/80 rounded-full animate-pulse"></div>
  </div>
</div>

<style>
  .loading {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.1);
    border-left-color: #4a6fa5;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>