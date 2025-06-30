import { push, location } from 'svelte-spa-router';
import { get } from 'svelte/store';
import {
  GetMangaImages,
  GetAllMangas,
  DeleteManga,
  GetImageDataUrl
} from '../../../../../wailsjs/go/viewer/Viewer';
import { mangaStore, updateMangaStore, type MangaState } from '../stores/mangaStore';

export class MangaService {
  static async loadManga(path: string) {
    try {
      updateMangaStore({ loading: true });
      
      // 解码路径参数
      const mangaPath = decodeURIComponent(path);
      
      // 获取所有漫画以支持导航功能
      const mangas = await GetAllMangas();
      const currentMangaIndex = mangas.findIndex(m => m.path === mangaPath);
      
      let mangaName: string;
      if (currentMangaIndex >= 0) {
        mangaName = mangas[currentMangaIndex].name;
      } else {
        mangaName = mangaPath.split('/').pop() || '';
      }
      
      // 更新状态
      updateMangaStore({
        mangaPath,
        mangaName,
        mangas,
        currentMangaIndex,
        selectedImages: []
      });
      
      // 获取所有图片
      await this.loadImages(mangaPath);
      
    } catch (error) {
      console.error('加载漫画失败:', error);
    } finally {
      updateMangaStore({ loading: false });
    }
  }

  static async loadImages(mangaPath: string) {
    try {
      // 获取所有图片路径
      const imagePaths = await GetMangaImages(mangaPath);
      
      // 并行加载所有图片，保持顺序
      const imagePromises = imagePaths.map(async (imagePath) => {
        try {
          return await GetImageDataUrl(imagePath);
        } catch (error) {
          console.error(`加载图片失败: ${imagePath}`, error);
          return null;
        }
      });

      // 等待所有图片加载完成
      const loadedImages = await Promise.all(imagePromises);
      
      // 过滤掉加载失败的图片（null值）
      const selectedImages = loadedImages.filter(img => img !== null);
      updateMangaStore({ selectedImages });
    } catch (error) {
      console.error("获取图片路径失败:", error);
    }
  }

  static handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.backToHome();
    }
  }



  static backToHome() {
    push('/');
  }

  static toggleNavigation() {
    const state = get(mangaStore);
    updateMangaStore({ showNavigation: !state.showNavigation });
  }

  static navigateToNextManga() {
    const state = get(mangaStore);
    
    if (state.currentMangaIndex < state.mangas.length - 1) {
      const nextManga = state.mangas[state.currentMangaIndex + 1];
      const encodedPath = encodeURIComponent(nextManga.path);
      
      // 使用替代路由方案处理相同路径不同参数的导航
      const currentLocation = get(location);
      if (currentLocation.includes('/manga/')) {
        // 如果当前已经在漫画页面，采用直接加载新数据的方式
        this.loadManga(nextManga.path);
        
        // 更新 URL 但不触发导航事件
        window.history.pushState(null, '', `/#/manga/${encodedPath}`);
      } else {
        // 否则正常导航
        push(`/manga/${encodedPath}`);
      }
    }
  }

  static navigateToPrevManga() {
    const state = get(mangaStore);
    
    if (state.currentMangaIndex > 0) {
      const prevManga = state.mangas[state.currentMangaIndex - 1];
      const encodedPath = encodeURIComponent(prevManga.path);
      
      // 使用替代路由方案处理相同路径不同参数的导航
      const currentLocation = get(location);
      if (currentLocation.includes('/manga/')) {
        // 如果当前已经在漫画页面，采用直接加载新数据的方式
        this.loadManga(prevManga.path);
        
        // 更新 URL 但不触发导航事件
        window.history.pushState(null, '', `/#/manga/${encodedPath}`);
      } else {
        // 否则正常导航
        push(`/manga/${encodedPath}`);
      }
    }
  }

  static async deleteAndViewNextManga() {
    const state = get(mangaStore);
    
    if (state.currentMangaIndex >= 0 && confirm(`确定要删除 "${state.mangaName}" 并查看下一部漫画吗？`)) {
      updateMangaStore({ loading: true });
      
      try {
        // 记录下一个漫画的位置，因为删除后数组会变化
        const hasNextManga = state.currentMangaIndex < state.mangas.length - 1;
        const nextMangaPath = hasNextManga ? state.mangas[state.currentMangaIndex + 1].path : null;
        
        // 执行删除操作
        const success = await DeleteManga(state.mangaPath);
        
        if (success) {
          if (nextMangaPath) {
            // 重要：在导航前设置 loading 为 false，防止新页面保持加载状态
            updateMangaStore({ loading: false });
            
            // 使用替代路由方案处理相同路径不同参数的导航
            const encodedPath = encodeURIComponent(nextMangaPath);
            const currentLocation = get(location);
            if (currentLocation.includes('/manga/')) {
              // 如果当前已经在漫画页面，采用直接加载新数据的方式
              this.loadManga(nextMangaPath);
              
              // 更新 URL 但不触发导航事件
              window.history.pushState(null, '', `/#/manga/${encodedPath}`);
            } else {
              // 否则正常导航
              push(`/manga/${encodedPath}`);
            }
          } else {
            // 如果没有下一部漫画，返回首页
            push('/');
          }
        } else {
          alert('删除失败!');
          updateMangaStore({ loading: false });
        }
      } catch (error) {
        console.error('删除漫画失败:', error);
        updateMangaStore({ loading: false });
      }
    }
  }
}