<template>
    <div v-if="tasks.length === 0">
        <div class="text-center h-80 flex items-center justify-center text-neutral-100">
            <p> no task</p>
        </div>
    </div>
    <table v-else class="w-full text-neutral-100 scroll-auto">
        <thead>
            <tr class="border-b border-neutral-300">
                <th>name</th>
                <th>url</th>
                <th>status</th>
                <th>progress</th>
                <th>startTime</th>
                <th>endTime</th>
                <th>action</th>
            </tr>
        </thead>

        <tbody>
            <tr v-for="task in tasks" :key="task.id" class="border-b border-neutral-300">
                <td :title="task.name">{{ task.name }}</td>
                <td :title="task.url">{{ task.url }}</td>
                <td>
                    <div class="flex items-center justify-center gap-2">
                        <component :is="getStatusIcon(task.status)?.icon" :size="16"
                            :class="getStatusIcon(task.status)?.class" />
                        <span>{{ formatStatus(task.status) }}</span>
                    </div>
                </td>
                <td>{{ calculateProgressPercentage(task.progress.current, task.progress.total) }}%</td>
                <td>{{ formatTime(task.startTime) }}</td>
                <td>{{ formatTime(task.completeTime) }}</td>
                <td><button v-if="task.status === 'pending' || task.status === 'downloading'"
                        class="text-red-500 cursor-pointer" @click="downloadStore.cancelTask(task.id)">cancel</button>
                </td>
            </tr>
        </tbody>

    </table>

</template>

<script setup lang="ts">
import { calculateProgressPercentage, formatStatus, formatTime, getStatusClass } from '../services';
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

</script>

<style scoped>
@reference "tailwindcss";

td,
th,
tr {
    @apply text-center text-xs p-2;
}

tr {
    @apply hover:bg-neutral-800;
}

td {
    @apply max-w-32 truncate
}
</style>