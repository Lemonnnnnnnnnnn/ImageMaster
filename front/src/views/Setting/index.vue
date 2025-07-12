<template>
    <div class="flex flex-col gap-8 p-8 text-white">
        <div class="flex flex-col gap-4">
            <div class="text-xl">下载目录</div>
            <div class="flex gap-4">
                <Input class="flex-1 cursor-pointer" v-model="downloadDir" placeholder="请选择下载目录"
                    @click="changeOutputDir" />
            </div>
        </div>

        <div class="flex flex-col gap-4">
            <div class="text-xl">漫画库</div>
            <div class="flex gap-2 ">
                <div v-for="library in libraries" :key="library" class="flex items-center gap-2">
                    <button :class="{ 'bg-neutral-500/50': library === activeLibrary }"
                        class=" cursor-pointer hover:bg-neutral-500/50 rounded-2xl border-1 border-neutral-300/50 py-2 px-4"
                        @click="changeActiveLibrary(library)">{{ library }}</button>
                </div>
            </div>
            <div class="flex justify-end">
                <button class="rounded-2xl border-1 border-neutral-300/50 py-2 px-4" @click="addLibrary">添加漫画库</button>
            </div>
        </div>

        <div class="flex flex-col gap-4">
            <div class="text-xl">代理设置</div>
            <div class="flex gap-4">
                <Input @blur="saveProxy" class="flex-1" v-model="proxyUrl" placeholder="请输入代理地址" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Input } from '@/components';
import { onMounted, ref } from 'vue';
import { GetLibraries, GetOutputDir, GetProxy, SetProxy, GetActiveLibrary, SetOutputDir, SetActiveLibrary, AddLibrary } from '../../../wailsjs/go/config/API';
import { toast } from 'vue-sonner';
import { debounce } from '@/utils';
import { LoadLibrary } from '../../../wailsjs/go/library/API';

let proxyUrl = ref("")
let downloadDir = ref("")
let libraries = ref<string[]>([])
let activeLibrary = ref("")

const saveProxy = debounce((e: Event) => {
    SetProxy((e.target as HTMLInputElement).value).then(() => {
        refreshConfig()
    })
}, 1000)

async function refreshConfig() {
    proxyUrl.value = await GetProxy();
    downloadDir.value = await GetOutputDir();
    libraries.value = await GetLibraries();
    activeLibrary.value = await GetActiveLibrary();
}

async function changeOutputDir() {
    SetOutputDir().then(() => {
        toast.success("设置成功！")
        refreshConfig()
    })
}

async function changeActiveLibrary(library: string) {
    SetActiveLibrary(library).then(() => {
        toast.success("设置成功！")
        refreshConfig()
    })
}

async function addLibrary() {
    AddLibrary().then(() => {
        toast.success("添加成功！")
        refreshConfig().then(() => {
            LoadLibrary(activeLibrary.value).then(() => {
                toast.success("加载成功！")
            })
        })
    })
}

onMounted(async () => {
    refreshConfig()
})
</script>

<style scoped></style>