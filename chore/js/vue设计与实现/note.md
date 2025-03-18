从高层的设计角度，来探讨框架需要关注的问题

https://www.cnblogs.com/wenruo/p/17047716.html

https://juejin.cn/post/7109822222549663758#heading-17

https://juejin.cn/post/7197980894363156540?searchId=20240624231144C761F2A716049568C738

https://juejin.cn/post/7088894106113409032?searchId=20240624231144C761F2A716049568C738

- 响应系统：监听可变数据，数据变化时触发回调函数
- 渲染器：将 VDOM 挂载或更新为真实 DOM ，其中涉及到 diff 算法
- 组件化：支持把一个大型系统拆分为若干组件，形成组件树
- 编译器：把 Vue 模板编译为 JS 代码 （对应 React 中的 JSX）
- **Reflect API 的作用就是：能改变对象 getter 里的 this(receiver)**
- **渲染时高效切换 DOM 事件**
  Vue 模板中绑定了事件，那渲染为真实 DOM 也需要绑定 DOM 事件。
  如果事件更新了，按照一般的思路是先 removeEventListener 然后再 addEventListener ，就是两次 DOM 操作 —— DOM 操作是昂贵的。
  Vue 对此进行了优化，极大减少了 DOM 操作。其实很简单：

  ```js
  invoker = { value: someFn }

  elem.addEventListener(type, invoker.value)

  // 如果事件更新，只修改 invoker.value 即可，不用进行 DOM 操作
  ```

- **异步更新**
  响应式原本是同步的，即 data 属性变化之后，effectFn 会同步触发执行。
  但如果多次修改 data 属性，会同步触发多次 effectFn 执行，如果用于渲染 DOM 就太浪费性能了。
  所以，Vue 在此基础上进行了优化，改为异步渲染，`即多次修改 data 属性，只会在最后一次触发 effectFn 执行，中间不会连续触发`。

## 框架设计概览

1. 权衡的艺术
   **框架的设计，本身就是一种权衡的艺术**

   - **命令式 vs 声明式**
     `命令式关注过程，声明式关注结果`
     `框架设计要做的就是，保持可维护性，同时让性能损失更少`
     前端开发中，无论是 jQuery Vue 还是 React ，其实都是两者的结合：用声明式去写 UI 配置，用命令式去做业务逻辑处理。
   - **虚拟 DOM**
     性能：innerHTML < 虚拟 DOM < 原生 JavaScript
     心智负担/可维护性：虚拟 DOM < innerHTML < 原生 JavaScript
   - **运行时和编译时**
     三种选择：纯运行时、运行时+编译时、纯编译时
     手写 vdom 太麻烦，所以通过模板 or JSX 完成书写，但是这些还需要经历一次编译
     Vue:运行时+编译时
     Svelte: 纯编译时

2. 框架设计的核心要素

   - **提升开发体验**
     用户未按要求使用时的错误提示，console 自定义 formatter 直观显示数据
   - **控制框架代码的体积**
     开发环境提供良好提示，不增加生产环境体积。（通过 `__DEV__` 变量判断）
   - **Tree-Shaking**
     如果一个函数有副作用，将不会被 Tree-Shaking 删除，通过 `/#__PURE__*/` 在打包工作（ rollup/webpack ）中声明没有副作用。
   - **输出产物**
     IIFE：rollup format: 'iife'
     ESM: format: 'esm'
     CommonJS: format: 'cjs'
   - **特性开关(类似 rust 的 feature flag)**
     用户关闭的特性，利用 Tree-Shaking 不打包到最终资源里。比如在 Vue3 中，对于 options API
   - **处理错误**
     执行用户提供的函数时，做统一的错误处理（try...catch），并提供给用户错误接口来处理。（可以错误上报等。）
   - **TypeScript 类型支持**

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

5. 对象的响应性实现原理 : `Proxy`
   Proxy 只能代理对象，不能代理非对象原始值，比如字符串。
   Proxy 会拦截对对象的基本语义，并重新定义对象的基本操作

6. 非对象的响应性实现原理：

- `引入 ref 的概念`
  既然原始值无法使用 Proxy 我们就只能把原始值包裹起来
  为了判断一个对象是否是原始值的包裹对象，`添加一个不可写不可枚举属性来判断`

  ```js
  function ref(val) {
    const wrapper = {
      value: val
    }
    Object.defineProperty(wrapper, '__v_isref', {
      value: true
    })
    return reactive(wrapper)
  }
  ```

- 响应丢失问题：通过扩展运算符获取响应式对象的值后，我们得到的值变成了普通对象
  可以在新建对象，然后把对应属性的 get 访问器设置为读取之前对象的值，这样就可以出发响应了
  如果属性多的时候，需要进行批量转换

```js
function toRef(obj, key) {
  const wrapper = {
    get value() {
      return obj[key]
    },
    set(val) {
      obj[key] = val
    }
  }
  Object.defineProperty(wrapper, '__v_isref', {
    value: true
  })
  return wrapper
}

function toRefs(obj) {
  const ret = {}
  for (let key in obj) {
    ret[key] = toRef(obj, key)
  }
  return ret
}
const newObj = { ...toRefs(obj) }

const newObj = {
  foo: toRef(obj, 'foo'),
  bar: toRef(obj, 'bar')
}
```

- 自动脱 ref
  我们在模板中使用 ref 对象值的时候，不需要用户再添加 .value 去使用，所以需要有一个自动脱 ref 的功能。
  思路就是`通过一个 Proxy 代理对象，在读取的值为 ref 时，再读取对象的 .value 值，同时设置值也应该有自动设置到 value 属性的功能`

```js
function proxyRefs(target) {
  return new Proxy(target, {
    get(target, key, receiver) {
      const value = Reflect.get(target, key, receiver)
      return value.__v_isRef ? value.value : value
    },
    set(target, key, newValue, receiver) {
      const value = target[key]
      if (value.__v_isRef) {
        value.value = newValue
        return true
      }
      Reflect.set(target, key, newValue, receiver)
    }
  })
}
```

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

- 渲染器需要有跨平台的能力
- 在浏览器端会渲染为真实的 DOM 元素
- 我们实现 createRenderer 函数，它会创建不同平台的渲染器。其中 render 用于浏览器。渲染分两种情况，挂载和后续渲染（存在旧节点），我们在内部使用 patch 去实现具体渲染（暂未实现）。

```js
function createRenderer() {
  // n1 旧node
  // n2 新node
  // container 容器
  // patch可以用户挂载 也可以用于后续渲染
  function patch(n1, n2, container) {}

  function render(vnode, container) {
    if (vnode) {
      // 如果有新 vnode 就和旧 vnode 一起用 patch 处理
      patch(container._vnode, node, container)
    } else {
      // 没有新 vnode 但是有旧 vnode 直接清空 DOM 即可
      if (container._vnode) {
        container.innerHTML = ''
      }
    }
    // 把旧 vnode 缓存到 container
    container._vnode = vnode
  }

  function hydrate(vnode, container) {
    // 服务端渲染
  }

  return {
    render,
    hydrate
  }
}
```

- 自定义渲染器
  把平台相关的函数提取出来，并通过参数传入，然后封装多平台通用的渲染器
  挂载要依赖平台的实现，比如创建元素，插入元素，我们把这些通过 options 统一传入

8. DOM 的挂载和更新的逻辑

- 事件的处理
  **存储一个事件处理函数，并把真正的事件函数赋值到该函数**
- 事件冒泡与更新时机问题
  记录事件触发的时间和事件绑定的时间，只有触发时间在绑定时间之后才会执行

9.  Diff 算法
    Vue3 使用了快速 diff 算法，参考了 ivi 和 inferno 。思路是：

    先进行双端比较
    剩余的部分计算出最长递增子序列（一个很常见的算法），以找到不用重建和移动的节点
    最后处理剩余部分

## 组件化

10. 组件的实现原理

- 渲染组件
  在渲染器内部的实现看，一个组件是一个特殊类型的虚拟 DOM 节点。
  之前在 patch 我们判断了 VNode 的 type 值来处理，现在来处理类型为对象的情况
- 组件状态与自更新
  在渲染时把组件的状态设置为响应式，并把渲染函数放在 effect 中执行，这样就实现了组件状态改变时重新渲染。`同时指定 scheduler 来让渲染队列在一个微任务中执行并进行去重`
- 组件实例与组件的生命周期
  当状态修改导致组件再次渲染时，patch 不应该还是挂载，所以我们需要维护一个实例，记录组件的状态，是否被挂载和上一次的虚拟 DOM 节点
- props 与组件的被动更新

- setup 函数
  返回一个函数，该函数将作为组件的 render 函数
  返回一个对象，其中包含的数据暴露给模板使用

11. 异步组件与函数式组件

12. 内建组件和模块
    KeepAlive、Teleport、Transition

- KeepAlive 一词借鉴了 HTTP 协议。
  KeepAlive 组件可以避免组件被频繁的销毁/重建。本质是缓存管理，再加上特殊的挂载卸载逻辑。
  卸载时将组件放入另一个容器中，再次挂载时再拿出来。对应生命周期为 activated 和 deactivated。
  我们创建内置组件 KeepAlive 后，需要在卸载和挂载时对它进行特殊处理。
  keepalive 的本质是缓存管理，在加上特殊的挂载/卸载逻辑。
  挂载一个被 keepalive 的组件时，它并不会真的被卸载，而会被移动到一个隐藏容器中。
  当重新”挂载“该组件时，它也不会被真的挂载，而是被从隐藏容器中取出，再搬运到原来的容器中。
  `这个过程对应到组件的 activated 和 deactivated 生命周期中。`

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
