import { writable } from 'svelte/store';

// 配置状态接口
export interface ConfigState {
  libraries: string[];
  outputDir: string;
  proxyURL: string;
  loading: boolean;
}

// 初始状态
const initialState: ConfigState = {
  libraries: [],
  outputDir: '',
  proxyURL: '',
  loading: false
};

// 创建可写store
export const configStore = writable<ConfigState>(initialState);

// 辅助函数
export const configActions = {
  setLoading: (loading: boolean) => {
    configStore.update(state => ({ ...state, loading }));
  },
  
  setLibraries: (libraries: string[]) => {
    configStore.update(state => ({ ...state, libraries }));
  },
  
  setOutputDir: (outputDir: string) => {
    configStore.update(state => ({ ...state, outputDir }));
  },
  
  setProxyURL: (proxyURL: string) => {
    configStore.update(state => ({ ...state, proxyURL }));
  }
};