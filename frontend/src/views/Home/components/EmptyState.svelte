<script lang="ts">
  import { Library, FolderOpen, Settings, Plus } from 'lucide-svelte';
  import { push } from 'svelte-spa-router';
  import Button from '../../../components/Button.svelte';
  
  export let type: 'no-libraries' | 'no-mangas' = 'no-libraries';
  
  // 导航函数
  function goToConfig() {
    push('/config');
  }
  
  // 空状态配置
  $: emptyConfig = {
    'no-libraries': {
      icon: Library,
      title: '尚未配置漫画库',
      description: '开始前，您需要添加至少一个漫画库目录',
      actionText: '前往设置',
      action: goToConfig,
      tips: [
        '漫画库是存储漫画文件的文件夹',
        '支持常见的图片格式（JPG、PNG、WEBP等）',
        '可以添加多个不同的漫画库目录'
      ]
    },
    'no-mangas': {
      icon: FolderOpen,
      title: '暂无漫画内容',
      description: '在已配置的漫画库中未找到任何漫画',
      actionText: '重新扫描',
      action: () => window.location.reload(),
      tips: [
        '确保漫画文件位于配置的库目录中',
        '检查文件格式是否被支持',
        '尝试刷新页面重新扫描'
      ]
    }
  }[type];
</script>

<!-- Fluent Design 空状态 -->
<div class="flex flex-col items-center justify-center py-20 px-6 text-center">
  <!-- 主要图标 -->
  <div class="w-24 h-24 bg-white-tertiary/10 rounded-fluent-xl flex items-center justify-center mb-8 relative">
    <!-- 背景光效 -->
    <div class="absolute inset-0 bg-gradient-to-br from-fluent-blue/20 to-transparent rounded-fluent-xl"></div>
    
    <!-- 图标 -->
    <svelte:component 
      this={emptyConfig.icon} 
      size={48} 
      class="text-white-tertiary relative z-10" 
    />
  </div>
  
  <!-- 标题和描述 -->
  <div class="max-w-md mb-8">
    <h2 class="text-2xl font-semibold text-white-primary mb-3">
      {emptyConfig.title}
    </h2>
    <p class="text-white-secondary leading-relaxed">
      {emptyConfig.description}
    </p>
  </div>
  
  <!-- 操作按钮 -->
  <div class="mb-8">
    <Button
      variant="filled"
      color="primary"
      onclick={emptyConfig.action}
      classes="px-6 py-3"
    >
      {#if type === 'no-libraries'}
        <Settings size={18} class="mr-2" />
      {:else}
        <Plus size={18} class="mr-2" />
      {/if}
      {emptyConfig.actionText}
    </Button>
  </div>
  
  <!-- 提示信息 -->
  <div class="max-w-lg">
    <div class="bg-glass-card/50 backdrop-blur-fluent-sm rounded-fluent-lg p-6 border border-white-tertiary/10">
      <h3 class="text-sm font-medium text-white-primary mb-3 flex items-center gap-2">
        <div class="w-2 h-2 bg-fluent-blue rounded-full"></div>
        小贴士
      </h3>
      <ul class="space-y-2 text-left">
        {#each emptyConfig.tips as tip}
          <li class="text-sm text-white-secondary flex items-start gap-2">
            <div class="w-1 h-1 bg-white-tertiary rounded-full mt-2 flex-shrink-0"></div>
            {tip}
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>