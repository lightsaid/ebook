import { clsx } from "clsx";
import { type MenuItem } from "@/routes";
import { Link, useLocation } from "react-router-dom";
import { Nav, Collapse, Overlay, Popover } from "react-bootstrap";


function FloatingMenu({item}: {item: MenuItem}) {
  return (
   <Popover id={`popover-${item.path}`} className="sidebar-popover">
      <Popover.Body className="p-0">
        <div className="popover-menu-title">{item.name}</div>
        {item.children?.map((child) => (
          <Link
            key={child.path}
            to={child.path!}
            className={clsx(
              "popover-item",
              location.pathname === child.path && "active"
            )}
          >
            {/* TODO: 要区分子级菜单 */}

            {child.name}
            
          </Link>
        ))}
      </Popover.Body>
    </Popover>
  )
}

export default FloatingMenu