// 格式化任务状态
export function formatStatus(status: string): string {
  const statusMap: Record<string, string> = {
    'pending': '等待中',
    'downloading': '下载中',
    'completed': '已完成',
    'failed': '失败',
    'cancelled': '已取消'
  };
  return statusMap[status] || status;
}

// 格式化时间
export function formatTime(timeStr: string): string {
  if (!timeStr) return '';
  const date = new Date(timeStr);
  return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
}

// 格式化进度信息
export function formatProgress(current: number, total: number): string {
  if (total <= 0) return '准备下载中...';
  const percentage = Math.round((current / total) * 100);
  return `${current}/${total} 张图片 (${percentage}%)`;
}

// 计算进度百分比
export function calculateProgressPercentage(current: number, total: number): number {
  if (total <= 0) return 0;
  return Math.round((current / total) * 100);
}

// 获取状态对应的CSS类名
export function getStatusClass(status: string): string {
  const statusClasses: Record<string, string> = {
    'pending': 'bg-blue-600 text-white',
    'downloading': 'bg-blue-600 text-white',
    'completed': 'bg-green-500 text-white',
    'failed': 'bg-red-600 text-white',
    'cancelled': 'bg-red-600 text-white'
  };
  return statusClasses[status] || 'bg-gray-600 text-white';
}