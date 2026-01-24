// export {} 的作用不是“导出”，而是强制把这个 .d.ts 文件变成一个“模块”，
// 从而让 declare global 生效且不污染顶层作用域。
export {}; // 必须要写，不然别的地方使用不了declare global里的类型


// 定义全局类型，必须要在 declare global 里
declare global {
  type BasicType = string | number | boolean | null | undefined; 
}