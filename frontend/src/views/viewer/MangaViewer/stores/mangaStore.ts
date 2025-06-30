import { writable } from 'svelte/store';

export interface MangaState {
  mangaPath: string;
  mangaName: string;
  selectedImages: string[];
  currentImageIndex: number;
  loading: boolean;
  mangas: any[];
  currentMangaIndex: number;
  showNavigation: boolean;
  viewMode: 'single' | 'scroll';
}

const initialState: MangaState = {
  mangaPath: '',
  mangaName: '',
  selectedImages: [],
  currentImageIndex: 0,
  loading: true,
  mangas: [],
  currentMangaIndex: -1,
  showNavigation: false,
  viewMode: 'scroll'
};

export const mangaStore = writable<MangaState>(initialState);

// 便捷的更新函数
export const updateMangaStore = (updates: Partial<MangaState>) => {
  mangaStore.update(state => ({ ...state, ...updates }));
};

// 重置状态
export const resetMangaStore = () => {
  mangaStore.set(initialState);
};