<script lang="ts">
  import Router from 'svelte-spa-router';
  import routes from './routes';
  import { Toaster } from 'svelte-sonner'
</script>

<main>
  <Router {routes} />
</main>
<Toaster />

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    background-color: #f5f5f5;
    color: #333;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  main {
    height: 100vh;
    width: 100%;
    /* overflow: hidden; */
    display: flex;
    flex-direction: column;
  }
</style>
