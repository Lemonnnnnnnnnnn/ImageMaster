<script lang="ts">
  import Home from './Home/index.svelte';
  
  // 移除Header导入，因为左侧导航栏已经提供导航功能
</script>

<!-- 直接显示内容，不需要额外的Header -->
<div class="min-h-screen bg-black-secondary">
  <slot>
    <!-- 默认显示 Home 组件 -->
    <Home />
  </slot>
</div>