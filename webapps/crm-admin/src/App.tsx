// src/App.tsx
import React from "react";
import {
  BrowserRouter as Router,
  useRoutes,
} from "react-router-dom";
import { router } from "@/routes"; // 导入路由配置

// 使用 useRoutes 渲染路由
const AppRoutes: React.FC = () => {
  const element = useRoutes(router);
  return element;
};

const App: React.FC = () => {
  return (
    // 使用 BrowserRouter
    <Router>
      <AppRoutes />
    </Router>
  );
};

export default App;