[TODO](https://juejin.cn/book/7294082310658326565?scrollMenuIndex=1)

是一本关于 react 组件的小册

## 1. 关于本小册

如何掌握好 React 呢？

我觉得就是这两方面：
一方面是 React 之上，学会写各种组件，并且能把这些组件封装成一个`组件库`、学习各种 React 相关的库。
一方面是 React 之下，能够自己调试源码，知道 React 是怎么运行的，能够实现一个`简易版 React`。

这两方面都掌握到一定程度，React 技术栈就算是通关了。

## 2. 一网打尽组件常用 Hook

- 官方文档也已经把 class 组件的语法划到了 legacy（遗产）的目录下。

- React.StrictMode 会导致额外的渲染
  在开发模式下，当组件内部使用了严格模式，`React 会故意将组件的挂载、更新和卸载操作执行两遍`。这样做的目的是为了帮助开发者发现那些可能不会在每次渲染中都表现出相同行为的副作用。如果一个副作用在两次渲染中表现不一致，那么它可能就是一个潜在的 bug 来源。
- 为什么要有 useLayoutEffect
  useEffect 的 effect 函数会在操作 dom 之后异步执行

  绝大多数情况下，用 useEffect，它能避免因为 effect 逻辑执行时间长导致页面卡顿（掉帧）。 但如果你遇到闪动的问题比较严重，那可以用 useLayoutEffect，但要注意 effect 逻辑不要执行时间太长。
  ![useLayoutEffect的执行是同步的，在重新渲染前](image.png)
  ![useEffect的执行是异步的，在重新渲染后](image-1.png)

  好处：浏览器会等 effect 逻辑执行完再渲染，好处自然就是不会闪动了。
  坏处：effect 逻辑要执行很久呢？就阻塞渲染了。
  useEffect 的 effect 函数是异步执行的，所以可能中间有次渲染，会闪屏，而 useLayoutEffect 则是同步执行的，所以不会闪屏，但如果计算量大可能会导致掉帧。

- 在 react 里，只要涉及到 state 的修改，就必须返回新的对象，不管是 useState 还是 useReducer
- react Context
  用 createContext 创建 context 对象，用 Provider 修改其中的值
  function 组件使用 useContext 的 hook 来取值，class 组件使用 Consumer 来取值。
- React.memo
  用 React.memo 的话，一般还会结合两个 hook：useMemo 和 useCallback。
  **React.memo 是防止 props 没变时的重新渲染，useMemo 和 useCallback 是防止 props 的不必要变化。**

  如果子组件用了 memo，那给它传递的对象、函数类的 props 就需要用 useMemo、useCallback 包裹，否则，每次 props 都会变，memo 就没用了。
  反之，如果 props 使用 useMemo、useCallback，但是子组件没有被 memo 包裹，那也没意义，因为不管 props 变没变都会重新渲染，只是做了无用功。

  **memo + useCallback、useMemo 是搭配着来的，少了任何一方，都会使优化失效。**

## 3. Hook 的闭包陷阱的成因和解决方案

- 闭包陷阱是什么：
  effect 函数等引用了 state，形成了闭包，但是并没有把 state 加到依赖数组里，导致执行 effect 时用的 state 还是之前的
- 本质原因：静态作用域
- 怎么办：

  1. 使用 `setState` 的函数的形式，从参数拿到上次的 state，这样就不会形成闭包了，或者用 useReducer，直接 dispatch action，而不是直接操作 state，这样也不会形成闭包
  2. 把`用到的 state 加到依赖数组里`，这样 state 变了就会重新跑 effect 函数，`引用新的 state`
  3. 使用 useRef 保存每次渲染的值，用到的时候从 `ref.current 取`

---

定时器的场景需要保证定时器只跑一次，不然重新跑会导致定时不准，所以需要用 useEffect + useRef 的方式来解决闭包陷阱问题。

我们还封装了 useInterval 的自定义 hook，这样可以不用在每个组件里都写一样的 useRef + useEffect 了，直接用这个自定义 hook 就行。

此外，关于要不要在渲染函数里直接修改 ref.current，其实都可以，直接改也行，包一层 useLayoutEffect 或者 useEffect 也行。

## 4. React 组件如何写 TypeScript 类型

- ReactNode > ReactElement > JSX.Element
- HTMLAttributes：组件可以传入 html 标签的属性，也可以指定具体的 ButtonHTMLAttributes、AnchorHTMLAttributes。

## 5. React 组件如何调试

## 6. 受控模式 VS 非受控模式

value 由用户控制就是非受控模式，由代码控制就是受控模式。

- 什么情况用受控模式
  需要对输入的值做处理之后设置到表单的时候，或者是你想实时同步状态值到父组件(比如把用户输入改为大写)
  ![ Form 组件内有一个 Store，会把表单值同步过去，然后集中管理和设置值](image-2.png)

- 非受控组件的 props 范式
  `defaultValue + onChange`
  这种情况，调用者只能设置 defaultValue 初始值，onChange 通知外部，组件内部的 state 值发生了变化。
- 受控组件的 props 范式
  `value + onChange`
  这种情况，调用者维护 value，onChange 通知外部需要改变 value。

一般的组件库，都会提供受控和非受控两种模式，比如 antd 的 Input 组件，就有 value 和 defaultValue 两个属性。
参数同时支持 value 和 defaultValue，`通过判断 value 是不是 undefined 来区分受控模式和非受控模式。`

- 抹平受控和非受控的差异的 hook：
  参见 ahooks 的 `useControllableValue` hook
  [useControllableValue](https://github.com/alibaba/hooks/blob/master/packages/hooks/src/useControllableValue/index.ts)
  用的时候就不用区分受控非受控了，直接 setState 就行

总结：
非受控模式就是完全用户自己修改 value，我们只是`设置个 defaultValue，可以通过 onChange 或者 ref 拿到表单值。`
受控模式是代码来控制 value，用户输入之后通过 `onChange 拿到值然后 setValue，触发重新渲染。`
单独用的组件，绝大多数情况下，用非受控模式就好了，因为你只是想获取到用户的输入。
如果需要结合 Form 表单用，那是要支持受控模式，因为 Form 会通过 Store 来统一管理所有表单项。
封装业务组件的话，用非受控模式或者受控都行。
有的团队就要求组件一定是受控的，然后在父组件里维护状态并同步到状态管理库，这样组件重新渲染也不会丢失数据。

# 7. 组件实战：迷你 Calendar

# 8. 组件实战：Calendar 日历组件(上)

# 9. 组件实战：Calendar 日历组件(下)

# 10. 快速掌握 Storybook

# 11 React 组件如何写单测？

# 12 深入理解 Suspense 和 ErrorBoundary

# 13 组件实战：Icon 图标组件

# 14 组件实战：Space 间距组件

# 15 React.Children 和它的两种替代方案

# 16 三个简单组件的封装

# 17 浏览器的 5 种 Observer

# 18 组件实战：Watermark 防删除水印组件

# 19 手写 react-lazyload

# 20 图解网页的各种距离

# 21 自定义 hook 练习

# 22 自定义 hook 练习(二)

# 23 react-spring 做弹簧动画

# 24 react-spring 结合 use-gesture 手势库实现交互动画

# 25 react-transition-group 和 react-spring 做过渡动画

# 26 快速掌握 tailwindcss

# 27 用 CSS Modules 避免样式冲突

# 28 CSS in JS: 快速掌握 styled-components

# 29 react-spring 实现滑入滑出的转场动画

# 30 组件实战：Message 全局提示组件

# 31 组件实战：Popover 气泡卡片组件

# 32 项目里如何快速定位组件源码？

# 33 一次超爽的 React 调试体验

# 34 组件实战：ColorPicker 颜色选择器（一）

# 35 组件实战：ColorPicker 颜色选择器（二）

# 36 组件实战：onBoarding 漫游式引导组件

# 37 组件实战：Upload 拖拽上传

# 38 组件实战：Form 表单组件

# 39 React 组件库都是怎么构建的

# 40 组件库实战：构建 esm 和 cjs 产物，发布到 npm

# 41 组件库实战：构建 umd 产物，通过 unpkg 访问

# 42 数据不可变：immutable 和 immer

# 43 基于 React Router 实现 keepalive

# 44 History api 和 React Router 实现原理

# 45 React Context 的实现原理和在 antd 里的应用

# 46 React Context 的性能缺点和解决方案

# 47 手写一个 Zustand

# 48 原子化状态管理库 Jotai

# 49 用 react-intl 实现国际化

# 50 国际化资源包如何通过 Excel 和 Google Sheet 分享给产品经理？

# 51 基于 react-dnd 实现拖拽排序

# 52 react-dnd 实战：拖拽版 TodoList

# 53 React Playground 项目实战：需求分析、实现原理

# 54 React Playground 项目实战：布局、代码编辑器

# 55 React Playground 项目实战：多文件切换

# 56 React Playground 项目实战：babel 编译、iframe 预览

# 57 React Playground 项目实战：文件增删改

# 58 React Playground 项目实战：错误显示、主题切换

# 59 React Playground 项目实战：链接分享、代码下载

# 60 React Playground 项目实战：Web Worker 性能优化

# 61 React Playground 项目实战：总结

# 62 手写 Mini React：思路分析

# 63 手写 Mini React：代码实现

# 64 手写 Mini React：和真实 React 源码的对比

# 65 React 18 的并发机制是怎么实现的？

# 66 Ref 的实现原理

# 67 低代码编辑器：核心数据结构、全局 store

# 68 低代码编辑器：拖拽组件到画布、拖拽编辑 json

# 69 低代码编辑器：画布区 hover 展示高亮框

# 70 低代码编辑器：画布区 click 展示编辑框

# 71 低代码编辑器：组件属性、样式编辑

# 72 低代码编辑器：预览、大纲

# 73 低代码编辑器：事件绑定

# 74 低代码编辑器：动作弹窗

# 75 低代码编辑器：自定义 JS

# 76 低代码编辑器：组件联动

# 77 低代码编辑器：拖拽优化、Table 组件

# 78 低代码编辑器：Form 组件、store 持久化

# 79 低代码编辑器：项目总结

# 80 快速掌握 React Flow 画流程图

# 81 React Flow 振荡器调音：项目介绍

# 82 React Flow 振荡器调音：流程图绘制

# 83 React Flow 振荡器调音：合成声音

# 84 AudioContext 实现在线钢琴

# 85 React 服务端渲染：从 SSR 到 hydrate

# 86 小册总结
