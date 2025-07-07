<template>
    <div class="flex flex-col gap-8 p-8">
        <div class="flex flex-col gap-4">
            <div class="text-xl">下载目录</div>
            <div class="flex gap-4">
                <Input class="flex-1 cursor-pointer" v-model="downloadDir" placeholder="请选择下载目录"
                    @click="changeOutputDir" />
                <!-- <Button @click="changeOutputDir">
                    <div class="flex items-center gap-2">
                        <Save :size="16" class="text-white" />
                        <span>保存目录</span>
                    </div>
                </Button> -->
            </div>
        </div>

        <div>
            <div class="text-xl">漫画库</div>
            <!-- <div class="flex gap-4">
                <Input class="flex-1 cursor-pointer" v-model="downloadDir" placeholder="请选择漫画库" @click="changeOutputDir" />
                <Button @click="changeOutputDir">
                    <div class="flex items-center gap-2">
                        <Save :size="16" class="text-white" />
                        <span>保存目录</span>
                    </div>
                </Button>
            </div> -->
        </div>

        <div class="flex flex-col gap-4">
            <div class="text-xl">代理设置</div>
            <div class="flex gap-4">
                <Input class="flex-1" v-model="proxyUrl" placeholder="请输入代理地址" />
                <Button @click="saveProxy">
                    <div class="flex items-center gap-2">
                        <Save :size="16" class="text-white" />
                        <span>保存设置</span>
                    </div>
                </Button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Input, Button } from '@/components';
import { onMounted, ref } from 'vue';
import { GetOutputDir, GetProxy, SetProxy } from '../../../wailsjs/go/storage/API';
import { Save } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import { SetOutputDir } from '../../../wailsjs/go/library/API';

let proxyUrl = ref("")
let downloadDir = ref("")
// let mangaDir = ref("")

function saveProxy() {
    SetProxy(proxyUrl.value).then(() => {
        toast.success("设置成功！")
        refreshConfig()
    })
}

async function refreshConfig() {
    proxyUrl.value = await GetProxy();
    downloadDir.value = await GetOutputDir();
}

async function changeOutputDir() {
    SetOutputDir().then(() => {
        toast.success("设置成功！")
        refreshConfig()
    })
}

onMounted(async () => {
    refreshConfig()
})
</script>

<style scoped></style>