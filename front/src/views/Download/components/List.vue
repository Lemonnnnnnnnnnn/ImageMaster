<template>
    <div v-if="tasks.length === 0">
        <div class="text-center h-80 flex items-center justify-center text-neutral-100">
            <p> no task</p>
        </div>
    </div>
    <table v-else class="w-full text-neutral-100 scroll-auto">
        <thead>
            <tr class="border-b border-neutral-300">
                <th>名字</th>
                <!-- <th>url</th> -->
                <th>状态</th>
                <th>进度</th>
                <!-- <th>startTime</th> -->
                <th>完成时间</th>
                <!-- <th>耗时</th> -->
                <th>操作</th>
            </tr>
        </thead>

        <tbody>
            <tr v-for="task in tasks" :key="task.id" class="border-b border-neutral-500/50">
                <td :title="task.name" class="max-w-48">{{ task.name }}</td>
                <!-- <td :title="task.url">{{ task.url }}</td> -->
                <td class="flex justify-center">
                    <!-- <div class="border-1 border-neutral-300 rounded-xl px-2 py-1 w-24"> -->
                    <div class="flex items-center justify-center gap-2">
                        <component :is="getStatusIcon(task.status)?.icon" :size="16"
                            :class="getStatusIcon(task.status)?.class" />
                        <span>{{ formatStatus(task.status) }}</span>
                    </div>
                    <!-- </div> -->
                </td>
                <!-- <td>{{ calculateProgressPercentage(task.progress.current, task.progress.total) }}%</td> -->
                <td>
                    <div class="border border-neutral-300 rounded-xl h-2 w-full">
                        <div class="bg-neutral-300 rounded-xl h-full transition-all duration-300"
                            :style="{ width: `${calculateProgressPercentage(task.progress.current, task.progress.total)}%` }">
                        </div>
                    </div>
                </td>
                <!-- <td>{{ formatTime(task.startTime) }}</td> -->
                <td>{{ task.status === 'completed' ? formatTime(task.completeTime) : '-' }}</td>
                <!-- <td>{{ calculateTimeDifference(task.startTime, task.status === 'completed' ? task.completeTime : Date.now()) }}</td> -->
                <td><button v-if="task.status === 'pending' || task.status === 'downloading'"
                        class="text-sky-500 cursor-pointer" @click="downloadStore.cancelTask(task.id)">取消</button>
                </td>
            </tr>
        </tbody>

    </table>

</template>

<script setup lang="ts">
import { calculateProgressPercentage, formatTime } from '../services';
import type { task } from '../../../../wailsjs/go/models';
import type { DownloadStore } from '../stores';
import { Loader, ArrowBigDownDash, CircleCheck, CircleX, CircleOff } from 'lucide-vue-next';

defineProps<{
    tasks: task.DownloadTask[]
    downloadStore: DownloadStore
}>();

function getStatusIcon(status: string) {
    if (status === "pending") {
        return { icon: Loader, class: 'animate-spin' };
    } else if (status === "parsing") {
        return { icon: Loader, class: 'animate-spin' };
    } else if (status === "downloading") {
        return { icon: ArrowBigDownDash, class: 'animate-bounce' };
    } else if (status === "completed") {
        return { icon: CircleCheck, class: '' };
    } else if (status === "failed") {
        return { icon: CircleX, class: '' };
    } else if (status === "cancelled") {
        return { icon: CircleOff, class: '' };
    }
}

function formatStatus(status: string): string {
    const statusMap: Record<string, string> = {
        'pending': '等待中',
        "parsing": "解析中",
        'downloading': '下载中',
        'completed': '已完成',
        'failed': '失败',
        'cancelled': '已取消'
    };
    return statusMap[status] || status;
}


</script>

<style scoped>
@reference "tailwindcss";

td,
th,
tr {
    @apply text-center text-xs px-2 py-4;
}

tr {
    @apply hover:bg-neutral-800;
}

td {
    @apply truncate
}
</style>