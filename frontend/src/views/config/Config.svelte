<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import { Folder, Globe, Library, Settings, Plus, Check } from 'lucide-svelte';
  import Header from "../../components/Header.svelte";
  import Button from "../../components/Button.svelte";
  import Card from "../../components/Card.svelte";
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

<!-- Fluent Design 黑色主题布局 -->
<div class="min-h-screen bg-black-secondary">
  <Header 
    title="应用设置" 
    subtitle="配置应用选项和管理漫画库"
  />
  
  <!-- 主要内容区域 -->
  <div class="container mx-auto max-w-4xl px-6 py-8 space-y-6">
    
    <!-- 输出目录设置卡片 -->
    <Card variant="elevated" padding="lg" classes="transition-fluent hover-lift">
      <div class="flex items-start gap-4">
        <!-- 图标 -->
        <div class="flex-shrink-0 w-12 h-12 bg-fluent-blue/20 rounded-fluent-lg flex items-center justify-center">
          <Folder size={24} class="text-fluent-blue" />
        </div>
        
        <!-- 内容区域 -->
        <div class="flex-1 min-w-0">
          <div class="mb-4">
            <h3 class="text-lg font-semibold text-white-primary mb-1">输出目录</h3>
            <p class="text-sm text-white-secondary">设置下载内容的保存位置</p>
          </div>
          
          <div class="flex items-center gap-4">
            <div class="flex-1">
              <div class="fluent-input bg-glass-input/50 px-4 py-3 rounded-fluent-md border border-white-tertiary/20">
                <span class="text-sm text-white-primary truncate">
                  {outputDir || "未设置输出目录"}
                </span>
              </div>
            </div>
            
            <Button
              onclick={changeOutputDir}
              disabled={loading}
              loading={loading}
              variant="filled"
              color="primary"
            >
              <Plus size={16} class="mr-2" />
              选择目录
            </Button>
          </div>
        </div>
      </div>
    </Card>

    <!-- 代理设置卡片 -->
    <Card variant="elevated" padding="lg" classes="transition-fluent hover-lift">
      <div class="flex items-start gap-4">
        <!-- 图标 -->
        <div class="flex-shrink-0 w-12 h-12 bg-fluent-green/20 rounded-fluent-lg flex items-center justify-center">
          <Globe size={24} class="text-fluent-green" />
        </div>
        
        <!-- 内容区域 -->
        <div class="flex-1 min-w-0">
          <div class="mb-4">
            <h3 class="text-lg font-semibold text-white-primary mb-1">代理设置</h3>
            <p class="text-sm text-white-secondary">配置网络代理以改善下载速度</p>
          </div>
          
          <div class="flex items-center gap-4">
            <div class="flex-1">
              <input
                type="text"
                bind:value={proxyURL}
                placeholder="例如：http://127.0.0.1:7890"
                disabled={loading}
                class="fluent-input w-full"
              />
            </div>
            
            <Button
              onclick={saveProxySettings}
              disabled={loading}
              loading={loading}
              variant="filled"
              color="success"
            >
              <Check size={16} class="mr-2" />
              保存设置
            </Button>
          </div>
        </div>
      </div>
    </Card>

    <!-- 漫画库管理卡片 -->
    <Card variant="elevated" padding="lg" classes="transition-fluent hover-lift">
      <div class="flex items-start gap-4">
        <!-- 图标 -->
        <div class="flex-shrink-0 w-12 h-12 bg-fluent-orange/20 rounded-fluent-lg flex items-center justify-center">
          <Library size={24} class="text-fluent-orange" />
        </div>
        
        <!-- 内容区域 -->
        <div class="flex-1 min-w-0">
          <div class="mb-4">
            <h3 class="text-lg font-semibold text-white-primary mb-1">漫画库管理</h3>
            <p class="text-sm text-white-secondary">添加和管理漫画存储目录</p>
          </div>
          
          {#if libraries.length === 0}
            <!-- 空状态 -->
            <div class="text-center py-8">
              <div class="w-16 h-16 mx-auto bg-white-tertiary/10 rounded-fluent-xl flex items-center justify-center mb-4">
                <Library size={32} class="text-white-tertiary" />
              </div>
              <h4 class="text-lg font-medium text-white-primary mb-2">尚未添加任何漫画库</h4>
              <p class="text-white-secondary mb-6">添加第一个漫画库来开始管理您的收藏</p>
              
              <Button
                onclick={addLibrary}
                disabled={loading}
                loading={loading}
                variant="filled"
                color="primary"
              >
                <Plus size={16} class="mr-2" />
                添加第一个漫画库
              </Button>
            </div>
          {:else}
            <!-- 漫画库列表 -->
            <div class="space-y-4">
              <div class="space-y-3">
                {#each libraries as lib, index}
                  <div 
                    class="flex items-center gap-3 p-3 bg-glass-card/30 rounded-fluent-md border border-white-tertiary/10 transition-fluent hover:bg-glass-card/50"
                    style="animation-delay: {index * 100}ms"
                  >
                    <div class="flex-shrink-0 w-8 h-8 bg-fluent-orange/20 rounded-fluent-md flex items-center justify-center">
                      <Folder size={16} class="text-fluent-orange" />
                    </div>
                    <span class="flex-1 text-sm text-white-primary truncate font-mono">
                      {lib}
                    </span>
                  </div>
                {/each}
              </div>
              
              <div class="flex justify-end pt-2">
                <Button
                  onclick={addLibrary}
                  disabled={loading}
                  loading={loading}
                  variant="outlined"
                  color="primary"
                >
                  <Plus size={16} class="mr-2" />
                  添加更多漫画库
                </Button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </Card>
    
    <!-- 提示信息 -->
    <Card variant="outlined" padding="md" classes="border border-fluent-blue/30 bg-fluent-blue/5">
      <div class="flex items-start gap-3">
        <div class="flex-shrink-0 w-5 h-5 bg-fluent-blue/20 rounded-full flex items-center justify-center mt-0.5">
          <Settings size={12} class="text-fluent-blue" />
        </div>
        <div class="flex-1">
          <p class="text-sm text-fluent-blue">
            <strong>提示:</strong> 设置完成后，请重启应用以确保所有配置生效。代理设置将影响图片下载速度，输出目录决定下载内容的保存位置。
          </p>
        </div>
      </div>
    </Card>
  </div>
</div>
