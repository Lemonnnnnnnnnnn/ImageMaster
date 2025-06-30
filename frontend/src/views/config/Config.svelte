<script lang="ts">
  import { onMount } from 'svelte';
  import Header from '../../components/Header.svelte';
  import NotificationArea from './components/NotificationArea.svelte';
  import LibrariesCard from './components/LibrariesCard.svelte';
  import OutputCard from './components/OutputCard.svelte';
  import ProxyCard from './components/ProxyCard.svelte';
  import { configStore, configActions, showSuccessMessage } from './stores/configStore';
  import { configService } from './services/configService';
  
  let librariesCardRef: LibrariesCard;
  let outputCardRef: OutputCard;
  let proxyCardRef: ProxyCard;

  // 响应式状态
  $: ({ libraries, outputDir, proxyURL, loading, error, success } = $configStore);

  onMount(async () => {
    await loadAllSettings();
  });

  async function loadAllSettings() {
    configActions.setLoading(true);
    
    try {
      // 并行加载所有设置
      const [librariesData, outputDirData, proxyData] = await Promise.all([
        configService.library.getLibraries(),
        configService.output.getOutputDir(),
        configService.proxy.getProxy()
      ]);
      
      configActions.setLibraries(librariesData);
      configActions.setOutputDir(outputDirData);
      configActions.setProxyURL(proxyData);
    } catch (err) {
      configActions.setError(`加载配置失败: ${err.message || '未知错误'}`);
    } finally {
      configActions.setLoading(false);
    }
  }

  function handleError(message: string) {
    configActions.setError(message);
  }

  function handleSuccess(message: string) {
    showSuccessMessage(message);
  }
</script>

<div class="config-container">
  <Header title="应用设置" />

  <NotificationArea {loading} {error} {success} />

  <div class="layout-grid">
    <LibrariesCard 
      bind:this={librariesCardRef}
      bind:libraries={$configStore.libraries}
      bind:loading={$configStore.loading}
      onError={handleError}
      onSuccess={handleSuccess}
    />
    
    <OutputCard 
      bind:this={outputCardRef}
      bind:outputDir={$configStore.outputDir}
      bind:loading={$configStore.loading}
      onError={handleError}
      onSuccess={handleSuccess}
    />
    
    <ProxyCard 
      bind:this={proxyCardRef}
      bind:proxyURL={$configStore.proxyURL}
      bind:loading={$configStore.loading}
      onError={handleError}
      onSuccess={handleSuccess}
    />
  </div>
</div>

<style>
  .config-container {
    padding: 16px;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .layout-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: 16px;
    grid-template-areas: 
      "libraries output"
      "libraries proxy";
  }
  
  /* 响应式布局调整 */
  @media (max-width: 768px) {
    .layout-grid {
      grid-template-columns: 1fr;
      grid-template-areas: 
        "libraries"
        "output"
        "proxy";
    }
  }
  
  @media (min-width: 1200px) {
    .layout-grid {
      grid-template-columns: 1fr 1fr 1fr;
      grid-template-areas: "libraries output proxy";
    }
  }
</style>