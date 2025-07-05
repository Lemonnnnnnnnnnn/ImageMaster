import { push } from 'svelte-spa-router';
import type { Manga } from '../stores/homeStore';

export class NavigationService {
  /**
   * 跳转到漫画查看页面
   */
  static viewManga(manga: Manga): void {
    // 将路径编码后传递给路由
    const encodedPath = encodeURIComponent(manga.path);
    push(`/manga/${encodedPath}`);
  }
}