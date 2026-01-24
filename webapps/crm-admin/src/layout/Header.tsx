import { useAppStore } from "@/store/app";
import { Navbar, Container, Nav, Dropdown } from "react-bootstrap";
import {  BsWindowDesktop, BsTextIndentLeft, BsTextIndentRight } from "react-icons/bs";

const Header: React.FC = () => {
  const { onMenuChange, isCollapsed } = useAppStore();
  return (
    <Navbar
      bg="white"
      variant="light"
      className="layout-header border-bottom"
      style={{ height: "56px" }}
    >
      <Container fluid>
        <Nav>
          {/* 菜单折叠按钮 */}
          <Nav.Link onClick={onMenuChange} style={{ fontSize: "20px" }}>
            { isCollapsed ? <BsTextIndentLeft /> :  <BsTextIndentRight />}
          </Nav.Link>
        </Nav>

        <Nav className="ms-auto me-sm-2">
          {/* 用户下拉菜单 */}
          <Dropdown as={Nav.Item} align="end">
            <Dropdown.Toggle
              as={Nav.Link}
              style={{ display: "flex", alignItems: "center" }}
            >
              <BsWindowDesktop
                style={{ fontSize: "20px", marginRight: "5px" }}
              />
              Admin
            </Dropdown.Toggle>
            <Dropdown.Menu>
              <Dropdown.Item href="#">个人中心</Dropdown.Item>
              <Dropdown.Divider />
              <Dropdown.Item href="#">退出登录</Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </Nav>
      </Container>
    </Navbar>
  );
};

export default Header;


