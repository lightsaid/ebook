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
