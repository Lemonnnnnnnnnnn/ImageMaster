import { defineStore } from 'pinia';
import {
  CancelCrawl,
  ClearHistory,
  GetActiveTasks,
  GetHistoryTasks
} from '../../../../wailsjs/go/api/CrawlerAPI';
import type { task } from '../../../../wailsjs/go/models';
import { ref } from 'vue';

// 轮询相关
const POLL_INTERVAL = 1000;
let pollTimer: ReturnType<typeof setInterval> | null = null;

export type DownloadStore = ReturnType<typeof useDownloadStore>;

export const useDownloadStore = defineStore('downloadStore', {
  state: () => ({
    activeTasks: ref<task.DownloadTask[]>([]),
    historyTasks: ref<task.DownloadTask[]>([]),
    activeTab: 'downloading' as 'downloading' | 'history',
    loading: false as boolean
  }),
  getters: {
    activeTasksCount: (state) => state.activeTasks.length
  },
  actions: {
    async initializeStore() {
      try {
        await this.pollTasks();
        this.startPolling();
      } catch (err) {
        console.error('初始化store失败:', err);
      }
    },
    async pollTasks() {
      console.log("pollTasks")
      try {
        const active = await GetActiveTasks();
        this.activeTasks = active;
        const history = await GetHistoryTasks();
        this.historyTasks = history;
      } catch (err) {
        console.error('轮询任务状态出错:', err);
      }
    },
    startPolling() {
      this.stopPolling();
      pollTimer = setInterval(this.pollTasks, POLL_INTERVAL);
    },
    stopPolling() {
      if (pollTimer) {
        clearInterval(pollTimer);
        pollTimer = null;
      }
    },
    async cancelTask(taskId: string) {
      try {
        await CancelCrawl(taskId);
        await this.pollTasks();
      } catch (err) {
        console.error('取消任务出错:', err);
        throw err;
      }
    },
    async clearHistory() {
      try {
        if (!confirm(`确定要清空任务吗？`)) {
          return false;
        }
        await ClearHistory();
        this.historyTasks = [];
      } catch (err) {
        console.error('清除历史出错:', err);
        throw err;
      }
    }
  }
});
