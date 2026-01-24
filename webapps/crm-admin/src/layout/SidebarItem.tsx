import { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { Nav, Collapse, OverlayTrigger, Popover } from "react-bootstrap";
import { useAppStore } from "@/store/app";
import { clsx } from "clsx";
import { type MenuItem } from "@/routes";
import { FaChevronRight } from "react-icons/fa6";

const SidebarItem: React.FC<{
  item: MenuItem;
  level: number;
}> = ({ item, level }) => {
  const { isCollapsed } = useAppStore();
  const location = useLocation();

  // 判断当前路由是否匹配该菜单或其子菜单
  const isChildActive = (menu: MenuItem): boolean => {
    if (location.pathname === menu.path) return true;
    return menu.children?.some((child) => isChildActive(child)) ?? false;
  };

  const isActive = location.pathname === item.path;
  const hasActiveChild = item.children ? isChildActive(item) : false;

  const [open, setOpen] = useState(false);

  // 逻辑 1 & 2: 监听路由和折叠状态
  useEffect(() => {
    if (isCollapsed) {
      setOpen(false); // 折叠时，关闭所有二级菜单内容
    } else {
      // 展开时，如果子项被激活，则自动展开父级
      if (hasActiveChild) {
        setOpen(true);
      }
    }
  }, [isCollapsed, hasActiveChild]);

  // 处理点击
  const handleToggle = () => {
    if (isCollapsed) return;
    setOpen(!open);
  };

  // 逻辑 3: 折叠时的悬浮窗渲染
  const renderFloatingMenu = (
    <Popover id={`popover-${item.path}`} className="sidebar-popover">
      <Popover.Body className="p-0">
        <div className="popover-menu-title">{item.name}</div>
        {item.children?.map((child) => (
          <Link
            key={child.path}
            to={child.path}
            className={clsx(
              "popover-item",
              location.pathname === child.path && "active"
            )}
          >
            {child.name}
          </Link>
        ))}
      </Popover.Body>
    </Popover>
  );

  const content = (
    <Nav.Link
      className={clsx(
        "nav-link-item",
        (isActive || (hasActiveChild && isCollapsed)) && "active"
      )}
      as={item.children ? "div" : Link}
      to={item.children ? undefined : item.path}
      onClick={handleToggle}
    >
      <div className="icon-wrapper">{item.icon && <item.icon />}</div>
      {!isCollapsed && <span className="menu-text">{item.name}</span>}
      {/* 展开状态下的箭头指示 */}
      {!isCollapsed && item.children && (
        <span className={clsx("arrow-icon", open && "rotated")}>
          <FaChevronRight  />
        </span>
      )}
    </Nav.Link>
  );

  return (
    <li className={clsx("sidebar-item-container", level > 1 && "sub-item")}>
      {/* 如果折叠且有子菜单，使用 OverlayTrigger */}
      {isCollapsed && item.children ? (
        <OverlayTrigger
          placement="right"
          delay={{ show: 50, hide: 100 }}
          overlay={renderFloatingMenu}
        >
          {content}
        </OverlayTrigger>
      ) : (
        content
      )}

      {/* 侧边栏展开时的二级列表 */}
      {item.children && !isCollapsed && (
        <Collapse in={open}>
          <div>
            <ul className="nav flex-column sidebar-submenu">
              {item.children.map((child) => (
                <SidebarItem key={child.path} item={child} level={level + 1} />
              ))}
            </ul>
          </div>
        </Collapse>
      )}
    </li>
  );
};

export default SidebarItem;
