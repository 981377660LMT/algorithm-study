## V2 to V3

> 可以看到优秀的设计理念、权衡

https://ahooks.js.org/zh-CN/guide/upgrade

- 所有的输入函数，永远使用最新的一份。
- 所有 Hooks 返回的函数，也有和 setState 一样的特性，地址不会变化。
- DOM 类 Hooks 支持 target 动态变化

### 核心升级理念

1. **简单性**

   - `核心代码极简，采用插件式架构`
     新版 useRequest 只做 Promise 管理的底层能力，更多高级能力可以基于 useRequest 封装高级 Hooks 来支持
   - `删除复杂功能`（如 pagination、loadMore），让基础 hooks 更纯粹
   - 文档循序渐进，便于学习

2. **更好的类型支持**

   - 删除 TypeScript 重载，使类型提示更清晰
   - API 设计更符合 TypeScript 最佳实践

3. **更可靠的运行时**

   - 全面支持 SSR
   - 解决严格模式问题
   - 修复热更新相关问题
   - 优化闭包问题

4. **更高的扩展性**
   - 基础 hooks 更专注
   - 鼓励用户基于基础 hooks 封装高级功能，不默认集成(请求库 axios)

## ahooks 函数处理规范

ahooks 通过对输入输出函数做特殊处理，尽力帮助大家避免闭包问题。
输入函数，我们通过 useRef 做一次记录，以保证在任何地方都能访问到最新的函数。
针对输出函数，我们通过 ahooks 的 useMemoizedFn 包裹，保证地址永远不会变化。

## React Hooks & SSR

由于 SSR 是在非浏览器环境执行 JS 代码，所以会出现很多问题。本文主要介绍 React Hooks 在 SSR 模式下常见问题及解决方案。

- 不要在组件初始化时直接用 window/document 等浏览器对象
  要在 useEffect 中使用，服务端不会执行

- 判断环境，安全访问 DOM/BOM

  ```tsx
  import React, { useState } from 'react'

  function isBrowser() {
    return !!(typeof window !== 'undefined' && window.document && window.document.createElement)
  }

  export default () => {
    const [state, setState] = useState(isBrowser() && document.visibilityState)

    return state
  }
  ```

- target 属性要支持函数写法
  ahooks 的 DOM 类 hooks，target 要用函数返回 DOM 节点
- LayoutEffect 用同构 `useIsomorphicLayoutEffect`

  ```tsx
  const useIsoLayoutEffect = isBrowser() ? useLayoutEffect : useEffect
  ```

## React Hooks & react-refresh（HMR）

https://ahooks.js.org/zh-CN/guide/blog/hmr

### 什么是 react-refresh（HMR）

- **react-refresh** 是 React 官方的“模块热替换（Hot Module Replacement, HMR）”方案，开发时修改代码可以**保留组件状态**，只更新变动部分，极大提升开发体验。
- 对于**函数组件**，react-refresh 会保留 useState/useRef 的值；对于**类组件**，会完全 remount（状态会丢失）。

---

### react-refresh 的工作机制

- **热更新时**，useState/useRef 的值不会变（状态保留），但 useEffect/useCallback/useMemo 等副作用会重新执行。
- 这样可以让你在开发时修改代码，页面状态不丢失，但副作用逻辑会刷新。

---

### react-refresh 导致的典型问题

1. **useEffect 重复执行，状态递增**

   ```jsx
   const [count, setCount] = useState(0)
   useEffect(() => {
     setCount(c => c + 1)
   }, [])
   // 热更新后 count 会不断递增，因为 effect 会重新执行，但 state 不变
   ```

2. **useUpdateEffect 等自定义 Hook 行为异常**

   - 依赖 ref 判断“是否首次执行”，但热更新时 ref 不会重置，导致副作用被错误触发。

3. **isUnmount 等 ref 标记失效**
   - 组件热更新后，ref 还保留旧值，导致逻辑异常（如 loading 状态一直为 true）。

---

### 解决方案

- **方案一：代码层面重置 ref**
  - `在 useEffect 里初始化 ref，确保热更新后 ref 状态正确`。
- **方案二：加注释强制 remount**
  - 在文件顶部加 `/* @refresh reset */`，每次热更新都重新挂载组件，彻底重置所有状态。

---

### 记忆口诀

> **“热更新保状态，副作用会重跑，ref 要重置，@refresh reset 可全清。”**

---

### 官方态度

- 这些行为是预期的，开发时要注意 Hooks 的副作用和 ref 的持久性，必要时用官方推荐的 reset 注释。

---

**总结：**  
react-refresh 极大提升了开发体验，但也带来了副作用和 ref 的“潜规则”。写 Hooks 时要特别注意热更新下的状态和副作用一致性，必要时用 `/* @refresh reset */` 强制重置。

## React Hooks & strict mode

https://ahooks.js.org/zh-CN/guide/blog/strict

React 的严格模式是一种开发时的检查工具，用于帮助发现不安全的生命周期、过时 API 和副作用等潜在问题。
通过 <React.StrictMode> 包裹组件即可开启，`所有不建议的 API 或写法，都会抛出警告`（只在开发模式生效）。
严格模式很重要的一个能力是`检测意外的副作用`。

concurrent 模式中，组件被分为两个阶段：

- 渲染（render）阶段：生成 DOM 树，会执行 `constructor`、componentWillMount、componentWillReceiveProps、componentWillUpdate、getDerivedStateFromProps、shouldComponentUpdate、render、`useState`、`useMemo`、`useCallback` 等生命周期

- 提交（commit）阶段：操作 DOM，触发 componentDidMount、componentDidUpdate、useEffect 等生命周期

**渲染阶段可能会被暂停、重新执行。**

```tsx
// 在 constructor 中发起网络请求，就可能被执行多次。所以不要在渲染阶段执行带有副作用的操作。
constructor(){
  services.getUserInfo().then(() => {
    .....
  });
}
```

假如你在渲染阶段执行了副作用操作，React 也是无法感知的。
**但是 React 在严格模式下，会故意重复执行渲染阶段的方法，使得我们在开发阶段能更容易发现这类 bug**。
**记住结论：在严格模式下，useState、useMemo、useReducer 的第一个参数、Hook 函数体都会被执行两次，不要在这里执行带有副作用的操作，`副作用只放在 useEffect/useLayoutEffect 里`。**

```tsx
const useTest = () => {
  const [state, setState] = useState(() => {
    console.log('get state') // 会执行两次
    return 'state'
  })

  const memoState = useMemo(() => {
    console.log('get memo state') // 会执行两次
    return 'state'
  }, [])

  console.log('render') // 会执行两次

  return state
}
```
