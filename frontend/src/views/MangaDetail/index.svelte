<script lang="ts">
    import { mangaStore } from "./stores";
    import { MangaService, ProgressService } from "./services";
    import { ArrowLeft } from "lucide-svelte";
    import { push } from "svelte-spa-router";
    import { onMount } from "svelte";
    import { Loading } from "../../components";

    let scrollContainer: HTMLElement | null = $state(null);
    let saveTimeout: number;
    let isRestoringProgress = false;
    let lastMangaPath = "";
    let hasRestoredForCurrentManga = false;
    const { loading, selectedImages, mangaPath, mangaName } = $mangaStore;

    onMount(() => {
        // 添加键盘事件监听
        window.addEventListener("keydown", MangaService.handleKeyDown);
        return () =>
            window.removeEventListener("keydown", MangaService.handleKeyDown);
    });

    // 防抖保存进度
    function debounceSaveProgress() {
        if (saveTimeout) {
            clearTimeout(saveTimeout);
        }

        saveTimeout = setTimeout(() => {
            if (scrollContainer && mangaPath && !isRestoringProgress) {
                const scrollPosition = scrollContainer.scrollTop;
                ProgressService.saveProgress(
                    mangaPath,
                    scrollPosition,
                    selectedImages.length,
                );
            }
        }, 1000); // 1秒防抖
    }

    // 处理滚动事件
    function handleScroll() {
        debounceSaveProgress();
    }

    // 恢复滚动位置
    function restoreScrollPosition() {
        if (!scrollContainer || !mangaPath || hasRestoredForCurrentManga)
            return;

        const progress = ProgressService.getProgress(mangaPath);
        if (progress && progress.scrollPosition > 0) {
            isRestoringProgress = true;
            hasRestoredForCurrentManga = true;

            // 延迟恢复滚动位置，确保图片已加载
            setTimeout(() => {
                if (scrollContainer) {
                    scrollContainer.scrollTop = progress.scrollPosition;
                    console.log(
                        `已恢复到上次阅读位置：${progress.scrollPosition}px`,
                    );
                }
                isRestoringProgress = false;
            }, 100);
        } else {
            hasRestoredForCurrentManga = true;
        }
    }

    let { params }: { params: { path: string } } = $props();

    $effect(() => {
        if (params.path) {
            MangaService.loadManga(params.path);
        }
    });

    // 监听漫画路径变化，只在新漫画时恢复位置
    $effect(() => {
        if (mangaPath && mangaPath !== lastMangaPath) {
            lastMangaPath = mangaPath;
            hasRestoredForCurrentManga = false;

            // 延迟恢复，等待图片加载
            if (selectedImages.length > 0 && scrollContainer) {
                setTimeout(restoreScrollPosition, 200);
            }
        }
    });

    $effect(() => {
        // 监听图片加载完成，如果还没恢复位置则尝试恢复
        if (
            selectedImages.length > 0 &&
            scrollContainer &&
            !hasRestoredForCurrentManga
        ) {
            setTimeout(restoreScrollPosition, 200);
        }
    });
</script>

<div class="flex flex-col h-full">
    <header class="p-4 bg-transparent flex items-center gap-4">
        <button
            onclick={() => {
                push("/");
            }}
            class="cursor-pointer"
        >
            <ArrowLeft size={20} class="text-white" />
        </button>
        <div class="text-white text-xs font-bold max-w-[500px] truncate">
            {mangaName}
        </div>
    </header>

    {#if loading}
        <Loading />
    {:else if selectedImages.length === 0}
        <div class="flex-grow flex flex-col items-center justify-center h-full">
            <p class="text-gray-100 mb-5">未找到图片</p>
        </div>
    {:else}
        <main
            bind:this={scrollContainer}
            onscroll={handleScroll}
            class="flex-grow overflow-y-auto p-5 flex flex-col items-center gap-5 scroll-smooth flex-1"
        >
            {#each selectedImages as image, i}
                <div class="max-w-[1200px] w-full">
                    <img
                        src={image}
                        alt="Manga page {i + 1}"
                        class="w-full h-auto block rounded"
                    />
                </div>
            {/each}
        </main>
    {/if}
</div>
