<script lang="ts">
  import Router from 'svelte-spa-router';
  import routes from './routes';
</script>

<main>
  <Router {routes} />
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  main {
    height: 100vh;
    /* width: 100vw; */
  }
</style>
