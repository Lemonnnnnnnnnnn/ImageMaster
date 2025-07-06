<template>
    <header class="py-2 px-4 bg-transparent flex justify-between items-center gap-4">
        <div class="flex items-center gap-2 shrink-0">
            <Button @click="router.push('/')">
                <div class="flex items-center gap-2">
                    <ArrowLeft :size="16" class="text-white" />
                    <span>返回</span>
                </div>

            </Button>

        </div>
        <div class="text-white text-sm font-bold flex-1 truncate">
            {{ mangaName }}
        </div>

        <div class="flex items-center gap-2 shrink-0">
            <Button @click="showNavigation = !showNavigation">
                <div class="flex items-center gap-2">
                    <EyeClosed v-if="!showNavigation" :size="16" class="text-white" />
                    <Eye :size="16" class="text-white" v-else />
                    <span>显示导航</span>
                </div>
            </Button>

            <Button @click="showQuickDownloadModal = true">
                <div class="flex items-center gap-2">
                    <Download :size="16" class="text-white" />
                    <span>快速下载</span>
                </div>
            </Button>

            <Button @click="() => mangaService.deleteAndViewNextManga()">
                <div class="flex items-center gap-2">
                    <Trash :size="16" class="text-white" />
                    <span>删除并看下一部</span>
                </div>
            </Button>
        </div>
    </header>
    <QuickDownloadModal v-model="showQuickDownloadModal" />
</template>

<script setup lang="ts">
import { ArrowLeft, Download, Eye, EyeClosed, Trash } from 'lucide-vue-next';
import { useRouter } from 'vue-router';
import { useMangaStore } from '../stores';
import { storeToRefs } from 'pinia';
import { Button, QuickDownloadModal } from '../../../components';
import { ref } from 'vue';

const router = useRouter();
const mangaStore = useMangaStore();
const { mangaName } = storeToRefs(mangaStore);

defineProps({
    mangaService: {
        type: Object,
        required: true
    }
})

let showNavigation = ref(false);
let showQuickDownloadModal = ref(false);
</script>

<style scoped></style>