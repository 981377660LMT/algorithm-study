React组件什么时候重新渲染

1. 当Parent re-render时，Child也会re-render，除非以下2种情况之一发生：
   - Child组件被React.memo装饰，并且传给Child组件的props与上次渲染相同
   - 父组件渲染的JSX tree中，在某个位置使用的<Child/>对象(react element，即React.createElement()返回的对象)与上次渲染时相同（在同一个位置使用同一个对象引用），这时Child组件不会re-render。
2. 子树：当Child的`状态更新并且re-render的时候`，Parent不会re-render，只有Child以及它的后裔组件会re-render。并且Child组件re-render的时候，收到的props与上次渲染相同。
3. 如果通过setState设置的`新state与当前state相同(Object.is)`，则不会触发re-render
