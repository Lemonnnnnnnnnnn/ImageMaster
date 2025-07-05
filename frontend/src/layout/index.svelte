<script lang="ts">
    import { SideBar } from "../components";
    import { Home, MangaDetail, Setting } from "../views";
    import Downloader2 from "../views/Downloader2/index.svelte";
    import Router from "svelte-spa-router";

    const routes = {
        "/": Home,
        "/download": Downloader2,
        "/setting": Setting,
        "/manga/:path": MangaDetail,
    };
</script>

<div class="flex h-screen bg-neutral-900">
    <SideBar />
    <main class="flex-1">
        <slot>
            <Router {routes} />
        </slot>
    </main>
</div>
