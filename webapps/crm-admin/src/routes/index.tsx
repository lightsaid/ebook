import type { RouteObject } from "react-router-dom";
import Workbench from "@/pages/Workbench";


// 定义路由表
const router: RouteObject[] = [
  {
    path: '/',
    element: <Workbench />, // 所有管理页面使用主布局
    children: [
      {
        index: true, // 默认子路由，对应 /
        element: <Workbench />,
      },
    ],
  },

];

export default router;