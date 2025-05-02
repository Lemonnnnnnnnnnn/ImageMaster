<script lang="ts">
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { 
    GetAllMangas, 
    SelectLibrary, 
    GetLibraries, 
    DeleteManga,
    GetImageDataUrl
  } from '../../../wailsjs/go/main/Viewer';

  let mangas = [];
  let loading = true;
  let libraries = [];
  let showScrollTop = false;
  let scrollY = 0;
  let mangaImages = new Map(); // 缓存图片

  onMount(async () => {
    // 加载库和漫画
    await loadLibraries();
    await loadMangas();
    loading = false;

    // 监听滚动事件
    window.addEventListener('scroll', handleScrollForMainPage);
    return () => window.removeEventListener('scroll', handleScrollForMainPage);
  });

  function handleScrollForMainPage() {
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

  function viewManga(manga) {
    // 将路径编码后传递给路由
    const encodedPath = encodeURIComponent(manga.path);
    push(`/manga/${encodedPath}`);
  }
</script>

<div class="home-container">
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
            >
              删除
            </button>
          </div>
        </div>
      {/each}
    </div>
  {:else if libraries.length > 0 && !loading}
    <div class="no-manga">
      <p>未找到漫画。请确保您添加的漫画库中包含子文件夹，且这些文件夹中含有图片文件。</p>
    </div>
  {/if}

  {#if showScrollTop}
    <button class="scroll-top-btn" on:click={scrollToTop}>
      ↑
    </button>
  {/if}
</div>

<style>
  .home-container {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .add-library-btn {
    padding: 8px 16px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .welcome, .no-manga {
    text-align: center;
    padding: 40px;
    background-color: #f5f5f5;
    border-radius: 8px;
    margin-top: 20px;
  }

  .libraries {
    margin-bottom: 20px;
    padding: 10px;
    background-color: #f5f5f5;
    border-radius: 4px;
  }

  .libraries ul {
    margin: 0;
    padding-left: 20px;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-left-color: #4CAF50;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .manga-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
    margin-top: 20px;
  }

  .manga-card {
    background-color: #fff;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s;
    cursor: pointer;
  }

  .manga-card:hover {
    transform: translateY(-5px);
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
    padding: 15px;
    position: relative;
  }

  .manga-info h3 {
    margin: 0 0 10px 0;
    font-size: 16px;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
  }

  .manga-info p {
    margin: 0;
    color: #666;
    font-size: 14px;
  }

  .delete-btn {
    position: absolute;
    top: 15px;
    right: 15px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 5px 10px;
    font-size: 12px;
    cursor: pointer;
    display: none;
  }

  .manga-card:hover .delete-btn {
    display: block;
  }

  .scroll-top-btn {
    position: fixed;
    bottom: 20px;
    right: 20px;
    width: 40px;
    height: 40px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 50%;
    font-size: 20px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
  }
</style> 