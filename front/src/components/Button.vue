<template>
    <button :class="type === 'primary' ? 'bg-sky-800 hover:bg-sky-700' : 'bg-transparent hover:bg-neutral-700/50'"
        class="text-xs text-white px-4 py-2 rounded-md border  cursor-pointer border-neutral-700/50 transition-all duration-300 active:outline-1">
        <slot></slot>
    </button>
</template>

<script setup lang="ts">
const props = defineProps({
    type: {
        type: String,
        default: 'default'
    },
    size: {
        type: String,
        default: 'md'
    }
})
</script>

<style scoped></style>