<template>
    <div class="flex flex-col h-full over">
        <Header />

        <Loading v-if="loading" />
        <div v-else-if="selectedImages.length === 0" class="flex-grow flex flex-col items-center justify-center h-full">
            <p class="text-gray-100 mb-5">未找到图片</p>
        </div>
        <main v-else ref="scrollContainer" @scroll="handleScroll"
            class="flex-grow overflow-y-auto p-5 flex flex-col items-center gap-5 scroll-smooth flex-1">
            <div v-for="(image, i) in selectedImages" :key="i">
                <div class="max-w-[1200px] w-full">
                    <img :src="image" :alt="`Manga page ${i + 1}`" class="w-full h-auto block rounded" />
                </div>
            </div>
        </main>
    </div>
</template>

<script setup lang="ts">
import { useMangaStore } from "./stores";
import { MangaService, ProgressService } from "./services";
import { Loading } from "../../components";
import { watch, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { Header } from "./components";

let scrollContainer = ref<HTMLElement | null>(null);
let saveTimeout: number;
let isRestoringProgress = false;
let lastMangaPath = "";
let hasRestoredForCurrentManga = false;


const mangaStore = useMangaStore();
const { loading, selectedImages, mangaPath, mangaName } =
    storeToRefs(mangaStore);

const route = useRoute();
const mangaService = new MangaService();

onMounted(() => {
    // 添加键盘事件监听
    window.addEventListener("keydown", mangaService.handleKeyDown);
    init();
    return () =>
        window.removeEventListener("keydown", mangaService.handleKeyDown);
});

watch(() => route.params.path, (newPath) => {
    if (newPath) {
        init();
    }
});

function init() {
    mangaService.loadManga(route.params.path as string);
}

// 监听漫画路径变化，只在新漫画时恢复位置
watch(() => mangaStore, () => {
    if (mangaPath.value && mangaPath.value !== lastMangaPath) {
        lastMangaPath = mangaPath.value;
        hasRestoredForCurrentManga = false;

        // 延迟恢复，等待图片加载
        if (selectedImages.value.length > 0 && scrollContainer) {
            setTimeout(restoreScrollPosition, 200);
        }
    }
});

watch(() => mangaStore, () => {
    // 监听图片加载完成，如果还没恢复位置则尝试恢复
    if (
        selectedImages.value.length > 0 &&
        scrollContainer &&
        !hasRestoredForCurrentManga
    ) {
        setTimeout(restoreScrollPosition, 200);
    }
});



// 防抖保存进度
function debounceSaveProgress() {
    if (saveTimeout) {
        clearTimeout(saveTimeout);
    }

    saveTimeout = setTimeout(() => {
        if (scrollContainer && mangaPath.value && !isRestoringProgress) {
            const scrollPosition = scrollContainer.value?.scrollTop;
            ProgressService.saveProgress(
                mangaPath.value,
                scrollPosition || 0,
                selectedImages.value.length,
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
    if (!scrollContainer || !mangaPath.value || hasRestoredForCurrentManga)
        return;

    const progress = ProgressService.getProgress(mangaPath.value);
    if (progress && progress.scrollPosition > 0) {
        isRestoringProgress = true;
        hasRestoredForCurrentManga = true;

        // 延迟恢复滚动位置，确保图片已加载
        setTimeout(() => {
            if (scrollContainer.value) {
                scrollContainer.value.scrollTop = progress.scrollPosition;
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
</script>

<style scoped></style>