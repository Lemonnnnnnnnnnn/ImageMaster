import { useHomeStore } from '../stores/homeStore';

export class ScrollService {
  private scrollHandler: (() => void) | null = null;
  private homeStore: ReturnType<typeof useHomeStore>;

  constructor() {
    this.homeStore = useHomeStore();
  }
  
  /**
   * 初始化滚动监听
   */
  initScrollListener(): () => void {
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
  private handleScroll(): void {
    const currentScrollY = window.scrollY;
    this.homeStore.scrollY = currentScrollY;
    this.homeStore.showScrollTop = currentScrollY > 300;
  }

  /**
   * 滚动到顶部
   */
  static scrollToTop(): void {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}