import { Navigate, type RouteObject } from "react-router-dom";
import Workbench from "@/pages/works/Workbench";
import MainLayout from "@/layout/MainLayout";
import BookList from "@/pages/books/BookList";
import AuthorList from "@/pages/books/AuthorList";
import CategoryList from "@/pages/books/CategoryList";
import PublisherList from "@/pages/books/PublisherList";
import OrderList from "@/pages/orders/OrderList";
import UserList from "@/pages/users/UserList";

import { BsWindowDesktop, BsBookHalf } from "react-icons/bs";

// 定义菜单类型，继承RouteObject类型
export type MenuItem = {
  name?: string; // 名称
  hidden?: boolean; // 是否在菜单上隐藏，默认都显示
  icon?: React.ComponentType<{ style?: React.CSSProperties }>; // icon 组件类型
  meta?: Record<string, BasicType>; // 元数据，比如code权限码
  children?: MenuItem[]; // 子菜单
} & RouteObject;

// 定义路由
export const router: MenuItem[] = [
  {
    path: "/",
    hidden: true,
    element: <MainLayout />, // 主布局
    children: [
      {
        name: "工作台",
        index: true, // 默认子路由，对应 /
        icon: BsWindowDesktop,
        element: <Workbench />,
      },
      {
        path: "books",
        name: "图书管理",
        icon: BsBookHalf,
        children: [
          {
            // 索引路由：负责重定向
            index: true,
            hidden: true,
            element: <Navigate to="/books/list" replace />,
          },
          {
            path: "list", // 自动拼接成"/books/list"
            name: "图书列表",
            element: <BookList />,
          },
          {
            path: "author",
            name: "作者列表",
            element: <AuthorList />,
          },
          {
            path: "category",
            name: "分类列表",
            element: <CategoryList />,
          },
          {
            path: "publisher",
            name: "出版社列表",
            element: <PublisherList />,
          },

          // 测试数据
          {
            path: "books2",
            name: "图书管理2",
            // icon: BsBookHalf,
            children: [
              {
                // 索引路由：负责重定向
                index: true,
                hidden: true,
                element: <Navigate to="/books/list" replace />,
              },
              {
                path: "list", // 自动拼接成"/books/list"
                name: "图书列表",
                element: <BookList />,
              },
              {
                path: "author",
                name: "作者列表",
                element: <AuthorList />,
              },
              {
                path: "category",
                name: "分类列表",
                element: <CategoryList />,
              },
              {
                path: "publisher",
                name: "出版社列表",
                element: <PublisherList />,
              },
            ],
          },
        ],
      },

      {
        path: "/orders",
        name: "订单管理",
        icon: BsWindowDesktop,
        children: [
          {
            index: true,
            hidden: true,
            element: <Navigate to="/orders/list" replace />,
          },
          {
            path: "list",
            name: "订单列表",
            element: <OrderList />,
          },
        ],
      },
      {
        path: "/users",
        name: "用户管理",
        icon: BsWindowDesktop,
        children: [
          {
            index: true,
            hidden: true,
            element: <Navigate to="/users/list" replace />,
          },
          {
            path: "list",
            name: "订单列表",
            element: <UserList />,
          },
        ],
      },
    ],
  },
];

/**
 * 根据路由定义生成菜单
 * 1. 过滤hidden的路由
 * 2. 过滤后，children仅有一个子项，组合新的菜单，保留一级即可
 */
export const getMenuData = (items: MenuItem[], parentPath = ""): MenuItem[] => {
  return items
    .filter((item) => !item.hidden) // 1. 过滤 hidden
    .map((item) => {
      // 解构出 index，防止它与 path 属性冲突
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { index, children: rawChildren, ...restItem } = item;

      // 处理当前项的完整路径
      const currentPath = item.index
        ? parentPath
        : resolvePath(parentPath, item.path);

      // 递归处理子项
      const filteredChildren = rawChildren
        ? getMenuData(rawChildren, currentPath)
        : undefined;

      // 2. 核心逻辑：如果子项过滤后仅剩一个有效项，则“提拔”该子项
      if (filteredChildren && filteredChildren.length === 1) {
        const onlyChild = filteredChildren[0];
        return {
          ...restItem, // 保留父级的 name, icon, meta 等
          path: onlyChild.path, // 使用子级的路径
          children: undefined, // 抹除子项，变为一级菜单
        } as MenuItem;
      }

      // 正常返回
      return {
        ...restItem,
        path: currentPath,
        children:
          filteredChildren && filteredChildren.length > 0
            ? filteredChildren
            : undefined,
      } as MenuItem;
    });
};

/**
 * 格式化路径工具函数
 */
const resolvePath = (base: string, path: string = "") => {
  if (path.startsWith("/")) return path;
  // 确保 base 不以 / 结尾，path 不以 / 开头，中间加 /
  const combined = `${base.endsWith("/") ? base.slice(0, -1) : base}/${
    path.startsWith("/") ? path.slice(1) : path
  }`;
  // 处理可能出现的双斜杠
  return combined.replace(/\/+/g, "/") || "/";
};
