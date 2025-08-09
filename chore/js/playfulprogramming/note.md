# playfulprogramming

https://playfulprogramming.com/search?q=*&filterTags=react

## 通过 as 属性来更改渲染的 HTML 元素

https://playfulprogramming.com/posts/react-as-prop
https://playfulprogramming.com/posts/vue-as-prop

```tsx
// 这种方法使我们能够创建一个具有类型安全、样式一致且易于复用的可重用组件

type PolymorphicProps<E extends React.ElementType> = React.PropsWithChildren<
  React.ComponentPropsWithoutRef<E> & {
    as?: E
  }
>

type TypographyProps<T extends React.ElementType> = PolymorphicProps<T> & {
  color?: string
}

export function Typography<T extends React.ElementType = 'p'>({
  as,
  color,
  ...props
}: TypographyProps<T>) {
  const Component = as || 'p'
  console.log(color)
  return <Component {...props} />
}

function Social() {
  return (
    <section>
      <Typography as="h1" className="mb-4">
        Connect
      </Typography>
      <Typography as="a" className="mb-4" href="https://www.christianvm.dev/" target="_blank">
        Link to my website
      </Typography>
    </section>
  )
}
```

## A11Y

工程中的无障碍是“创建可被尽可能广泛能力范围内的人使用的产品的过程”。

## react 合集

https://playfulprogramming.com/search?q=*&filterTags=react&sort=newest

---

可变变量最大的问题是它们不是线程安全的。线程安全代码的定义是：“线程安全代码仅以确保所有线程正确运行并满足其设计规范且无意外交互的方式操作共享数据结构。”（来源：线程安全 - 维基百科）
不可变性通过确保数据结构不能被修改，只能被读取来解决这个问题。
