import { clsx } from "clsx";
import { BsHouseDoor } from "react-icons/bs";
import { useAppStore } from "@/store/app";
import SidebarItem from "./SidebarItem";
import { getMenuData, router } from "@/routes";

const Sidebar: React.FC = () => {

  const appMenus = getMenuData(router[0].children!);
  console.log(appMenus)

  const { isCollapsed } = useAppStore();

  console.log("渲染 - Sidebar")
 
  return (
    <div
      className={clsx(
        "layout-sidebar d-flex flex-column flex-shrink-0 bg-body-tertiary",
        isCollapsed ? "collapsed" : "normal"
      )}
    >
      <a
        href="/"
        className="logo d-flex border-bottom align-items-center mb-md-0 me-md-auto link-body-emphasis text-decoration-none"
      >
        <BsHouseDoor />
        {isCollapsed ? null : <span className="fs-4">LOGO</span>}
      </a>

      <ul className="nav nav-pills flex-column mb-auto">
        {appMenus.map((item, index) => (
          <SidebarItem
            item={item}
            level={1}
            key={index}
          />
        ))}
      </ul>
    </div>
  );
};


export default Sidebar;