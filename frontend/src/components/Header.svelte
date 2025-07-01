<script lang="ts">
  import { push } from 'svelte-spa-router';
  // import { QuickDownloaderService } from '../views/downloader/services/quickDownloaderService';
  import { onMount } from 'svelte';
  import Button from './Button.svelte';

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

  // function openQuickDownloader() {
  //   QuickDownloaderService.openModal();
  // }

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

<div class="fixed top-0 left-0 right-0 z-[1000] transition-transform duration-300 ease-in-out bg-gray-100 border-b border-gray-300 shadow-md" class:translate-y-0={isVisible} class:-translate-y-full={!isVisible} bind:this={headerElement}>
  <div class="flex justify-between items-center px-5 py-2.5 max-w-6xl mx-auto sm:gap-2.5">
    <h1 class="m-0 text-2xl text-gray-800">{title}</h1>
    <div class="flex gap-2.5">
      <!-- 始终显示所有按钮 -->
      <Button variant="filled" color="primary" onclick={goToViewer}>漫画查看器</Button>
        <Button variant="filled" color="gray" onclick={goToDownloader}>图片下载器</Button>
        <Button variant="filled" color="secondary" onclick={goToConfig}>应用设置</Button>
    </div>
  </div>
</div>

<!-- 提供占位空间，防止内容被固定定位的header覆盖 -->
<div class="h-16 sm:h-25"></div>