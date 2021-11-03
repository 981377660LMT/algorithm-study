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
   vue:Vue 的 compile 阶段的 optimize 标记了 static 点,可以减少 differ 次数,而且是采用双向遍历方法;
   react:React 首位是除删除外是固定不动的,然后依次遍历对比;
6. Vuex 和 Redux
   Vuex 更加灵活一些，组件中既可以 dispatch action 也可以 commit updates，而 Redux 中只能进行 dispatch，并不能直接调用 reducer 进行修改
   Redux 使用的是不可变数据，而 Vuex 的数据是可变的。Redux 每次都是用新的 state 替换旧的 state，而 Vuex 是直接修改

都是处理 UI 层的框架

相同点:

1. 组件化
2. 数据驱动视图
3. 自顶向下
4. vdom 操纵 DOM
5. Vue 相当于扩展了 html、而 React 相当于扩展了 js。
6. 如果你希望快速构建应用，那么应选择 Vue、如果你希望构建复杂的应用，那么应选择 React。
