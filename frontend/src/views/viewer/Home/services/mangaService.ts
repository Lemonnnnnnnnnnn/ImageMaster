import {
  GetAllMangas,
  DeleteManga,
  GetImageDataUrl,
  LoadAllLibraries,
} from '../../../../../wailsjs/go/library/API';
import { GetLibraries } from '../../../../../wailsjs/go/storage/API';
import { mangas, libraries, loading, mangaImages } from '../stores/homeStore';
import { get } from 'svelte/store';
import type { Manga } from '../stores/homeStore';

export class MangaService {
  /**
   * 初始化数据加载
   */
  static async initialize(): Promise<void> {
     loading.set(true);
    await LoadAllLibraries();
    await this.loadLibraries();
    await this.loadMangas();
    loading.set(false);
  }

  /**
   * 加载库列表
   */
  static async loadLibraries(): Promise<void> {
    const librariesData = await GetLibraries();
    libraries.set(librariesData);
  }

  /**
   * 加载漫画列表
   */
  static async loadMangas(): Promise<void> {
    loading.set(true);
    const mangasData = await GetAllMangas();
    mangas.set(mangasData);

    // 预加载每个漫画的预览图
    const imageCache = get(mangaImages);
    for (let manga of mangasData) {
      if (!imageCache.has(manga.previewImg)) {
        const imageUrl = await GetImageDataUrl(manga.previewImg);
        imageCache.set(manga.previewImg, imageUrl);
      }
    }
    mangaImages.set(imageCache);
    loading.set(false);
  }

  /**
   * 删除漫画
   */
  static async deleteManga(manga: Manga): Promise<boolean> {
    if (!confirm(`确定要删除 "${manga.name}" 吗？这将永久删除该文件夹及其内容！`)) {
      return false;
    }

    loading.set(true);
    const success = await DeleteManga(manga.path);
    
    if (success) {
      const currentMangas = get(mangas);
      mangas.set(currentMangas.filter((m) => m.path !== manga.path));
    } else {
      alert('删除失败！');
    }
    
    loading.set(false);
    return success;
  }

  /**
   * 获取漫画预览图
   */
  static getMangaImage(previewImg: string): string {
    const imageCache = get(mangaImages);
    return imageCache.get(previewImg) || '';
  }
}