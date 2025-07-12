<template>
    <main>
        <EmptyState v-if="activeLibrary === '' && !loading" type="no-libraries" />
        <Loading v-if="loading" />
        <MangaGrid v-else-if="mangas.length > 0" />
        <EmptyState v-else-if="activeLibrary !== '' && !loading" type="no-mangas" />
    </main>

</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { onMounted, ref } from "vue";
import { Loading } from "../../components";
import { EmptyState, MangaGrid } from "./components";
import { MangaService } from "./services/mangaService";
import { useHomeStore } from "./stores/homeStore";
import { GetActiveLibrary } from "../../../wailsjs/go/config/API";

const homeStore = useHomeStore();
const { loading, mangas } = storeToRefs(homeStore);

let activeLibrary = ref("");

async function getActiveLibrary() {
    const library = await GetActiveLibrary();
    activeLibrary.value = library;
}

onMounted(async () => {
    // 初始化数据加载
    const mangaService = new MangaService();
    await mangaService.initialize();
    getActiveLibrary();
});

</script>

<style scoped></style>