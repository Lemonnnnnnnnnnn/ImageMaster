import { writable } from 'svelte/store';

export interface Manga {
  name: string;
  path: string;
  previewImg: string;
  imagesCount: number;
}

export interface Library {
  // 根据实际Library结构定义
  [key: string]: any;
}

// 状态管理
export const mangas = writable<Manga[]>([]);
export const libraries = writable<string[]>([]);
export const loading = writable<boolean>(true);
export const mangaImages = writable<Map<string, string>>(new Map());

// 滚动相关状态
export const showScrollTop = writable<boolean>(false);
export const scrollY = writable<number>(0);