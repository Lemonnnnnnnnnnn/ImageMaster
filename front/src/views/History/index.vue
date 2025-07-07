<template>
    <div class="p-8">
        <div class="flex items-center justify-end gap-2">
            <Button @click="clearHistory">
                <div class="flex items-center text-white gap-2">
                    <Trash :size="16" class="text-white" />
                    <span>清空记录</span>
                </div>
            </Button>
        </div>
        <TaskList :tasks="historyTasks" />
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ClearHistory, GetHistoryTasks } from '../../../wailsjs/go/api/CrawlerAPI';
import type { task } from '../../../wailsjs/go/models';
import { Button, TaskList } from '@/components';

let historyTasks = ref<task.DownloadTask[]>([])

async function clearHistory() {
    try {
        if (!confirm(`确定要清空任务吗？`)) {
            return false;
        }
        await ClearHistory();
        await loadData()
    } catch (err) {
        console.error('清除历史出错:', err);
        throw err;
    }
}

async function loadData() {
    const history = await GetHistoryTasks();
    historyTasks.value = history
    console.log(historyTasks)
}

onMounted(() => {
    loadData()
})
</script>

<style scoped></style>