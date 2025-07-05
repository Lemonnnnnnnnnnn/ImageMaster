import { wrap } from 'svelte-spa-router/wrap';

import Layout from './layout/index.svelte';
// 导入页面组件
import NotFound from './views/NotFound.svelte';

// // 导入漫画查看器组件
// import MangaViewer from './views/viewer/MangaViewer/index.svelte';
// import ViewerContainer from './views/viewer/ViewerContainer.svelte';

// // 导入下载器组件
// import Downloader from './views/downloader/Downloader.svelte';

// // 导入配置组件
// import Config from './views/config/Config.svelte';

// 定义路由
const routes = {
  // 主页 - 显示漫画查看器
  '/*': Layout,
  // // 漫画查看页面 - 带参数路由
  // '/manga/:path': wrap({
  //   component: MangaViewer,
  //   // 确保 path 参数能够正确传递，处理 URL 编码
  //   conditions: [
  //     (detail) => {
  //       try {
  //         return !!detail.params?.path;
  //       } catch (e) {
  //         return false;
  //       }
  //     }
  //   ]
  // }),

  // // 下载器页面
  // '/downloader': Downloader,

  // // 配置页面
  // '/config': Config,

  // 通配符路由 - 处理所有未匹配的 URL
  '*': NotFound,
};

export default routes;
