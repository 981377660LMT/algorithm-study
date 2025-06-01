https://stackoverflow.com/questions/58123398/when-to-use-jsx-element-vs-reactnode-vs-reactelement

1. ReactElement 是一个具有 type、props 和 key 属性的对象
2. 一个 JSX.Element 是 ReactElement<any，any>。它的存在是因为各种库可以以自己的方式实现 JSX
3. ReactNode 是 ReactElement、 字符串 、 数字 、Iterable<ReactNode>、ReactPortal、boolean、null 或 undefined

```jsx
<div> // <- ReactElement
  <Component> // <- ReactElement
    {condition && 'text'} // <- ReactNode
  </Component>
</div>
```
