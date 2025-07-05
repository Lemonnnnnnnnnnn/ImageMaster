<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { ChevronUp } from 'lucide-svelte';
  
  let isVisible = false;
  let scrollY = 0;
  
  // 显示/隐藏按钮的阈值
  const threshold = 300;
  
  // 监听滚动事件
  function handleScroll() {
    scrollY = window.scrollY;
    isVisible = scrollY > threshold;
  }
  
  // 平滑滚动到顶部
  function scrollToTop() {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  }
  
  // 处理键盘事件
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      scrollToTop();
    }
  }
  
  onMount(() => {
    window.addEventListener('scroll', handleScroll, { passive: true });
    handleScroll(); // 初始检查
  });
  
  onDestroy(() => {
    window.removeEventListener('scroll', handleScroll);
  });
</script>

<!-- Fluent Design 回到顶部按钮 -->
{#if isVisible}
  <div class="fixed bottom-6 left-6 z-50 animate-in slide-in-from-bottom-4 duration-300">
    <button
      class="group w-12 h-12 bg-glass-overlay backdrop-blur-fluent rounded-fluent-xl border border-white-tertiary/20 shadow-fluent-lg hover:shadow-fluent-xl transition-all duration-fluent-normal hover:scale-105 focus:outline-none focus:ring-2 focus:ring-fluent-blue/50 focus:ring-offset-2 focus:ring-offset-black-secondary active:scale-95"
      on:click={scrollToTop}
      on:keydown={handleKeyDown}
      aria-label="回到顶部"
      title="回到顶部"
    >
      <!-- 背景光效 -->
      <div class="absolute inset-0 bg-gradient-to-br from-fluent-blue/20 to-transparent rounded-fluent-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
      
      <!-- 图标 -->
      <div class="relative flex items-center justify-center w-full h-full">
        <ChevronUp size={20} class="text-white-primary group-hover:text-fluent-blue transition-colors duration-300" />
      </div>
      
      <!-- 进度环（可选，显示滚动进度） -->
      <svg class="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 48 48">
        <circle
          cx="24"
          cy="24"
          r="20"
          fill="none"
          stroke="rgba(255, 255, 255, 0.1)"
          stroke-width="2"
        />
        <circle
          cx="24"
          cy="24"
          r="20"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-dasharray={2 * Math.PI * 20}
          stroke-dashoffset={2 * Math.PI * 20 * (1 - Math.min(scrollY / (document.documentElement.scrollHeight - window.innerHeight), 1))}
          stroke-linecap="round"
          class="text-fluent-blue opacity-60 transition-all duration-300"
        />
      </svg>
    </button>
  </div>
{/if}