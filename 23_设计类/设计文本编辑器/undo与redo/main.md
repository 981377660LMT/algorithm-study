# 单人编辑下的 undo 与 redo 逻辑

[请问编辑器中的 undo 和 redo 操作是如何实现的？](https://www.zhihu.com/question/52997094/answer/560529788)

https://stackoverflow.com/questions/54416318/how-to-make-a-undo-redo-function

## 备忘录模式

The `Memento Pattern`, where you capture the whole **current state**. It's easy to implement, but memory-inefficient since you need to store similar copies of the whole state.

这种思路是使用持久化数据结构(immutable data structure)，例如 rope/可持久化线段树，**记录每个版本的状态**
优点是容易实现，缺点是空间复杂度高，因为需要存储大量的数据(即使是可持久化数据结构)。

## 命令模式

The `Command Pattern`, where you capture **commands/actions** that affect the state (the current action and it's inverse action). Harder to implement since for for each undoable action in your application you must explicitly code it's inverse action, but it's far more memory-efficient since you only store the actions that affect the state.

这种思路是使用 command/action **记录每个版本的变化**，例如对顶栈/链表/piece table(改进版本的 gap buffer)
优点是空间复杂度低，缺点是对**每一个 action，必须要是可撤销的**，这就要求必须每次清晰地写出撤销 action 的逻辑，即 Undo 和 Redo 的实际逻辑应该在每种 Command 中实现
注意这里的 action 一般是 `{type:...,payload:...}`，每一种 type 对应一个对顶栈存放 action，相当于邻接表

# 多人编辑下的 undo 与 redo 逻辑

[对可多人协同编辑的在线编辑器，如何设计其 undo/redo 的逻辑？](https://www.zhihu.com/question/367915946/answer/2240528814)

多人编辑与单人编辑的区别是 每个 undo redo 不是线性变化的

需要遵循的设计原则
