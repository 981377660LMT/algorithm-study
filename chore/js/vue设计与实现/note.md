从高层的设计角度，来探讨框架需要关注的问题

https://www.cnblogs.com/wenruo/p/17047716.html
https://juejin.cn/post/7109822222549663758#heading-17
https://juejin.cn/post/7197980894363156540?searchId=20240624231144C761F2A716049568C738
https://juejin.cn/post/7088894106113409032?searchId=20240624231144C761F2A716049568C738

## 框架设计概览

1. 权衡的艺术

   - 命令式 vs 声明式
     `命令式关注过程，声明式关注结果`
     `框架设计要做的就是，保持可维护性，同时让性能损失更少`
   - 虚拟 DOM
     性能：innerHTML < 虚拟 DOM < 原生 JavaScript
     心智负担/可维护性：虚拟 DOM < innerHTML < 原生 JavaScript
   - 运行时和编译时
     三种选择：纯运行时、运行时+编译时、纯编译时
     手写 vdom 太麻烦，所以通过模板 or JSX 完成书写，但是这些还需要经历一次编译
     Vue:运行时+编译时
     Svelte: 纯编译时

2. 框架设计的核心要素

   - 提升开发体验
     用户未按要求使用时的错误提示，console 自定义 formatter 直观显示数据
   - 控制框架代码的体积
     开发环境提供良好提示，不增加生产环境体积。（通过 `__DEV__` 变量判断）
   - Tree-Shaking
     如果一个函数有副作用，将不会被 Tree-Shaking 删除，通过 `/#__PURE__*/` 在打包工作（ rollup/webpack ）中声明没有副作用。
   - 输出产物
     IIFE：rollup format: 'iife'
     ESM: format: 'esm'
     CommonJS: format: 'cjs'
   - 特性开关(类似 rust 的 feature flag)
     用户关闭的特性，利用 Tree-Shaking 不打包到最终资源里。比如在 Vue3 中，对于 options API
   - 处理错误
     执行用户提供的函数时，做统一的错误处理（try...catch），并提供给用户错误接口来处理。（可以错误上报等。）
   - TypeScript 类型支持

3. Vue.js 3 的设计思路

   - 声明式的描述 UI
     模板描述，或者虚拟 DOM 描述
   - 渲染器
     将`虚拟` DOM 渲染成`真实` DOM
   - 组件的本质
   - 模板的工作原理
   - Vue 是各模块组成的有机体
     模板 --[编译器]--> 渲染函数 --[渲染器]--> 真实 DOM
     如果在编译器中增加优化，比如静态节点的判断，就可以在渲染器减少一些工作。

## 响应式系统

4.  响应系统的作用与实现

- 响应式数据与副作用
  副作用函数就是会对外部造成影响的函数，比如修改了全局变量。
  响应式：修改了某个值的时候，某个会`读取该值的副作用函数能够自动重新执行`。
- 响应系统的简单实现
  1、副作用读取值的时候，把函数放到值的某个桶里
  2、重新给值赋值的时候，执行桶里的函数
  在 Vue2 中通过 Object.defineProperty 实现，Vue3 通过 Proxy 实现。
- 设计一个完善的响应系统
  - `在全局维护一个变量来存储这个副作用函数`
  - WeakMap 来存储对象持有的副作用函数
  - 把 get 和 set 中的操作分别封装到 track 和 trigger
  - 切换分支与 cleanup：每次副作用函数重新执行的时候，我们要先把它从所有与之关联的依赖集合中删除。执行后会建立新的关联。
  - `嵌套的 effect 与 effect 栈`
    如果有嵌套的 effect 执行，我们就需要在保存当前 effect 函数的同时，记录之前的 effect 函数，并在当前的函数之前完之后，把上一层的 effect 赋值为 activeEffect。很简单的会想到用栈来实现这个功能。

5. 对象的响应性实现原理
6. 非对象的响应性实现原理

**Q & A**

- 设计存储空间的数据结构

  ```ts
  WeakMap<object, Map<PropertyKey, Set<() => void>>>
  ```

- 分支切换导致冗余依赖的问题
  每次副作用函数重新执行的时候，我们要先把它从所有与之关联的依赖集合中删除。执行后会建立新的关联

- 嵌套副作用函数的问题
  `执行 fn 前入栈，执行 fn 后出栈`

  ```ts
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    effectStack.push(effectFn)
    fn()
    effectStack.pop()
    activeEffect = effectStack.length > 0 ? effectStack[effectStack.length - 1] : undefined
  }
  ```

- 副作用函数无限调用自身导致栈溢出
  在 trigger 中执行副作用函数的时候，`不执行当前正在处理的副作用函数`，即 activeEffect
- 响应式系统的可调度性
  1. 把 options 挂在 effectFn 上
  2. 如果一个副作用函数存在调度器，就用调度器执行副作用函数
- 过期的副作用(避免竞态问题)
  第一次每次执行回调的请求之前给 watch 传一个过期函数，然后 watch 把它保存起来，然后在这个过程中`如果再次执行 watch 了，就会执行之前保存的过期函数，就会把上次的请求设置为不合法`

## 渲染器(Renderer)

将虚拟 DOM 渲染成真实 DOM

7. 渲染器的设计
8. DOM 的挂载和更新的逻辑
9. Diff 算法

## 组件化

10. 组件的实现原理

- setup 函数
  返回一个函数，该函数将作为组件的 render 函数
  返回一个对象，其中包含的数据暴露给模板使用

11. 异步组件与函数式组件
12. 内建组件和模块
    KeepAlive、Teleport、Transition

- keepalive 的本质是缓存管理，在加上特殊的挂载/卸载逻辑。
  挂载一个被 keepalive 的组件时，它并不会真的被卸载，而会被移动到一个隐藏容器中。
  当重新”挂载“该组件时，它也不会被真的挂载，而是被从隐藏容器中取出，再搬运到原来的容器中。
  这个过程对应到组件的 activated 和 deactivated 生命周期中。
- teleport 的本质是将组件的 DOM 结构移动到另一个地方。
  该组件会直接把它的插槽内容渲染到 body 下，而不会按照模板的 dom 层级渲染，这样就实现了跨层级渲染。
- transition
  过渡组件的实现原理是：
  当 dom 元素被挂载时，将动效附加到该 dom 元素上
  当 dom 元素被卸载时，不要立即卸载 dom 元素，而是等到附加到该 dom 元素上的动效执行完成后再卸载它

## 编译器

将模板编译成渲染函数

13. 编译器核心技术概述
14. 解析器
15. 编译优化
    `做静态分析，尽可能的在编译阶段就知道哪些东西会变，哪些东西不会变。`
    `于是就可以对不变的东西做下缓存，不需要更新的不更新`

- Block 和 PatchFlags
- 静态提升
- 预字符串化
- 缓存内联事件处理函数

  ```js
  function render(ctx, cache) {
    // cache 数组来自组件实例
    return (
      openBlock(),
      {
        // 将内联事件处理函数缓存到 cache 数组中
        onChange: cache[0] || (cache[0] = $event => ctx.a + ctx.b)
      }
    )
  }
  ```

- v-once 指令

## 服务端渲染

16. 同构渲染
