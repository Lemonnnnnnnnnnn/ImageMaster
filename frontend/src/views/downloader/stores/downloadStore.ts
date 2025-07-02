import { writable, derived } from 'svelte/store';
import type { models } from '../../../../wailsjs/go/models';
import {
  GetActiveTasks,
  GetHistoryTasks,
  CancelDownload,
  ClearHistory
} from '../../../../wailsjs/go/api/DownloadAPI';
import { GetOutputDir } from '../../../../wailsjs/go/storage/StorageAPI';
import { toast } from 'svelte-sonner';

// 基础状态
export const activeTasks = writable<models.DownloadTask[]>([]);
export const historyTasks = writable<models.DownloadTask[]>([]);
export const activeTab = writable<'downloading' | 'history'>('downloading');
export const outputDir = writable<string>('');
export const loading = writable<boolean>(false);

// 轮询相关
const POLL_INTERVAL = 1000;
let pollTimer: ReturnType<typeof setInterval> | null = null;

// 派生状态
export const activeTasksCount = derived(activeTasks, $activeTasks => $activeTasks.length);

// 初始化函数
export async function initializeStore() {
  try {
    const dir = await GetOutputDir();
    outputDir.set(dir);
    await pollTasks();
    startPolling();
  } catch (err) {
    console.error('初始化store失败:', err);
  }
}

// 轮询任务状态
export async function pollTasks() {
  try {
    const active = await GetActiveTasks();
    
    activeTasks.set(active);
    
    // 如果当前显示历史标签，也获取历史任务
    const currentTab = await new Promise<'downloading' | 'history'>(resolve => {
      const unsubscribe = activeTab.subscribe(tab => {
        resolve(tab);
        // unsubscribe();
      });
    });
    
    if (currentTab === 'history') {
      const history = await GetHistoryTasks();
      historyTasks.set(history);
    }
  } catch (err) {
    console.error('轮询任务状态出错:', err);
  }
}

// 开始轮询
export function startPolling() {
  stopPolling();
  pollTimer = setInterval(pollTasks, POLL_INTERVAL);
}

// 停止轮询
export function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
}

// 取消下载任务
export async function cancelTask(taskId: string) {
  try {
    await CancelDownload(taskId);
    await pollTasks();
  } catch (err) {
    console.error('取消任务出错:', err);
    throw err;
  }
}

// 切换标签
export async function switchTab(tab: 'downloading' | 'history') {
  activeTab.set(tab);
  
  if (tab === 'history') {
    try {
      const history = await GetHistoryTasks();
      historyTasks.set(history);
    } catch (err) {
      console.error('获取历史任务失败:', err);
    }
  }
}

// 清除历史记录
export async function clearHistory() {
  try {
    await ClearHistory();
    historyTasks.set([]);
  } catch (err) {
    console.error('清除历史出错:', err);
    throw err;
  }
}