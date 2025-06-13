<script lang="ts">
  import { onMount } from 'svelte';
  import { push, location } from 'svelte-spa-router';
  import { 
    GetMangaImages, 
    GetAllMangas,
    DeleteManga,
    GetImageDataUrl,
    LoadAllLibraries
  } from '../../../wailsjs/go/viewer/Viewer';
  import QuickDownloader from '../../components/QuickDownloader.svelte';

  // 从路由参数获取漫画路径
  export let params: {
    path?: string
  } = {};
  
  let mangaPath = '';
  let mangaName = '';
  let selectedImages = [];
  let currentImageIndex = 0;
  let loading = true;
  let mangas = []; // 用于上一部/下一部功能
  let currentMangaIndex = -1;
  let showNavigation = false; // 控制导航按钮的显示
  let viewMode = 'scroll'; // 'single' 单图浏览模式 或 'scroll' 滚动浏览模式
  
  // 响应式声明，监听 params 变化
  $: if (params.path) {
    loadManga(params.path);
  }

  // 页面加载函数，现在可以被 params 变化触发
  async function loadManga(path: string) {
    try {
      loading = true;
      
      // 解码路径参数
      mangaPath = decodeURIComponent(path);
      
      // 获取所有漫画以支持导航功能
      mangas = await GetAllMangas();
      currentMangaIndex = mangas.findIndex(m => m.path === mangaPath);
      
      if (currentMangaIndex >= 0) {
        mangaName = mangas[currentMangaIndex].name;
      } else {
        mangaName = mangaPath.split('/').pop();
      }
      
      // 重置图片和选中索引
      selectedImages = [];
      currentImageIndex = 0;
      
      // 获取所有图片
      await loadImages();
      
    } catch (error) {
      console.error('加载漫画失败:', error);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    // 添加键盘事件监听
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  });

  async function loadImages() {
    try {
      // 获取所有图片路径
      const imagePaths = await GetMangaImages(mangaPath);
      
      // 并行加载所有图片，保持顺序
      const imagePromises = imagePaths.map(async (imagePath) => {
        try {
          return await GetImageDataUrl(imagePath);
        } catch (error) {
          console.error(`加载图片失败: ${imagePath}`, error);
          return null; // 返回 null 表示加载失败
        }
      });

      // 等待所有图片加载完成
      const loadedImages = await Promise.all(imagePromises);
      
      // 过滤掉加载失败的图片（null值）
      selectedImages = loadedImages.filter(img => img !== null);
    } catch (error) {
      console.error("获取图片路径失败:", error);
    }
  }

  function handleKeyDown(event) {
    if (viewMode === 'single') {
      if (event.key === 'ArrowRight' || event.key === ' ') {
        nextImage();
      } else if (event.key === 'ArrowLeft') {
        prevImage();
      } else if (event.key === 'Escape') {
        backToHome();
      }
    } else {
      if (event.key === 'Escape') {
        backToHome();
      }
    }
  }

  function toggleViewMode() {
    viewMode = viewMode === 'single' ? 'scroll' : 'single';
  }

  function nextImage() {
    if (currentImageIndex < selectedImages.length - 1) {
      currentImageIndex++;
    } else if (currentMangaIndex < mangas.length - 1) {
      // 当前漫画已浏览完，询问是否进入下一部
      if (confirm('已浏览完当前漫画，是否查看下一部？')) {
        navigateToNextManga();
      }
    }
  }

  function prevImage() {
    if (currentImageIndex > 0) {
      currentImageIndex--;
    }
  }

  function backToHome() {
    push('/');
  }

  function toggleNavigation() {
    showNavigation = !showNavigation;
  }

  function navigateToNextManga() {
    if (currentMangaIndex < mangas.length - 1) {
      const nextManga = mangas[currentMangaIndex + 1];
      const encodedPath = encodeURIComponent(nextManga.path);
      
      // 使用替代路由方案处理相同路径不同参数的导航
      const currentLocation = $location;
      if (currentLocation === `/manga/${params.path}`) {
        // 如果当前已经在漫画页面，采用 params.path 改变触发重新加载
        loadManga(nextManga.path);
        
        // 更新 URL 但不触发导航事件
        window.history.pushState(null, '', `/#/manga/${encodedPath}`);
      } else {
        // 否则正常导航
        push(`/manga/${encodedPath}`);
      }
    }
  }

  function navigateToPrevManga() {
    if (currentMangaIndex > 0) {
      const prevManga = mangas[currentMangaIndex - 1];
      const encodedPath = encodeURIComponent(prevManga.path);
      
      // 使用替代路由方案处理相同路径不同参数的导航
      const currentLocation = $location;
      if (currentLocation === `/manga/${params.path}`) {
        // 如果当前已经在漫画页面，采用 params.path 改变触发重新加载
        loadManga(prevManga.path);
        
        // 更新 URL 但不触发导航事件
        window.history.pushState(null, '', `/#/manga/${encodedPath}`);
      } else {
        // 否则正常导航
        push(`/manga/${encodedPath}`);
      }
    }
  }

  async function deleteAndViewNextManga() {
    if (currentMangaIndex >= 0 && confirm(`确定要删除 "${mangaName}" 并查看下一部漫画吗？`)) {
      loading = true;
      
      try {
        // 记录下一个漫画的位置，因为删除后数组会变化
        const hasNextManga = currentMangaIndex < mangas.length - 1;
        const nextMangaPath = hasNextManga ? mangas[currentMangaIndex + 1].path : null;
        
        // 执行删除操作
        const success = await DeleteManga(mangaPath);
        
        if (success) {
          if (nextMangaPath) {
            // 重要：在导航前设置 loading 为 false，防止新页面保持加载状态
            loading = false;
            
            // 使用替代路由方案处理相同路径不同参数的导航
            const encodedPath = encodeURIComponent(nextMangaPath);
            const currentLocation = $location;
            if (currentLocation === `/manga/${params.path}`) {
              // 如果当前已经在漫画页面，采用直接加载新数据的方式
              loadManga(nextMangaPath);
              
              // 更新 URL 但不触发导航事件
              window.history.pushState(null, '', `/#/manga/${encodedPath}`);
            } else {
              // 否则正常导航
              push(`/manga/${encodedPath}`);
            }
          } else {
            // 如果没有下一部漫画，返回首页
            push('/');
          }
        } else {
          alert('删除失败!');
          loading = false;
        }
      } catch (error) {
        console.error('删除漫画失败:', error);
        loading = false;
      }
    }
  }

  // 处理快速下载成功事件
  function handleDownloadStarted(event) {
    const { taskId, url } = event.detail;
    console.log(`下载任务已创建: ${taskId}, URL: ${url}`);
    // 可以在这里添加通知或其他处理逻辑
  }
</script>

<div class="manga-viewer">
  <div class="viewer-header">
    <div class="left-controls">
      <button on:click={backToHome}>返回</button>
      <button on:click={toggleViewMode}>
        {viewMode === 'single' ? '切换为滚动模式' : '切换为单图模式'}
      </button>
    </div>
    <div class="center-title">
      <!-- <h2>{mangaName} ({currentImageIndex + 1}/{selectedImages.length})</h2> -->
      <h2>{mangaName} ({selectedImages.length})</h2>
    </div>
    <div class="right-controls">
      <button on:click={toggleNavigation}>
        {showNavigation ? '隐藏导航' : '显示导航'}
      </button>
      <button class="delete-btn" on:click={deleteAndViewNextManga}>删除并看下一部</button>
    </div>
  </div>
  
  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>
  {:else if selectedImages.length === 0}
    <div class="no-images">
      <p>未找到图片</p>
      <button on:click={backToHome}>返回</button>
    </div>
  {:else if viewMode === 'single'}
    <div class="single-view">
      <img src={selectedImages[currentImageIndex]} alt="Manga page {currentImageIndex + 1}" />
      
      <div class="single-view-controls">
        <button on:click={prevImage} disabled={currentImageIndex === 0}>上一页</button>
        <span>{currentImageIndex + 1} / {selectedImages.length}</span>
        <button on:click={nextImage} disabled={currentImageIndex === selectedImages.length - 1}>下一页</button>
      </div>
    </div>
  {:else}
    <div class="scroll-view">
      {#each selectedImages as image, i}
        <div class="scroll-image-container">
          <img src={image} alt="Manga page {i + 1}" />
        </div>
      {/each}
    </div>
  {/if}
  
  {#if showNavigation}
    <div class="navigation-panel">
      <button on:click={navigateToPrevManga} disabled={currentMangaIndex <= 0}>
        上一部漫画
      </button>
      <button on:click={navigateToNextManga} disabled={currentMangaIndex >= mangas.length - 1}>
        下一部漫画
      </button>
    </div>
  {/if}
  
  <!-- 快速下载组件 -->
  <QuickDownloader on:downloadStarted={handleDownloadStarted} />
</div>

<style>
  .manga-viewer {
    height: 100%;
    display: flex;
    flex-direction: column;
    background-color: #1a1a1a;
    color: #f0f0f0;
  }
  
  .viewer-header {
    padding: 10px 20px;
    background-color: #2a2a2a;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .center-title {
    text-align: center;
    flex-grow: 1;
  }
  
  .center-title h2 {
    margin: 0;
    font-size: 18px;
    /* white-space: nowrap; */
    overflow: hidden;
    text-overflow: ellipsis;
  }

  
  .left-controls, .right-controls {
    display: flex;
    gap: 10px;
    min-width: 240px;
  }
  
  button {
    padding: 8px 12px;
    background-color: #4a6fa5;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:disabled {
    background-color: #555;
    cursor: not-allowed;
  }
  
  .delete-btn {
    background-color: #c62828;
  }
  
  .loading, .no-images {
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
  
  .single-view {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    position: relative;
  }
  
  .single-view img {
    max-width: 100%;
    max-height: calc(100vh - 120px);
    object-fit: contain;
  }
  
  .single-view-controls {
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 20px;
    align-items: center;
    padding: 10px;
    background-color: rgba(42, 42, 42, 0.8);
    border-radius: 8px;
  }
  
  .scroll-view {
    flex-grow: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
  
  .scroll-image-container {
    max-width: 800px;
    width: 100%;
  }
  
  .scroll-image-container img {
    width: 100%;
    height: auto;
    display: block;
    border-radius: 4px;
  }
  
  .navigation-panel {
    position: fixed;
    bottom: 20px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    background-color: rgba(42, 42, 42, 0.8);
    padding: 10px;
    border-radius: 8px;
  }
</style>