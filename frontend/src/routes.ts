import { wrap } from 'svelte-spa-router/wrap';


// 导入漫画查看器组件
import MangaViewer from './views/MangaViewer/index.svelte';
import ViewerContainer from './views/ViewerContainer.svelte';

// 导入下载器组件
import Downloader from './views/Downloader/Downloader.svelte';

// 导入配置组件
import Config from './views/Config/Config.svelte';

// 定义路由
const routes = {
  // 主页 - 显示漫画查看器
  '/': wrap({
    component: ViewerContainer
  }),
  
  // 漫画查看页面 - 带参数路由
  '/manga/:path': wrap({
    component: MangaViewer,
    // 确保 path 参数能够正确传递，处理 URL 编码
    conditions: [
      (detail) => {
        try {
          return !!detail.params?.path;
        } catch (e) {
          return false;
        }
      }
    ]
  }),

  // 下载器页面
  '/downloader': Downloader,
  
  // 配置页面
  '/config': Config,
};

export default routes;