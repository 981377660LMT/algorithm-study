1. 这个内部的 Component 类不仅仅是我想出来的一个概念；它更代表了 React 中 VDOM 节点中状态的存储方式。
   **当 React 决定是时候渲染某个组件时，它会从该节点中提取 Hook 状态。**

   ```js
   // Global reference to current component
   let currentComponent = null

   // Component class to hold hook state array
   class Component {
     constructor() {
       this.state = []
       this.currentHookIndex = 0
     }
     render(renderFn) {
       // Reset state for this render
       currentComponent = this
       // Reset hook index for this render
       this.currentHookIndex = 0
       // Call the component function
       const result = renderFn()
       return result
     }
   }

   function useState(init) {
     const component = currentComponent
     const idx = component.currentHookIndex
     component.state[idx] = component.state[idx] ?? { val: init }
     // Increment for next hook call
     component.currentHookIndex++
     return [component.state[idx].val, data => (component.state[idx].val = data)]
   }

   function Test() {
     const [data, setData] = useState(1)
     console.log(data)
     setData(data + 1)
   }

   // Create component and run renders
   const component = new Component()
   component.render(Test) // 1
   component.render(Test) // 2
   component.render(Test) // 3
   ```

2. 对链表的重构
