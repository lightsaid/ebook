import { useAppStore } from "@/store/app";
import "./layout.scss";
import Header from "./Header";
import Sidebar from "./Sidebar";
import RouterView from "./RouterView";


// TODO:
// 1. 通过浏览器地址输入的链接，定位到菜单并展开
// 2. 展开/收缩功能，二级菜单需要联动起来，当从收缩点击展开时保留当前导航的菜单展开，其他菜单折叠
//    当从展开点击收缩时，所有二级菜单折叠起来，仅保留一级菜单的图表；这里一切需要动画效果，过度丝毫。
// 3. 折叠时，仅显示一级菜单的icon，鼠标移动到icon上，需要显示二级菜单
// 4. Tab 支持多页，是否要缓存状态？

export default function MainLayout() {
  const { isCollapsed } = useAppStore();
  console.log("渲染 - MainLayout")

  return (
    <div className="layout-container container-fluid">
      {/* 侧边栏，固定定位 */}
      <Sidebar />

      {/* 主内容区域，根据侧边栏状态调整左侧 margin */}
      <div className={`layout-body ${isCollapsed ? "collapsed" : ""}`}>
        {/* 顶部导航 */}
        <Header />

        {/* 页面内容 */}
        <main className="layout-main">
          {/* React Router 的内容出口，使用 Suspense 处理懒加载 */}
          <RouterView />
        </main>
      </div>
    </div>
  );
}

