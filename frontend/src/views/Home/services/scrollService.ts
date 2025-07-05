import { showScrollTop, scrollY } from '../stores/homeStore';

export class ScrollService {
  private static scrollHandler: (() => void) | null = null;

  /**
   * 初始化滚动监听
   */
  static initScrollListener(): () => void {
    this.scrollHandler = this.handleScroll.bind(this);
    window.addEventListener('scroll', this.scrollHandler);
    
    // 返回清理函数
    return () => {
      if (this.scrollHandler) {
        window.removeEventListener('scroll', this.scrollHandler);
        this.scrollHandler = null;
      }
    };
  }

  /**
   * 处理滚动事件
   */
  private static handleScroll(): void {
    const currentScrollY = window.scrollY;
    scrollY.set(currentScrollY);
    showScrollTop.set(currentScrollY > 300);
  }

  /**
   * 滚动到顶部
   */
  static scrollToTop(): void {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}