https://www.bilibili.com/video/BV1EA41127Q3
go context 其实也是一种设计模式

react 的组件树和 go 的 Context 树非常类似：

- 保存值时，声明 Context Provider 的组件上挂上 value，子组件 Consumer 的时候向上查找
- 而且都有生命周期控制

---

https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247483677&idx=1&sn=d1c0e52b1fd31932867ec9b1d00f4ec2

## 1. 数据结构

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

## 2. 几种context

- EmptyContext
- CancelContext
- TimeoutContext
- ValueContext
