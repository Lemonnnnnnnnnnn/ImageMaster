<template>
    <main>
        <EmptyState v-if="libraries.length === 0 && !loading" type="no-libraries" />
        <Loading v-if="loading" />
        <MangaGrid v-else-if="mangas.length > 0" />
        <EmptyState v-else-if="libraries.length > 0 && !loading" type="no-mangas" />
    </main>

</template>

<script setup lang="ts">
import { useHomeStore } from "./stores/homeStore";
import { MangaService } from "./services/mangaService";
import { ScrollService } from "./services/scrollService";
import { Loading } from "../../components";
import { EmptyState, MangaGrid } from "./components";
import { onMounted, onUnmounted } from "vue";
import { storeToRefs } from "pinia";

let cleanupScrollListener: (() => void) | null = null;
const homeStore = useHomeStore();
const { loading, libraries, mangas } = storeToRefs(homeStore);

onMounted(async () => {
    // 初始化数据加载
    const mangaService = new MangaService();
    const scrollService = new ScrollService();
    await mangaService.initialize();
    // 初始化滚动监听
    cleanupScrollListener = scrollService.initScrollListener();
});

onUnmounted(() => {
    // 清理滚动监听器
    if (cleanupScrollListener) {
        cleanupScrollListener();
    }
});
</script>

<style scoped></style>