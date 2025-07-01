<script lang="ts">
  import { SetProxy, GetProxy } from '../../../../wailsjs/go/viewer/Viewer';
  import Button from '../../../components/Button.svelte';
  import Card from '../../../components/Card.svelte';
  
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

<Card classes="p-4 preset-outlined-surface-500 overflow-hidden flex flex-col">
  <header class="card-header flex items-center">
    <span class="text-xl mr-2">⚙️</span>
    <h2 class="text-base font-semibold">代理设置</h2>
  </header>
  <section class="p-4 flex-1 flex">
    <div class="w-full">
      <div class="flex gap-2 items-end">
        <div class="flex-1 mb-2">
          <label for="proxy" class="label">
            <span>代理服务器 URL</span>
          </label>
          <input 
            type="text" 
            id="proxy" 
            bind:value={proxyURL} 
            placeholder="例如: http://127.0.0.1:7890"
            disabled={loading}
            class="input"
          />
        </div>
        <Button 
          onclick={saveProxySettings} 
          disabled={loading}
          loading={loading}
          variant="filled"
          color="primary"
          classes="flex items-center"
        >
          <span class="mr-1">💾</span>
          <span>保存</span>
        </Button>
      </div>
      <p class="text-xs text-surface-500 my-1 leading-relaxed">支持 HTTP 和 SOCKS 代理，格式为 http://host:port 或 socks5://host:port</p>
    </div>
  </section>
</Card>