<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { loading, libraries, mangas } from "./stores/homeStore";
    import { MangaService } from "./services/mangaService";
    import { ScrollService } from "./services/scrollService";
    import { Loading } from "../../components";
    import { EmptyState, MangaGrid } from "./components";

    let cleanupScrollListener: (() => void) | null = null;

    onMount(async () => {
        // 初始化数据加载
        await MangaService.initialize();

        // 初始化滚动监听
        cleanupScrollListener = ScrollService.initScrollListener();
    });

    onDestroy(() => {
        // 清理滚动监听器
        if (cleanupScrollListener) {
            cleanupScrollListener();
        }
    });
</script>

<main>
    {#if $libraries.length === 0 && !$loading}
        <EmptyState type="no-libraries" />
    {/if}

    {#if $loading}
        <Loading />
    {:else if $mangas.length > 0}
        <MangaGrid />
    {:else if $libraries.length > 0 && !$loading}
        <EmptyState type="no-mangas" />
    {/if}
</main>
