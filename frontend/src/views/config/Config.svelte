<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import Header from "../../components/Header.svelte";
  import Button from "../../components/Button.svelte";
  import { configStore, configActions } from "./stores/configStore";
  import { SelectLibrary, SetOutputDir } from "../../../wailsjs/go/library/API";
  import {
    GetLibraries,
    SetProxy,
    GetProxy,
    GetOutputDir,
  } from "../../../wailsjs/go/storage/API";

  // 响应式状态
  $: ({ libraries, outputDir, proxyURL, loading } = $configStore);

  onMount(async () => {
    await loadAllSettings();
  });

  async function loadAllSettings() {
    configActions.setLoading(true);

    try {
      // 并行加载所有设置
      const [librariesData, outputDirData] = await Promise.all([
        GetLibraries(),
        GetOutputDir(),
      ]);

      let proxyData = "";
      try {
        proxyData = await GetProxy();
      } catch (err) {
        console.error("无法加载代理设置:", err);
      }

      configActions.setLibraries(librariesData);
      configActions.setOutputDir(outputDirData);
      configActions.setProxyURL(proxyData);
    } catch (err: any) {
      toast.error(`加载配置失败: ${err.message || "未知错误"}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  // 输出目录相关
  async function changeOutputDir() {
    configActions.setLoading(true);

    try {
      const newDir = await SetOutputDir();
      if (newDir) {
        configActions.setOutputDir(newDir);
        toast.success("成功更改输出目录");
      }
    } catch (err: any) {
      toast.error(`更改输出目录失败: ${err.message || "未知错误"}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  // 代理设置相关
  async function saveProxySettings() {
    configActions.setLoading(true);

    try {
      await SetProxy(proxyURL);
      toast.success("成功保存代理设置");
    } catch (err: any) {
      toast.error(`保存代理设置失败: ${err.message || "未知错误"}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  // 漫画库相关
  async function loadLibraries() {
    configActions.setLoading(true);
    try {
      const libs = await GetLibraries();
      configActions.setLibraries(libs);
    } catch (err: any) {
      toast.error(`加载漫画库失败: ${err.message || "未知错误"}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  async function addLibrary() {
    configActions.setLoading(true);

    try {
      const newLib = await SelectLibrary();
      if (newLib) {
        await loadLibraries();
        toast.success("成功添加新的漫画库");
      }
    } catch (err: any) {
      toast.error(`添加漫画库失败: ${err.message || "未知错误"}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  function updateProxyURL(event: Event) {
    const target = event.target as HTMLInputElement;
    configActions.setProxyURL(target.value);
  }
</script>

<div class="bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/40">
  <div class="container mx-auto max-w-4xl px-6 py-8">
    <Header title="应用设置" />
    
    <!-- 主要内容区域 -->
    <div class="mt-8 space-y-8">
      
      <!-- 输出目录设置卡片 -->
      <div class="config-card group">
        <div class="config-card-header">
          <div class="flex items-center gap-3">
            <div class="config-icon-container">
              <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 1v6" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">输出目录</h3>
              <p class="text-sm text-gray-600 mt-1">设置下载内容的保存位置</p>
            </div>
          </div>
        </div>
        
        <div class="config-card-content">
          <div class="flex items-center gap-4">
            <div class="flex-1">
              <div class="fluent-input-display">
                <span class="text-sm text-gray-700 truncate">
                  {outputDir || "未设置输出目录"}
                </span>
              </div>
            </div>
            
            <Button
              onclick={changeOutputDir}
              disabled={loading}
              {loading}
              variant="filled"
              color="primary"
              classes="fluent-button-primary"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              选择目录
            </Button>
          </div>
        </div>
      </div>

      <!-- 代理设置卡片 -->
      <div class="config-card group">
        <div class="config-card-header">
          <div class="flex items-center gap-3">
            <div class="config-icon-container">
              <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">代理设置</h3>
              <p class="text-sm text-gray-600 mt-1">配置网络代理以改善下载速度</p>
            </div>
          </div>
        </div>
        
        <div class="config-card-content">
          <div class="flex items-center gap-4">
            <div class="flex-1">
              <div class="fluent-input-container">
                <input
                  type="text"
                  bind:value={proxyURL}
                  placeholder="例如：http://127.0.0.1:7890"
                  disabled={loading}
                  class="fluent-input"
                />
                <div class="fluent-input-border"></div>
              </div>
            </div>
            
            <Button
              onclick={saveProxySettings}
              disabled={loading}
              {loading}
              variant="filled"
              color="primary"
              classes="fluent-button-success"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              保存设置
            </Button>
          </div>
        </div>
      </div>

      <!-- 漫画库管理卡片 -->
      <div class="config-card group">
        <div class="config-card-header">
          <div class="flex items-center gap-3">
            <div class="config-icon-container">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">漫画库管理</h3>
              <p class="text-sm text-gray-600 mt-1">添加和管理漫画存储目录</p>
            </div>
          </div>
        </div>
        
        <div class="config-card-content">
          {#if libraries.length === 0}
            <div class="empty-state">
              <div class="empty-state-icon">
                <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h4 class="text-lg font-medium text-gray-900 mb-2">尚未添加任何漫画库</h4>
              <p class="text-gray-600 mb-6">添加第一个漫画库来开始管理您的收藏</p>
              
              <Button
                onclick={addLibrary}
                disabled={loading}
                {loading}
                variant="filled"
                color="primary"
                classes="fluent-button-primary"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                添加第一个漫画库
              </Button>
            </div>
          {:else}
            <div class="space-y-4">
              <div class="library-list">
                {#each libraries as lib, index}
                  <div class="library-item" style="animation-delay: {index * 50}ms">
                    <div class="library-item-icon">
                      <svg class="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
                      </svg>
                    </div>
                    <span class="library-item-path">{lib}</span>
                  </div>
                {/each}
              </div>
              
              <div class="flex justify-end pt-2">
                <Button
                  onclick={addLibrary}
                  disabled={loading}
                  {loading}
                  variant="outlined"
                  color="primary"
                  classes="fluent-button-outlined"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                  添加更多漫画库
                </Button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
