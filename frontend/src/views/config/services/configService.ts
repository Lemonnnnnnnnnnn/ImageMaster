import { 
  SelectLibrary,
  SetOutputDir
} from '../../../../wailsjs/go/library/API';

import {
  GetLibraries,
  GetOutputDir,
  SetProxy,
  GetProxy
} from '../../../../wailsjs/go/storage/API';

// 漫画库相关服务
export const libraryService = {
  async getLibraries(): Promise<string[]> {
    return await GetLibraries();
  },
  
  async selectLibrary(): Promise<string | null> {
    return await SelectLibrary();
  }
};

// 输出目录相关服务
export const outputService = {
  async getOutputDir(): Promise<string> {
    return await GetOutputDir();
  },
  
  async setOutputDir(): Promise<string | null> {
    return await SetOutputDir();
  }
};

// 代理设置相关服务
export const proxyService = {
  async getProxy(): Promise<string> {
    try {
      return await GetProxy();
    } catch (err) {
      console.error('无法加载代理设置:', err);
      return '';
    }
  },
  
  async setProxy(proxyURL: string): Promise<void> {
    await SetProxy(proxyURL);
  }
};

// 统一的配置服务
export const configService = {
  library: libraryService,
  output: outputService,
  proxy: proxyService
};