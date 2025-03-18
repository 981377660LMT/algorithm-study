Vue.js 3.0 在模块的拆分和设计上做得非常合理。模块之间的耦合度非常低，很多模块可以独立安装使用，而不需要依赖完整的 Vue.js 运行时，例如 @vue/reactivity 模块

# 第一篇 框架设计概览

## 第 1 章 权衡的艺术

声明式的更新性能消耗 = 找出差异的性能消耗 + 直接修改的性能消耗

有没有什么办法能够让我们不用付出太多的努力（写声明式代码），还能够保证应用程序的性能下限，让应用程序的性能不至于太差，甚至想办法`逼近命令式代码的性能呢？这其实就是虚拟 DOM要解决的问题`

使用 innerHTML 更新页面的过程是重新构建 HTML 字符串，再重新设置 DOM 元素的 innerHTML 属性，这其实是在说，哪怕我们`只更改了一个文字，也要重新设置 innerHTML 属性`。而重新设置 innerHTML 属性就等价于`销毁所有旧的 DOM 元素，再全量创建新的 DOM 元素`

- **纯运行时**的框架： render(vNode, container) 的写法
- 能不能引入编译的手段，把 HTML 标签编译成树型结构的数据对象，这样不就可以继续使用 Render 函数了吗

```ts
const html = `
<div>
<span>hello world</span>
</div>
`
// 调用 Compiler 编译得到树型结构的数据对象
const obj = Compiler(html)
// 再调用 Render 进行渲染
Render(obj, document.body)
```

这时我们的框架就变成了一个**运行时+编译时(vue3)**的框架
它既支持运行时，用户可以直接提供数据对象从而无须编译；又支持编译时，用户可以提供 HTML 字符串，我们将其编译为数据对象后再交给运行时处理

- **纯编译时(Svelte)**
  不支持任何运行时内容，用户的代码通过编译器编译后才能运行。
  将 HTML 字符串编译为命令式代码的过程

## 第 2 章 框架设计的核心要素

- 用户体验：在 Vue.js 的源码中，我们经常能够看到 warn 函数的调用
- 控制台输出结果：浏览器允许我们编写自定义的 formatter，从而自定义输出形式。
  vue 中使用 initCustomFormatter
  勾选“Console”→“Enable custom formatters” 选项后，输出内容变得非常直观
- 控制框架代码的体积：每一个 warn 函数的调用都会配合 `__DEV__` 常量的检查
  构建生产环境资源时，`__DEV__` 常量替换为字面量 false，这时我们发现这
  段分支代码永远都不会执行，因为判断条件始终为假，这段永远不会
  执行的代码称为 **dead code**，它不会出现在最终产物中，在构建资源的
  时候就会被移除，因此在 vue.global.prod.js 中是不会存在这段代码的
  这样我们就做到了在开发环境中为用户提供友好的警告信息的同时，不会增加生产环境代码的体积
- 框架要做到良好的 Tree-Shaking
  Tree-Shaking 指的就是消除那些永远不会被执行的代码，也就是排除 dead code
  想要实现 Tree-Shaking，必须满足一个条件，即模块必须是 ESM（ES Module），因为 **Tree-Shaking 依赖 ESM 的静态结构**
  Tree-Shaking 中的第二个关键点——副作用。如果一个函数调用会产生`副作用`，那么就不能将其移除
  JavaScript 本身是动态语言，因此想要静态地分析`哪些代码是 dead code 很有难度`
  因为静态地分析 JavaScript 代码很困难，所以像 rollup.js 这类工具都会提供一个机制，让`我们能明确地告诉 rollup.js：“放心吧，这段代码不会产生副作用，你可以移除它。”`

  ```js
  import { foo } from './utils'

  foo()
  ```

  基于这个案例，我们应该明白，在编写框架的时候需要合理使用
  `/*#__PURE__*/` 注释
  这会不会对编写代码造成很大的心智负担呢？其实不会，`因为通常产生副作用的代码都是模块内函数的顶级调用`

- 框架应该输出怎样的构建产物(iife、cjs、esm，其中 esm 又有两个)

  - HTML 页面中使用 <script src> 标签: IIFE`(立即调用的函数表达式)`

  ```js
  var Vue = (function (exports) {
    // ...
    exports.createApp = createApp
    // ...
    return exports
  })({})
  ```

  - Node.js: CommonJS

    ```js
    const Vue = require('vue') // 服务端渲染
    ```

  - 浏览器 esm: ESM
    有**两个** esm 版本

    - vue.esm-browser.js
      html 标签引入

      ```html
      <script type="module" src="vue.esm-browser.js"></script>
      ```

    - vue.esm-bundler.js
      webpack 等打包工具使用

    无论是 rollup.js 还是 webpack，在寻找资源时，如果 package.json 中存在 module 字段，那么会优先使用 `module 字段指向的资源来代替 main 字段指向的资源`

    ```json
    {
      "main": "index.js",
      "module": "dist/vue.runtime.esm-bundler.js"
    }
    ```

    **区别是什么：**
    brower 版本通过`__DEV__`常量来判断是否是开发环境，而 bundler 版本则通过`process.env.NODE_ENV`来判断是否是开发环境

- 特性开关
  本质上是利用 rollup.js 的预定义常量插件来实现。
  ```js
  // support for 2.x options
  if (__FEATURE_OPTIONS_API__) {
    currentInstance = instance
    pauseTracking()
    applyOptions(instance, Component)
    resetTracking()
    currentInstance = null
  }
  ```
- 错误处理
  代替用户统一处理错误，为用户提供统一的错误处理接口 registerErrorHandler
- 良好的 TypeScript 类型支持
  对 TS 类型的支持是否完善也成为评价一个框架的重要指标
  使用 TS 编写代码与对 TS 类型支持友好是两件事
  Vue.js 源码中的 **runtimecore/src/apiDefineComponent.ts 文件**，`整个文件里真正会在浏览器中运行的代码其实只有 3 行，但是全部的代码接近 200 行`，其实这些代码都是在为类型支持服务。由此可见，框架想要做到完善的类型支持，需要付出相当大的努力。

## 第 3 章 Vue.js 3 的设计思路

- 声明式地描述 UI
  前端页面都涉及哪些内容

  - DOM 元素：例如是 div 标签还是 a 标签。
  - 属性：如 a 标签的 href 属性，再如 id、class 等通用属性。
  - 事件：如 click、keydown 等。
  - 元素的层级结构：DOM 树的层级结构，既有子节点，又有父节点

  Vue.js 3 除了支持使用模板描述 UI 外，还支持使用虚拟 DOM 描述 UI ，使用 JavaScript 对象描述 UI 更加灵活 -> h 函数(一个辅助创建虚拟 DOM 的工具函数)
  Vue.js 会根据`组件的 render 函数的返回值拿到虚拟 DOM`，然后就可以把组件的内容渲染出来了

- 渲染器 renderer

```js
function renderer(vnode, container) {
  if (typeof vnode.tag === 'string') {
    // 说明 vnode 描述的是标签元素
    mountElement(vnode, container)
  } else if (typeof vnode.tag === 'function') {
    // 说明 vnode 描述的是组件
    mountComponent(vnode, container)
  }
}
```

- 组件就是一组 DOM 元素的封装
  组件的返回值(或者 render 函数的返回值)是虚拟 DOM

## 第 4 章 响应系统的作用与实现

## 第 5 章 非原始值的响应式方案

代理，指的是对一个对象基本语义的代理。它允许我们拦截并重新定义对一个对象的基本操作。在实现代理的过程中，我们遇到了访问器属性的 this 指向问题，这需要`使用 Reflect.* 方法并指定正确的 receiver 来解决`。

对一个普通对象的所有可能的读取操作

- 访问属性：obj.foo。
- 判断对象或原型上是否存在给定的 key：key in obj。
- 使用 for...in 循环遍历对象：for (const key in obj)
  {}。
- 删除: delete obj.foo

`避免循环调用导致的调用栈溢出`
使用一个标记变量 shouldTrack 来代表是否允许进行追踪，然后重写了上述这些方法，目的是，当这些方法间接读取 length 属性值时，我们会先将
shouldTrack 的值设置为 false，即禁止追踪。这样就可以断开 length 属性与副作用函数之间的响应联系，从而避免循环调用导致的调用栈溢出。

## 第 6 章 原始值的响应式方案

- JavaScript 的 Proxy 无法提供对原始值的代理，所以我们需要使用一层对象作为包裹，间接实现原始值的响应式方案
- toRef 以及 toRefs 这两个函数解决响应丢失问题
- 暴露到模板中的响应式数据自动脱 ref

# 第 7 章 渲染器的设计

**渲染器的作用是把虚拟 DOM 渲染为特定平台上的真实元素**

# 第 8 章 挂载与更新

渲染器的核心功能：挂载与更新

- 挂载
  首先讨论了如何挂载子节点，以及节点的属性。
  对于子节点，只需要递归地调用 patch 函数完成挂载即可。
  而节点的属性比想象中的复杂，它涉及两个重要的概念：HTML Attributes 和 DOM Properties。为元素设置属性时，我们不能总是使用 setAttribute 函数，也不能总是通过元素的 DOM Properties 来设置。至于如何正确地为元素设置属性，取决于被设置属性的特点。例如，表单元素的 el.form 属性是只读的，因此只能使用 setAttribute 函数来设置。

`el.className` 修改 class 属性的性能最优，

- 卸载
  直接使用 innerHTML 来清空容器元素存在诸多问题

  - 容器的内容可能是由某个或多个组件渲染的，当卸载操作发生时，应该正确地调用这些组件的 beforeUnmount、unmounted 等生命周期函数
  - 不会移除绑定在 DOM 元素上的事件处理函数

  - 事件更新:`不用(remove+add，而是直接 update)`
    为了提升性能，我们伪造了 invoker 函数，并把真正的事件处理函数存储在 invoker.value 属性中，当事件需要更新时，只更新 invoker.value 的值即可，这样可以避免一次 removeEventListener 函数的调用。
  - 屏蔽所有绑定时间晚于事件触发时间的事件处理函数的执行。
  - 对虚拟节点中的 children 属性进行了规范化，规定 vnode.children 属性只能有三种类型: string、array、null

# 第 9 章 简单 Diff 算法

渲染器通过 key 属性找到可复用的节点，然后尽可能地通过 DOM 移动操作来完成更新，避免过多地对 DOM 元素进行销毁和重建

# 第 10 章 双端 Diff 算法

在新旧两组子节点的四个端点之间分别进行比较，并试图找到可复用的节点

# 第 11 章 快速 Diff 算法

借鉴了文本 Diff 中的预处理思路，先处理新旧两组子节点中相同的前置节点和相同的后置节点。当前置节点和后置节点全部处理完毕后，如果无法简单地通过挂载新节点或者卸载已经不存在的节点来完成更新，则需要根据节点的索引关系，构造出一个最长递增子序列。**最长递增子序列所指向的节点即为不需要移动的节点。**

# 第 12 章 组件的实现原理

- 使用虚拟节点的 `vnode.type 属性来存储组件对象`，渲染器根据虚拟节点的该属性的类型来判断它是否是组件。如果是组件，则渲染器会使用
  mountComponent 和 patchComponent 来完成组件的挂载和更新
- 在组件挂载阶段，会为组件创建一个用于渲染其内容的副作用函数 effect，该副作用函数会与组件自身的响应式数据建立响应联系；无论对响应式数据进行多少次修改，副作用函数都只会重新执行一次( scheduler:queueJob)
- 组件实例上是一个对象，包含了组件运行过程中的状态，例如组件是否挂载、组件自身的响应式数据，以及组件所渲染的内容（即 subtree）等。有了组件实例后，在渲染副作用函数内，我们就可以根据组件实例上的状态标识，来决定应该进行全新的挂载，还是应该打补丁
- 渲染上下文（renderContext），它实际上是组件实例的代理对象
- setup 函数的返回值可以是两种类型，如果返回函数，则将该函数作为组件的渲染函数；如果返回数据对象，则将该对象暴露到渲染上下文中
- emit 函数包含在 setupContext 对象中，可以通过 emit 函数发射组件的自定义事件。通过 v-on 指令为组件绑定的事件在经过编译后，会以 onXxx 的形式存储到 props 对象中。`当 emit 函数执行时，会在 props 对象中寻找对应的事件处理函数并执行它`
- 通过 onMounted 注册的生命周期函数会被注册到当`前组件实例的 instance.mounted 数组中`。为了维护当前正在初始化的组件实例，我们定义了全局变量 currentInstance，以及用来设置该变量的 setCurrentInstance 函数。

# 第 13 章 异步组件与函数式组件

- 异步的方式加载并渲染一个组件。这在代码分割、服务端下发组件等场景中尤为重要
- 在框架层面为异步组件提供更好的封装支持，与之对应的能力如下
  允许用户指定加载`出错`时要渲染的组件。
  允许用户指定 `Loading` 组件，以及展示该组件的`延迟`时间(网络状况良好的情况下，异步组件的加载速度会非常快，这会导致 Loading 组件刚完成渲染就立即进入卸载阶段，于是出现闪烁的情况，加 delay 避免 Loading 组件导致的闪烁问题)。
  允许用户设置加载组件的`超时`时长。
  组件加载失败时，为用户提供`重试`的能力。
  因此，框架有必要内建异步组件的实现。
- defineAsyncComponent 函数，用来定义异步组件

```js
const AsyncComponent = defineAsyncComponent({
  loader: () => import('./AsyncComponent.vue'),
  loadingComponent: LoadingComponent,
  errorComponent: ErrorComponent,
  delay: 200, // 延迟 200ms 展示 loading 组件，避免Loading 组件闪烁
  timeout: 3000, // 加载组件的超时时长
  suspensible: false,
  onError(error, retry, fail, attempts) {
    if (error.message.match(/fetch/)) {
      return retry()
    }
    if (attempts <= 3) {
      return retry()
    }
    return fail()
  }
})
```

    - 指定了超时时长，则开启一个定时器计时，超时后将 timeout 设置为 true，被卸载时清除定时器
    - 如果用户指定了延迟时间，则开启延迟定时器。定时器到时后，再将 loading.value 的值设置为 true
    - 出错后重试，将onError函数暴露给用户

- **在 Vue.js 3 中使用函数式组件，主要是因为它的简单性，而不是因为它的性能好**
  函数式组件没有自身状态，但它仍然可以接收由外部传入的 props

# 第 14 章 内建组件和模块

与渲染器的结合非常紧密，因此需要框架提供底层的实现与支持

- KeepAlive

  借鉴于 HTTP 协议：在 HTTP 协议中，KeepAlive 又称 HTTP 持久连接（HTTP persistent connection），其作用是`允许多个请求或响应共用一个 TCP 连接`。在没有 KeepAlive 的情况下，一个 HTTP 连接会在每次请求/响应结束后关闭，`当下一次请求发生时，会建立一个新的 HTTP 连接`。`频繁地销毁、创建 HTTP 连接`会带来额外的性能开销，KeepAlive 就是为了解决这个问题而生的
  与 HTTP 中的 KeepAlive 类似，Vue.js 内建的 KeepAlive 组件可以`避免一个组件被频繁地销毁/重建`
  应用场景: `Tab 页签`等

  - KeepAlive 组件的实现`需要渲染器层面的支持`。这是因为被 KeepAlive 的组件在卸载时，我们不能真的将其卸载，否则就无法维持组件的当前状态了 -> activated 和 deactivated (假卸载、假挂载)，类似 hide/restore

- Teleport
  目标：跨越 DOM 层级渲染
  代码实现：在实现 Teleport 时，我们将 `Teleport 组件的渲染逻辑从渲染器中分离出来`，避免渲染器逻辑代码“膨胀”，且可以 Tree-Shaking 删除 Teleport 相关的代码
  Teleport 本质上是渲染器逻辑的合理抽象，它完全可以作为渲染器的一部分而存在
- Transition

# 第 15 章 编译器核心技术概览

分为三个步骤。
(1) 分析模板，将其解析为模板 AST。
(2) 将模板 AST 转换为用于描述渲染函数的 JavaScript AST。
(3) 根据 JavaScript AST 生成渲染函数代码。

# 第 16 章 解析器

# 第 17 章 编译优化

编译器将模板编译为渲染函数的过程中，尽可能多地提取关键信息，并以此指导生成最优代码的过程。
编译优化的策略与具体实现是由框架的设计思路所决定的，不同的框架具有不同的设计思路，因此编译优化的策略也不尽相同。但优化的方向基本一致，即`尽可能地区分动态内容和静态内容，并针对不同的内容采用不同的优化策略`。

1. 传统 Diff 算法的问题
   渲染器在运行时得不到足够的信息。传统 Diff 算法无法利用编译时提取到的任何关键信息，这导致渲染器在运行时不可能去做相关的优化
   而 Vue.js 3 的编译器会将编译时得到的`关键信息“附着”在它生成的虚拟 DOM 上`，这些信息会通过虚拟 DOM 传递给渲染器
2. Block 与 PatchFlags
   补丁标志理解为一系列数字标记，并根据数字值的不同赋予它不同的含义，示例如下。

   数字 1：代表节点有动态的 textContent（例如上面模板中的 p 标签）。
   数字 2：代表元素有动态的 class 绑定。
   数字 3：代表元素有动态的 style 绑定。
   数字 4：其他……

   动态子节点存储到该虚拟节点的 dynamicChildren 数组内 -> 我们把带有该属性的虚拟节点称为“块”，即 Block

3. 静态提升
   把纯静态的节点提升到渲染函数之外
   当响应式数据变化，并使得渲染函数重新执行时，`并不会重新创建静态的虚拟节点`，从而避免了额外的性能开销
4. 预字符串化
   预字符串化能够将这些`静态节点序列化为字符串`，并生成一个 Static 类型的 VNode
   ```js
   const hoistStatic = createStaticVNode('<p></p><p></p><p></p>...20 个...<p></p>')
   ```
5. 缓存内联事件处理函数

   ```js
   function render(ctx, cache) {
     return h(Comp, {
       // 将内联事件处理函数缓存到 cache 数组中
       onChange: cache[0] || (cache[0] = $event => ctx.a + ctx.b)
     })
   }
   ```

   无论执行多少次渲染函数，props 对象中 onChange 属性的值始终不变，于是就不会触发 Comp 组件更新了

6. v-once
   缓存全部或部分虚拟节点

# 第 18 章 服务端渲染

SSR vs CSR vs 同构渲染

四个维度：SEO、白屏问题、占用服务端资源、用户体验

- 将组件渲染为 HTML 字符串
- 客户端激活
  在页面中的 DOM 元素与虚拟节点对象之间建立联系；
  为页面中的 DOM 元素添加事件绑定。

```js
// html 代表由服务端渲染的字符串
const html = renderComponentVNode(compVNode)

// 假设客户端已经拿到了由服务端渲染的字符串
// 获取挂载点
const container = document.querySelector('#app')
// 设置挂载点的 innerHTML，模拟由服务端渲染的内容
container.innerHTML = html

// 接着调用 hydrate 函数完成激活
renderer.hydrate(compVNode, container)
```
