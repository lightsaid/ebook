# crm-admin

React + TypeScript + Vite + react-bootstrap + bootstrap5 + zustand + react-router

### 创建项目

```bash
npm create vite@v6.5.0 crm-admin -- --template react-ts
```

因为 react-bootstrap 没有明确指出支持 react19，所以这里降级 react。

```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.66",
    "@types/react-dom": "^18.2.22"
  }
}
```

安装常用依赖库

```bash
npm install
npm ls react react-dom # 检查一下
npm i react-router react-router-dom
npm install react-bootstrap bootstrap
npm install zustand
npm install react-icons --save
npm install sass@1.78.0 # 安装和bootstrap兼容的版本
npm install clsx
```

### 配置别名

- npm i @types/node -D

vite.config.ts

```ts
 resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	}
```

tsconfig.app.json

```json
{
  "compilerOptions": {
    // ...
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

### 如何实现多 Tab 页，页面不刷新

考量：

在 Vue 里有`<keep-alive><router-view /></keep-alive>`可以实现这个功能；

但是 React 路由切换 -> 组件卸载 -> 再回来 -> 重新挂载，那么状态滚动条全无了。

解决方案思考：

1. 寻找一个类似 Vue 的 KeepAlive 的第三方库，如：react-activation

2. 用`display:none`控制显示，感觉复杂度较高

3. 使用状态管理如 zustand 保存页面状态，页面恢复时重新设置，根据业务复杂而定，非常麻烦

暂不做。

### 笔记

- src/routes/index.tsx 不想eslint报变量未使用错误

```ts
// 解构出 index，防止它与 path 属性冲突
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const { index, children: rawChildren, ...restItem } = item;
```

- src/layout/SidebarItem.tsx 下面的to属性，用法正确，但是ts无法正确预判，所以阻止误报 

```tsx
<Nav.Link
  className={clsx(
    "nav-link-item",
    (isActive || (hasActiveChild && isCollapsed)) && "active"
  )}
  as={item.children ? "div" : Link}
  // 停止to属性误报
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  to={item.children ? undefined : item.path}
  onClick={handleToggle}
>
  <div className="icon-wrapper">{item.icon && <item.icon />}</div>
  {!isCollapsed && <span className="menu-text">{item.name}</span>}
  {/* 展开状态下的箭头指示 */}
  {!isCollapsed && item.children && (
    <span className={clsx("arrow-icon", open && "rotated")}>
      <FaChevronRight size={15} />
    </span>
  )}
</Nav.Link>
```
