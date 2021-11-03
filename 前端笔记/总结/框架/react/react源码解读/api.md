1. 组件类
   React.memo: 第二个参数 返回 true 组件不渲染 ， 返回 false 组件重新渲染。 shouldComponentUpdate: 返回 true 组件渲染 ， 返回 false 组件不渲染。
   StrictMode
   用于检测 react 项目中的潜在的问题，。与 Fragment 一样， StrictMode 不会渲染任何可见的 UI 。它为其后代元素触发额外的检查和警告。
   严格模式检查仅在开发模式下运行；它们不会影响生产构建。

   StrictMode 目前有助于：
   ① **识别不安全的生命周期**。
   ② 关于使用过时字符串 ref API 的警告
   ③ 关于使用废弃的 findDOMNode 方法的警告
   ④ 检测意外的副作用
   ⑤ 检测过时的 context API

2. 工具类
   - createElement 把我们写的 jsx，变成 element 对象; 而 **cloneElement** 的作用是以 element 元素为样板克隆并返回新的 React 元素。返回元素的 props 是将新的 props 与原始元素的 props 浅层合并后的结果。比如说，我们可以在组件中，**劫持 children element，然后通过 cloneElement 克隆 element，混入 props**。经典的案例就是 react-router 中的 Swtich 组件，通过这种方式，来匹配唯一的 Route 并加以渲染。
   - isValidElement
     这个方法可以用来检测是否为 react element 元素,接受待验证对象，返回 true 或者 false。这个 api 可能对于业务组件的开发，作用不大，因为对于组件内部状态，都是已知的，我们根本就不需要去验证，是否是 react element 元素。但是，对于一起公共组件或是开源库，isValidElement 就很有作用了。
     React.Children 提供了用于处理 this.props.children **不透明数据结构**的实用方法。
   - Children.map
   - Children.forEach
   - Children.count
   - Children.toArray
   - Children.only
     ```JS
     这个数据结构，我们不能正常的遍历了，即使遍历也不能遍历，每一个子元素。此时就需要 react.Chidren 来帮忙了。
     function Index(){
         return <div style={{ marginTop:'50px' }} >
             <WarpComponent>
                 { Array(3).fill(0).map(()=><Text/>) }
                 <span>hello,world</span>
             </WarpComponent>
         </div>
     }
     ```
3. ReatDOM

```JS
ReactDOM.createPortal(child, container)
```
