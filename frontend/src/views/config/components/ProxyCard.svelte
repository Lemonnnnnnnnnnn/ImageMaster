<script lang="ts">
  import { SetProxy, GetProxy } from '../../../../wailsjs/go/viewer/Viewer';
  
  export let proxyURL = '';
  export let loading = false;
  export let onError = (message: string) => {};
  export let onSuccess = (message: string) => {};

  async function loadProxySettings() {
    try {
      proxyURL = await GetProxy();
    } catch (err) {
      console.error('无法加载代理设置:', err);
      proxyURL = '';
    }
  }

  async function saveProxySettings() {
    loading = true;
    
    try {
      await SetProxy(proxyURL);
      onSuccess('成功保存代理设置');
    } catch (err: any) {
      onError(`保存代理设置失败: ${err.message || '未知错误'}`);
    } finally {
      loading = false;
    }
  }

  // 导出加载函数供父组件调用
  export { loadProxySettings };
</script>

<div class="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
  <div class="flex items-center p-3 bg-gray-50 border-b border-gray-200">
    <span class="text-xl mr-2">⚙️</span>
    <h2 class="m-0 text-base font-semibold text-gray-800">代理设置</h2>
  </div>
  <div class="p-4 flex-1 flex">
    <div class="w-full">
      <div class="flex gap-2 items-end">
        <div class="flex-1 mb-2">
          <label for="proxy" class="block mb-1 text-sm text-gray-600">代理服务器 URL</label>
          <input 
            type="text" 
            id="proxy" 
            bind:value={proxyURL} 
            placeholder="例如: http://127.0.0.1:7890"
            disabled={loading}
            class="w-full px-3 py-2 border border-gray-300 rounded text-sm transition-all duration-200 focus:border-blue-500 focus:shadow-md focus:outline-none placeholder-gray-400"
          />
        </div>
        <button on:click={saveProxySettings} disabled={loading} class="flex items-center justify-center px-4 py-2 bg-blue-600 text-white border-none rounded cursor-pointer text-sm font-medium transition-all duration-200 whitespace-nowrap hover:bg-blue-700 active:scale-95 disabled:bg-gray-400 disabled:cursor-not-allowed disabled:transform-none">
          <span class="mr-1">💾</span>
          <span>保存</span>
        </button>
      </div>
      <p class="text-xs text-gray-500 my-1 leading-relaxed">支持 HTTP 和 SOCKS 代理，格式为 http://host:port 或 socks5://host:port</p>
    </div>
  </div>
</div>