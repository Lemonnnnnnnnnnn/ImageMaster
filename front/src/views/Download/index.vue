<template>
    <div class="p-8 h-screen flex flex-col">
        <div class="flex gap-4 mt-2">
            <Input v-model="url" class="flex-1" help="please input the target manga url" />
            <Button @click="handleDownload" type="primary">
                <div class="flex items-center gap-2">
                    <Download :size="16" class="text-white" />
                    <span>Download</span>
                </div>
            </Button>
        </div>

        <div class="h-1 border-b border-neutral-300 w-full my-4"></div>

        <div class="flex items-center justify-between gap-2">
            <Switch activeLabel="active task" inactiveLabel="history task" v-model="active" />
            <Button type="danger" @click="downloadStore.clearHistory">
                <div class="flex items-center gap-2">
                    <Trash :size="16" class="text-white" />
                    <span>clear</span>
                </div>
            </Button>
        </div>
        <div class="flex-1 overflow-auto">
            <List class="mt-2" :tasks="showTasks" :downloadStore="downloadStore" />
        </div>
    </div>

</template>

<script setup lang="ts">
import { Button, Input, Switch } from '@/components';
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { toast } from 'vue-sonner';
import { createDownloadHandler } from './services';
import { List } from './components';
import { useDownloadStore } from './stores';
import { Download } from 'lucide-vue-next';
import { storeToRefs } from 'pinia';

const downloadStore = useDownloadStore();

let active = ref(true);

const { historyTasks, activeTasks } = storeToRefs(downloadStore);

const showTasks = computed(() => {
    return active.value ? activeTasks.value : historyTasks.value;
});

onMounted(async () => {
    await downloadStore.initializeStore();
    console.log("组件已挂载，轮询已开始");
});

onUnmounted(() => {
    downloadStore.stopPolling();
    console.log("组件已销毁，轮询已停止");
});

const url = ref('');

// 处理下载
async function handleDownload() {
    if (!url.value.trim()) {
        toast.error('请输入网址');
        return;
    }

    await downloadHandler(url.value.trim());
}

// 创建下载处理器
const downloadHandler = createDownloadHandler({
    onError: (errorMsg) => {
        toast.error(errorMsg);
    },
});




</script>

<style scoped></style>