https://www.bilibili.com/video/BV1EA41127Q3
go context 其实也是一种设计模式

react 的组件树和 go 的 Context 树非常类似：

- 保存值时，声明 Context Provider 的组件上挂上 value，子组件 Consumer 的时候向上查找
- 而且都有生命周期控制
