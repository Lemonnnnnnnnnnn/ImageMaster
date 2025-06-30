import { writable } from 'svelte/store';

// 配置状态接口
export interface ConfigState {
  libraries: string[];
  outputDir: string;
  proxyURL: string;
  loading: boolean;
  error: string;
  success: string;
}

// 初始状态
const initialState: ConfigState = {
  libraries: [],
  outputDir: '',
  proxyURL: '',
  loading: false,
  error: '',
  success: ''
};

// 创建可写store
export const configStore = writable<ConfigState>(initialState);

// 辅助函数
export const configActions = {
  setLoading: (loading: boolean) => {
    configStore.update(state => ({ ...state, loading }));
  },
  
  setError: (error: string) => {
    configStore.update(state => ({ ...state, error, success: '' }));
  },
  
  setSuccess: (success: string) => {
    configStore.update(state => ({ ...state, success, error: '' }));
  },
  
  clearMessages: () => {
    configStore.update(state => ({ ...state, error: '', success: '' }));
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

// 成功消息自动清除
let successTimeout: NodeJS.Timeout;

export const showSuccessMessage = (message: string) => {
  configActions.setSuccess(message);
  if (successTimeout) clearTimeout(successTimeout);
  successTimeout = setTimeout(() => {
    configActions.setSuccess('');
  }, 3000);
};