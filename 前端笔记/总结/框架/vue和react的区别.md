react 和 vue 最大的区别在状态管理方式上，vue 是通过响应式，react 是通过 setState 的 api。
这个是最大的区别，因为它导致了后面 react 架构的变更。
react 的 setState 的方式，导致它并不知道哪些组件变了，需要渲染整个 vdom 才行。但是这样计算量又会比较大，会阻塞渲染，导致动画卡顿。
所以 react 后来改造成了 fiber 架构，目标是`可打断的计算`。
为了这个目标，`不能变对比变更新 dom 了`，所以把`渲染分为了 render 和 commit` 两个阶段，render 阶段通过 schedule 调度来进行 reconcile，也就是找到变化的部分，创建 dom，`打上增删改的 tag`，等全部计算完之后，commit 阶段一次性更新到 dom。
打断之后要找到父节点、兄弟节点，所以 vdom 也被改造成了 fiber 的数据结构，有了 parent、sibling 的信息。

**在 dom 操作前，会异步调用 useEffect 的回调函数，异步是因为不能阻塞渲染。**
**在 dom 操作之后，会同步调用 useLayoutEffect 的回调函数，并且更新 ref。**

1. 监听数据变化原理不同
   **Vue 跟 React 的最大区别在于数据的 reactivity，就是反应式系统上**
   都是**数据驱动的单向数据流**，但是
   Vue:**数据响应式**，数据改动时，界面就会自动更新；组件粒度更新
   React:**手动更新**，需要调用方法 SetState；应用粒度更新
2. 组件通信不同
   Vue 是通过 props、on、children、EventBus、vuex、$root 等方式实现组件通信
   React 可以通过 props 向子组件传递数据或者回调或者 context 实现组件通信
3. 模板渲染不同
   vue:template/jsx
   react:jsx
4. 编程风格
   vue:composition API
   react:hooks
5. dom diff
   vue:Vue 的 compile 阶段的 optimize 标记了 static 点,可以减少 differ 次数,而且是`采用双向遍历方法`;
   react:React 首位是除删除外是固定不动的,然后`向右`依次遍历对比;
6. Vuex 和 Redux
   Vuex 更加灵活一些，组件中既可以 dispatch action 也可以 commit updates，而 Redux 中只能进行 dispatch，并不能直接调用 reducer 进行修改
   Redux 使用的是不可变数据，而 Vuex 的数据是可变的。Redux 每次都是用新的 state 替换旧的 state，而 Vuex 是直接修改

都是处理 UI 层的框架

相同点:

1. 组件化
2. 数据驱动视图
3. 自顶向下
4. react 和 vue 都是基于 vdom 的前端框架，之所以用 vdom 是因为可以`精准的对比`关心的属性，而且还可以`跨平台渲染`。vdom 的渲染就是根据不同的类型来用不同的 dom api 来操作 dom。渲染组件的时候，如果是函数组件，就执行它拿到 vdom。class 组件就创建实例然后调用 render 方法拿到 vdom。vue 的那种 option 对象的话，就调用 render 方法拿到 vdom。
5. Vue 相当于扩展了 html、而 React 相当于扩展了 js。
