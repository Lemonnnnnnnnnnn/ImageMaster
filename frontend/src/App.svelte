<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    GetAllMangas, 
    GetMangaImages, 
    SelectLibrary, 
    GetLibraries, 
    DeleteManga,
    GetImageDataUrl
  } from '../wailsjs/go/main/App';

  let mangas = [];
  let selectedManga = null;
  let selectedImages = [];
  let currentImageIndex = 0;
  let showViewer = false;
  let loading = true;
  let libraries = [];

  // 向上滚动按钮显示控制
  let showScrollTop = false;
  let scrollY = 0;

  let mangaImages = new Map(); // 缓存图片

  let showMangaView = false; // 控制是否显示漫画查看页面

  onMount(async () => {
    // 加载库和漫画
    await loadLibraries();
    await loadMangas();
    loading = false;

    // 监听滚动事件
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  });

  function handleScroll() {
    scrollY = window.scrollY;
    showScrollTop = scrollY > 300;
  }

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  async function loadLibraries() {
    libraries = await GetLibraries();
  }

  async function loadMangas() {
    loading = true;
    mangas = await GetAllMangas();
    console.log(mangas);
    
    // 预加载每个漫画的预览图
    for (let manga of mangas) {
      if (!mangaImages.has(manga.previewImg)) {
        mangaImages.set(manga.previewImg, await GetImageDataUrl(manga.previewImg));
      }
    }
    
    loading = false;
  }

  async function chooseLibrary() {
    loading = true;
    const newLib = await SelectLibrary();
    if (newLib) {
      await loadLibraries();
      await loadMangas();
    }
    loading = false;
  }

  async function viewManga(manga) {
    loading = true;
    selectedManga = manga;
    selectedImages = [];
    
    // 获取所有图片路径
    const imagePaths = await GetMangaImages(manga.path);
    
    // 加载所有图片
    for (let i = 0; i < imagePaths.length; i++) {
      try {
        const imageData = await GetImageDataUrl(imagePaths[i]);
        selectedImages = [...selectedImages, imageData];
        
        // 加载第一张后即显示页面
        if (i === 0) {
          showMangaView = true;
        }
      } catch (error) {
        console.error("加载图片失败:", error);
      }
    }
    
    loading = false;
  }

  function backToHome() {
    showMangaView = false;
    selectedManga = null;
    selectedImages = [];
  }

  async function deleteManga(event, manga) {
    event.stopPropagation();
    if (confirm(`确定要删除 "${manga.name}" 吗？这将永久删除该文件夹及其内容！`)) {
      loading = true;
      const success = await DeleteManga(manga.path);
      if (success) {
        mangas = mangas.filter(m => m.path !== manga.path);
      } else {
        alert('删除失败！');
      }
      loading = false;
    }
  }
</script>

<main>
  {#if !showMangaView}
    <!-- 主页面 -->
    <div class="header">
      <h1>漫画查看器</h1>
      <button on:click={chooseLibrary} class="add-library-btn">添加漫画库</button>
    </div>

    {#if libraries.length === 0 && !loading}
      <div class="welcome">
        <h2>欢迎使用漫画查看器</h2>
        <p>请点击"添加漫画库"按钮选择一个包含漫画的文件夹</p>
      </div>
    {:else if libraries.length > 0}
      <div class="libraries">
        <h3>已添加的漫画库：</h3>
        <ul>
          {#each libraries as lib}
            <li>{lib}</li>
          {/each}
        </ul>
      </div>
    {/if}

    {#if loading}
      <div class="loading">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    {:else if mangas.length > 0}
      <div class="manga-grid">
        {#each mangas as manga}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <div class="manga-card" on:click={() => viewManga(manga)}>
            <div class="manga-preview">
              <img src={mangaImages.get(manga.previewImg) || ''} alt={manga.name} />
            </div>
            <div class="manga-info">
              <h3>{manga.name}</h3>
              <p>{manga.imagesCount} 张图片</p>
              <button 
                class="delete-btn" 
                on:click={(e) => deleteManga(e, manga)}
                title="删除"
              >
                ✕
              </button>
            </div>
          </div>
        {/each}
      </div>
    {:else if libraries.length > 0 && !loading}
      <div class="no-mangas">
        <p>未找到漫画。请确保您的漫画库中包含有图片的子文件夹。</p>
      </div>
    {/if}

    {#if showScrollTop}
      <button class="scroll-top-btn" on:click={scrollToTop}>
        ↑
      </button>
    {/if}
  {:else}
    <!-- 漫画查看页面 -->
    <div class="manga-view">
      <div class="manga-view-header">
        <button class="back-btn" on:click={backToHome}>← 返回</button>
        <h2>{selectedManga?.name}</h2>
        <div class="placeholder"></div>
      </div>
      
      {#if loading}
        <div class="loading">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>
      {:else}
        <div class="manga-view-content">
          {#each selectedImages as image}
            <div class="manga-page">
              <img src={image} alt="漫画页面" />
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    background-color: #f5f5f5;
    color: #333;
  }

  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  h1 {
    margin: 0;
    font-size: 24px;
  }

  .add-library-btn {
    background-color: #4CAF50;
    color: white;
    border: none;
    padding: 10px 15px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .add-library-btn:hover {
    background-color: #45a049;
  }

  .welcome {
    text-align: center;
    margin: 50px 0;
    padding: 30px;
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  }

  .libraries {
    margin-bottom: 20px;
    background-color: white;
    padding: 15px;
    border-radius: 8px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  }

  .libraries h3 {
    margin-top: 0;
    margin-bottom: 10px;
    font-size: 16px;
  }

  .libraries ul {
    margin: 0;
    padding-left: 20px;
  }

  .libraries li {
    margin-bottom: 5px;
    word-break: break-all;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
  }

  .spinner {
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-radius: 50%;
    border-top: 4px solid #3498db;
    width: 30px;
    height: 30px;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .manga-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
  }

  .manga-card {
    background-color: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    transition: transform 0.2s;
    cursor: pointer;
    position: relative;
  }

  .manga-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.2);
  }

  .manga-preview {
    height: 200px;
    overflow: hidden;
  }

  .manga-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .manga-info {
    padding: 10px;
    position: relative;
  }

  .manga-info h3 {
    margin: 0 0 5px 0;
    font-size: 16px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .manga-info p {
    margin: 0;
    font-size: 14px;
    color: #666;
  }

  .delete-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    background-color: rgba(255, 0, 0, 0.7);
    color: white;
    border: none;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s;
  }

  .manga-card:hover .delete-btn {
    opacity: 1;
  }

  .delete-btn:hover {
    background-color: red;
  }

  .no-mangas {
    text-align: center;
    margin: 50px 0;
    color: #666;
  }

  .scroll-top-btn {
    position: fixed;
    bottom: 30px;
    right: 30px;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background-color: #333;
    color: white;
    border: none;
    font-size: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 2px 10px rgba(0,0,0,0.3);
    z-index: 100;
  }

  .manga-view {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: #f5f5f5;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .manga-view-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 20px;
    background-color: white;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    z-index: 100;
  }

  .manga-view-header h2 {
    margin: 0;
    font-size: 18px;
    text-align: center;
  }

  .back-btn {
    background: none;
    border: none;
    color: #333;
    font-size: 16px;
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 4px;
  }

  .back-btn:hover {
    background-color: #f0f0f0;
  }

  .placeholder {
    width: 80px; /* 与back-btn大致相同宽度，用于居中标题 */
  }

  .manga-view-content {
    flex: 1;
    overflow-y: auto;
    padding: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .manga-page {
    width: 100%;
    display: flex;
    justify-content: center;
  }

  .manga-page img {
    max-width: 100%;
    height: auto;
    display: block;
  }
</style>
