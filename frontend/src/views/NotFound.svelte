<script lang="ts">
  import { push } from 'svelte-spa-router';
</script>

<div class="text-center p-8 my-8 mx-auto max-w-lg">
  <h1 class="text-red-600 mb-4">页面未找到</h1>
  <p class="mb-4">抱歉，您请求的页面不存在。</p>
  <button on:click={() => push('/')} class="mt-4 px-4 py-2 bg-blue-500 text-white border-none rounded cursor-pointer text-base hover:bg-blue-600">返回首页</button>
</div>