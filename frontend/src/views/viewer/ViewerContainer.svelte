<script lang="ts">
  import Header from '../../components/Header.svelte';
  import Home from './Home/index.svelte';

  // 删除不再需要的导航函数，因为已由Header组件提供
</script>

<div class="w-full h-full flex flex-col">
  <Header title="漫画查看器" />
  
  <slot>
    <!-- 默认显示 Home 组件 -->
    <Home />
  </slot>
</div>