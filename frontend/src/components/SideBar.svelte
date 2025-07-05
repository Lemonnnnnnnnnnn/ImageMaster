<script lang="ts">
    import { AlignJustify, Download, Home, Settings } from "lucide-svelte";
    import { fade, fly } from "svelte/transition";
    import { push } from "svelte-spa-router";
    import { location } from "svelte-spa-router";

    let collapsed = $state(false);

    let menuList = $derived([
        {
            icon: Home,
            label: "Home",
            path: "/",
            active: $location === "/",
        },
        {
            icon: Download,
            label: "Download",
            path: "/download",
            active: $location === "/download",
        },
        {
            icon: Settings,
            label: "Setting",
            path: "/setting",
            active: $location === "/setting",
        },
        
    ]);
</script>

<aside
    class="{collapsed
        ? 'w-12'
        : 'w-48'} border-r border-neutral-500/20 transition-all duration-200"
>
    <!-- <div class="justify-between flex items-center">
        {#if !collapsed}
            <div class="p-4 text-neutral-300 text-sm">Image Master</div>
        {/if} -->
    <div class="border-b border-neutral-500/20">
        <button
            class="p-4 hover:bg-neutral-700/50 transition-colors duration-200 cursor-pointer"
            onclick={() => (collapsed = !collapsed)}
        >
            <AlignJustify class="w-4 h-4 text-neutral-300" />
        </button>
    </div>
    <!-- </div> -->

    <div class="border-b border-neutral-500/20">
        {#each menuList as menu}
            <button
                class="relative w-full hover:bg-neutral-700/50 transition-colors duration-200 cursor-pointer {menu.active
                    ? 'bg-neutral-700/50'
                    : ''}"
                onclick={() => {
                    push(menu.path);
                }}
            >
                {#if menu.active}
                    <div
                        class="absolute left-0 top-0 w-0.5 h-full rounded-full bg-blue-500 transition-all duration-500 animate-shrink"
                    ></div>
                {/if}
                <div class="p-4 flex items-center gap-2">
                    <menu.icon class="w-4 h-4 text-neutral-300" />
                    {#if !collapsed}
                        <span
                            in:fly={{ x: -50 }}
                            class="text-neutral-300 text-xs ml-2"
                            >{menu.label}</span
                        >
                    {/if}
                </div>
            </button>
        {/each}
    </div>
</aside>

<style>

</style>
