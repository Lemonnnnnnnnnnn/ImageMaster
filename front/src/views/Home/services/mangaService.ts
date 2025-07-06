import {
  GetAllMangas,
  DeleteManga,
  GetImageDataUrl,
  LoadAllLibraries,
} from '../../../../wailsjs/go/library/API';
import { GetLibraries } from '../../../../wailsjs/go/storage/API';
import { useHomeStore } from '../stores/homeStore';
// import { get } from 'svelte/store';
import type { Manga } from '../stores/homeStore';

export class MangaService {
  private homeStore: ReturnType<typeof useHomeStore>;
  constructor() {
    this.homeStore = useHomeStore();
  }
  /**
   * 初始化数据加载
   */
  async initialize(): Promise<void> {
    this.homeStore.loading = true;
    await LoadAllLibraries();
    await this.loadLibraries();
    await this.loadMangas();
    this.homeStore.loading = false;
  }

  /**
   * 加载库列表
   */
  async loadLibraries(): Promise<void> {
    const librariesData = await GetLibraries();
    this.homeStore.libraries = librariesData;
  }

  /**
   * 加载漫画列表
   */
  async loadMangas(): Promise<void> {
    this.homeStore.loading = true;
    const mangasData = await GetAllMangas();
    this.homeStore.mangas = mangasData;

    // 预加载每个漫画的预览图
    const imageCache = this.homeStore.mangaImages;
    for (let manga of mangasData) {
      if (!imageCache.has(manga.previewImg)) {
        const imageUrl = await GetImageDataUrl(manga.previewImg);
        imageCache.set(manga.previewImg, imageUrl);
      }
    }
    this.homeStore.mangaImages = imageCache;
    this.homeStore.loading = false;
  }

  /**
   * 删除漫画
   */
  async deleteManga(manga: Manga): Promise<boolean> {
    if (!confirm(`确定要删除 "${manga.name}" 吗？这将永久删除该文件夹及其内容！`)) {
      return false;
    }

    this.homeStore.loading = true;
    const success = await DeleteManga(manga.path);

    if (success) {
      const currentMangas = this.homeStore.mangas;
      this.homeStore.mangas = currentMangas.filter((m) => m.path !== manga.path);
    } else {
      alert('删除失败！');
    }

    this.homeStore.loading = false;
    return success;
  }

  /**
   * 获取漫画预览图
   */
  getMangaImage(previewImg: string): string {
    const imageCache = this.homeStore.mangaImages;
    return imageCache.get(previewImg) || '';
  }
}