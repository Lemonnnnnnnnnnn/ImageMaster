<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import {
    GetAllMangas,
    GetLibraries,
    DeleteManga,
    GetImageDataUrl,
    LoadAllLibraries,
  } from "../../../wailsjs/go/viewer/Viewer";

  let mangas: any[] = [];
  let loading = true;
  let libraries = [];
  let showScrollTop = false;
  let scrollY = 0;
  let mangaImages = new Map(); // 缓存图片

  // @ts-expect-error
  onMount(async () => {
    // 加载库和漫画
    await LoadAllLibraries();
    await loadLibraries();
    await loadMangas();
    loading = false;

    // 监听滚动事件
    window.addEventListener("scroll", handleScrollForMainPage);
    return () => window.removeEventListener("scroll", handleScrollForMainPage);
  });

  function handleScrollForMainPage() {
    scrollY = window.scrollY;
    showScrollTop = scrollY > 300;
  }

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
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
        mangaImages.set(
          manga.previewImg,
          await GetImageDataUrl(manga.previewImg),
        );
      }
    }

    loading = false;
  }

  async function deleteManga(event: any, manga: any) {
    event.stopPropagation();
    if (
      confirm(`确定要删除 "${manga.name}" 吗？这将永久删除该文件夹及其内容！`)
    ) {
      loading = true;
      const success = await DeleteManga(manga.path);
      if (success) {
        mangas = mangas.filter((m) => m.path !== manga.path);
      } else {
        alert("删除失败！");
      }
      loading = false;
    }
  }

  function viewManga(manga : any) {
    // 将路径编码后传递给路由
    const encodedPath = encodeURIComponent(manga.path);
    push(`/manga/${encodedPath}`);
  }
</script>

<div class="p-5 max-w-6xl mx-auto">
  <div class="flex justify-between items-center mb-5">
    <h1>漫画查看器</h1>
  </div>

  {#if libraries.length === 0 && !loading}
    <div class="text-center p-10 bg-gray-100 rounded-lg mt-5">
      <h2>欢迎使用漫画查看器</h2>
      <p>请前往"应用设置"页面添加一个漫画库</p>
    </div>
  {/if}

  {#if loading}
    <div class="flex flex-col items-center justify-center p-10">
      <div
        class="w-10 h-10 border-4 border-gray-200 border-l-blue-500 rounded-full animate-spin"
      ></div>
      <p>加载中...</p>
    </div>
  {:else if mangas.length > 0}
    <div class="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-5">
      {#each mangas as manga}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
          class="border border-gray-300 rounded-lg overflow-hidden transition-all duration-200 cursor-pointer hover:-translate-y-1 hover:shadow-lg"
          on:click={() => viewManga(manga)}
        >
          <div class="h-48 overflow-hidden">
            <img
              src={mangaImages.get(manga.previewImg) || ""}
              alt={manga.name}
              class="w-full h-full object-cover"
            />
          </div>
          <div class="p-2">
            <h3
              class="m-0 mb-1 whitespace-nowrap overflow-hidden text-ellipsis"
            >
              {manga.name}
            </h3>
            <p class="m-0 mb-2 text-gray-600">{manga.imagesCount} 张图片</p>
            <button
              on:click={(e) => deleteManga(e, manga)}
              class="px-2 py-1 bg-red-600 text-white border-none rounded cursor-pointer"
            >
              删除
            </button>
          </div>
        </div>
      {/each}
    </div>
  {:else if libraries.length > 0 && !loading}
    <div class="text-center p-10 bg-gray-100 rounded-lg mt-5">
      <p>
        未找到漫画。请确保您添加的漫画库中包含子文件夹，且这些文件夹中含有图片文件。
      </p>
    </div>
  {/if}

  {#if showScrollTop}
    <button
      on:click={scrollToTop}
      class="fixed bottom-8 right-8 w-12 h-12 rounded-full bg-blue-500 text-white text-2xl border-none cursor-pointer flex items-center justify-center shadow-lg"
    >
      ↑
    </button>
  {/if}
</div>
