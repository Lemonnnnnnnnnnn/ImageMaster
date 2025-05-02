import { wrap } from 'svelte-spa-router/wrap';

// 导入页面组件
import Home from './views/Home.svelte';
import MangaViewer from './views/MangaViewer.svelte';
import NotFound from './views/NotFound.svelte';

// 定义路由
const routes = {
  // 主页 - 显示所有漫画
  '/': Home,
  
  // 漫画查看页面 - 带参数路由
  '/manga/:path': wrap({
    component: MangaViewer,
    // 确保 path 参数能够正确传递，处理 URL 编码
    conditions: [
      (detail) => {
        try {
          return !!detail.params.path;
        } catch (e) {
          return false;
        }
      }
    ]
  }),

  // 通配符路由 - 处理所有未匹配的 URL
  '*': NotFound,
};

export default routes; 