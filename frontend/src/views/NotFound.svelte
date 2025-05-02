<script lang="ts">
  import { push } from 'svelte-spa-router';
</script>

<div class="not-found">
  <h1>页面未找到</h1>
  <p>抱歉，您请求的页面不存在。</p>
  <button on:click={() => push('/')}>返回首页</button>
</div>

<style>
  .not-found {
    text-align: center;
    padding: 2rem;
    margin: 2rem auto;
    max-width: 500px;
  }

  h1 {
    color: #e74c3c;
    margin-bottom: 1rem;
  }

  button {
    margin-top: 1rem;
    padding: 0.5rem 1rem;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
  }

  button:hover {
    background-color: #2980b9;
  }
</style> 