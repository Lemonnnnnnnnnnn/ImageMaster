<script lang="ts">
  import { onMount } from 'svelte';
  import { push, location } from 'svelte-spa-router';
  import { 
    GetMangaImages, 
    GetAllMangas,
    DeleteManga,
    GetImageDataUrl
  } from '../../wailsjs/go/main/Viewer';

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
      
      // 加载图片
      for (let imagePath of imagePaths) {
        try {
          const imageData = await GetImageDataUrl(imagePath);
          selectedImages = [...selectedImages, imageData];
        } catch (error) {
          console.error("加载图片失败:", error);
        }
      }
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
            backToHome();
          }
        } else {
          alert('删除失败！');
          loading = false;
        }
      } catch (error) {
        console.error('删除漫画时出错:', error);
        alert('删除过程中发生错误');
        loading = false;
      }
    }
  }
</script>

<div class="manga-viewer">
  <div class="viewer-header">
    <button class="back-button" on:click={backToHome}>返回</button>
    <h2>{mangaName}</h2>
    <div class="header-controls">
      <button on:click={toggleViewMode} class="view-toggle-btn">
        {viewMode === 'single' ? '切换到滚动模式' : '切换到单图模式'}
      </button>
      <button on:click={toggleNavigation} class="nav-toggle-btn">
        {showNavigation ? '隐藏导航' : '显示导航'}
      </button>
      <button on:click={deleteAndViewNextManga} class="delete-btn">删除</button>
    </div>
  </div>

  {#if loading}
    <div class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>
  {:else if selectedImages.length > 0}
    {#if viewMode === 'single'}
      <!-- 单图浏览模式 -->
      <div class="image-container">
        {#if currentImageIndex > 0}
          <button class="nav-button prev" on:click={prevImage}>←</button>
        {/if}
        
        <img 
          src={selectedImages[currentImageIndex]} 
          alt={`${mangaName} - 图片 ${currentImageIndex + 1}`} 
          on:click={nextImage}
        />
        
        {#if currentImageIndex < selectedImages.length - 1}
          <button class="nav-button next" on:click={nextImage}>→</button>
        {/if}
      </div>
      
      <div class="progress-info">
        {currentImageIndex + 1} / {selectedImages.length}
      </div>
    {:else}
      <!-- 滚动浏览模式 - 原始实现方式 -->
      <div class="scroll-container">
        {#each selectedImages as image, index}
          <div class="manga-page">
            <img src={image} alt={`${mangaName} - 图片 ${index + 1}`} />
          </div>
        {/each}
      </div>
    {/if}
    
    {#if showNavigation}
      <div class="navigation-panel">
        <div class="manga-navigation">
          <button 
            on:click={navigateToPrevManga} 
            disabled={currentMangaIndex <= 0}
          >
            上一部漫画
          </button>
          <button 
            on:click={navigateToNextManga} 
            disabled={currentMangaIndex >= mangas.length - 1}
          >
            下一部漫画
          </button>
        </div>
        
        {#if viewMode === 'single'}
          <div class="image-navigation">
            <input 
              type="range" 
              min="0" 
              max={selectedImages.length - 1} 
              bind:value={currentImageIndex}
            />
          </div>
        {/if}
      </div>
    {/if}
  {:else}
    <div class="no-images">
      <p>未找到图片</p>
      <button on:click={backToHome}>返回首页</button>
    </div>
  {/if}
</div>

<style>
  .manga-viewer {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background-color: #333;
    color: white;
  }

  .viewer-header {
    display: flex;
    align-items: center;
    padding: 10px 20px;
    background-color: #222;
    z-index: 100;
  }

  .back-button {
    margin-right: 20px;
    padding: 5px 10px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  h2 {
    flex: 1;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .header-controls {
    display: flex;
    gap: 10px;
  }

  .view-toggle-btn {
    padding: 5px 10px;
    background-color: #9b59b6;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .nav-toggle-btn {
    padding: 5px 10px;
    background-color: #2ecc71;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .delete-btn {
    padding: 5px 10px;
    background-color: #e74c3c;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .loading, .no-images {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 1;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.1);
    border-left-color: #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* 单图浏览模式样式 */
  .image-container {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
  }

  .image-container img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    cursor: pointer;
  }

  .nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background-color: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    font-size: 24px;
    padding: 20px 15px;
    cursor: pointer;
    z-index: 10;
  }

  .nav-button.prev {
    left: 0;
    border-radius: 0 5px 5px 0;
  }

  .nav-button.next {
    right: 0;
    border-radius: 5px 0 0 5px;
  }

  .progress-info {
    text-align: center;
    padding: 10px;
    background-color: #222;
  }

  /* 滚动浏览模式样式 */
  .scroll-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .manga-page {
    margin-bottom: 20px;
    max-width: 100%;
  }

  .manga-page img {
    max-width: 100%;
    height: auto;
    display: block;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  }

  .navigation-panel {
    padding: 10px 20px;
    background-color: #222;
    border-top: 1px solid #444;
  }

  .manga-navigation {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .manga-navigation button {
    padding: 5px 10px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .manga-navigation button:disabled {
    background-color: #95a5a6;
    cursor: not-allowed;
  }

  .image-navigation {
    padding: 10px 0;
  }

  .image-navigation input {
    width: 100%;
  }
</style> 