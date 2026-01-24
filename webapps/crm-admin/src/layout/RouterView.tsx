import { Suspense } from "react";
import { Outlet } from "react-router-dom";

function RouterView() {
  return (
    <Suspense fallback={<div>加载中...</div>}>
      <Outlet />
    </Suspense>
  );
}

export default RouterView;
