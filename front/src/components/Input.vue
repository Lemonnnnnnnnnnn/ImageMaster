<template>
    <div>
        <div v-if="help" class="text-xs text-neutral-500 mb-1">{{ help }}</div>
        <input type="text" :value="inputValue" :placeholder="placeholder"
            class="w-full px-2 py-2 text-white text-sm outline-0 border-b-1 border-neutral-300 focus:border-b-blue-500 border-solid rounded-sm bg-neutral-700 focus:bg-neutral-800" />
    </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from "vue";

const props = defineProps({
    inputValue: { type: String, default: '' },
    help: { type: String, default: '' },
    placeholder: { type: String, default: '' },
});

const emits = defineEmits(['update:inputValue']);   
</script>

<style scoped></style>