<script lang="ts">
  import { push } from 'svelte-spa-router';
  import { onMount } from 'svelte';

  // 删除activeView参数，只保留标题
  export let title: string = '图像大师';
  
  // 导航函数
  function goToViewer() {
    push('/');
  }
  
  function goToDownloader() {
    push('/downloader');
  }
  
  function goToConfig() {
    push('/config');
  }

  // 控制顶部导航栏的显示和隐藏
  let isVisible = true;
  let lastScrollY = 0;
  let headerElement;

  // 监听滚动事件
  function handleScroll() {
    const currentScrollY = window.scrollY;
    
    // 向上滚动时显示，向下滚动时隐藏
    if (currentScrollY < lastScrollY) {
      // 向上滚动
      isVisible = true;
    } else if (currentScrollY > lastScrollY && currentScrollY > 60) {
      // 向下滚动且滚动位置超过一定值
      isVisible = false;
    }
    
    lastScrollY = currentScrollY;
  }

  onMount(() => {
    // 添加滚动事件监听
    window.addEventListener('scroll', handleScroll, { passive: true });
    
    // 组件销毁时移除监听
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  });
</script>

<div class="header-wrapper" class:visible={isVisible} bind:this={headerElement}>
  <div class="header">
    <h1>{title}</h1>
    <div class="nav-buttons">
      <!-- 始终显示所有按钮 -->
      <button on:click={goToViewer} class="nav-btn viewer-btn">漫画查看器</button>
      <button on:click={goToDownloader} class="nav-btn downloader-btn">图片下载器</button> 
      <button on:click={goToConfig} class="nav-btn config-btn">应用设置</button>
    </div>
  </div>
</div>

<!-- 提供占位空间，防止内容被固定定位的header覆盖 -->
<div class="header-placeholder"></div>

<style>
  .header-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 1000;
    transform: translateY(0); /* 默认显示 */
    transition: transform 0.3s ease;
    background-color: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .header-wrapper:not(.visible) {
    transform: translateY(-100%); /* 隐藏时向上移动 */
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 20px;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .header-placeholder {
    height: 65px; /* 与header高度相匹配 */
  }
  
  h1 {
    margin: 0;
    font-size: 1.5rem;
    color: #333;
  }
  
  .nav-buttons {
    display: flex;
    gap: 10px;
  }
  
  .nav-btn {
    padding: 8px 16px;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.2s ease;
  }
  
  .viewer-btn {
    background-color: #9C27B0;
  }
  
  .downloader-btn {
    background-color: #4CAF50;
  }
  
  .config-btn {
    background-color: #3498db;
  }
  
  .nav-btn:hover {
    opacity: 0.9;
    transform: translateY(-2px);
  }
  
  @media (max-width: 600px) {
    .header {
      flex-direction: column;
      gap: 10px;
    }
    
    .header-placeholder {
      height: 100px;
    }
  }
</style> 