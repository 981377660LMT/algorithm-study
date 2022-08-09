# 单人编辑下的 undo 与 redo 逻辑

[请问编辑器中的 undo 和 redo 操作是如何实现的？](https://www.zhihu.com/question/52997094/answer/560529788)

https://stackoverflow.com/questions/54416318/how-to-make-a-undo-redo-function

## 备忘录模式

The `Memento Pattern`, where you capture the whole **current state**. It's easy to implement, but memory-inefficient since you need to store similar copies of the whole state.

这种思路是使用持久化数据结构(immutable data structure)，例如 `rope(vim 内部)`/`可持久化线段树(每行一个不可变数组)`/`piece table(vscode 内部,gap buffer 改进版)`，**记录每个版本的状态**
优点是容易实现，缺点是空间复杂度高，因为需要存储大量的数据(即使是可持久化数据结构)。

## 命令模式

The `Command Pattern`, where you capture **commands/actions** that affect the state (the current action and it's inverse action). Harder to implement since for for each undoable action in your application you must explicitly code it's inverse action, but it's far more memory-efficient since you only store the actions that affect the state.

这种思路是使用 command/action **记录每个版本的变化**，例如对顶栈/链表/
优点是空间复杂度低，缺点是对**每一个 action，必须要是可撤销的**，这就要求必须每次清晰地写出撤销 action 的逻辑，即 Undo 和 Redo 的实际逻辑应该在每种 Command 中实现
注意这里的 action 一般是 `{type:...,payload:...}`，每一种 type 对应一个对顶栈存放 action，相当于邻接表

# 多人编辑下的 undo 与 redo 逻辑

[对可多人协同编辑的在线编辑器，如何设计其 undo/redo 的逻辑？](https://www.zhihu.com/question/367915946/answer/2240528814)
[redux 干的事情实质上就是一种 event sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

多人编辑与单人编辑的区别:多人协作编辑时，任何一个闭麦很久的人都可以随时撤销自己很早之前的上一次改动，`历史记录已经不再是严格的栈了，相当于把压在栈底的任意一个状态抽掉`

核心的设计准则其实只有两条

1. **用户只能撤销自己的改动**，这个对其他人状态的命运有很大的关系。
2. **用户从状态 A 独立撤销 N 次之后再重做 N 次，要能回到 A**。注意这条规则并不只是一条产品需求，更是—条重要的技术需求。(`「回放可互相抵消的 operation 逆操作」`)

那么怎样才能实现只撤销自己产生的操作呢？
OT（Operational Transformation）和 CRDT（Conflict-Free Replicated Data Type）算法

## 围绕 operation 设计 ，是以 operation 为抓手，打通对轻量 operation 数据的变换(OT 派)

PT (Operational transformation) 转变算法 ，Google Doc 的协同编辑就是基于该算法的
所有的数据存储和传输都是 Operation 的形式(oplist)

## 围绕 model 设计，以 model 为抓手，实现出可以任意拥抱变化的 model 数据结构(CRDT 派)

CRDT（Conflict-Free Replicated Data Type）无冲突拷贝数据结构 ，atom 的 teletype 协同编辑就是基于它来的。
案例:
[yjs](https://docs.yjs.dev/)
[automerge](https://github.com/automerge/automerge)

**在 2021 年，这个领域的未来已经明显是属于 CRDT 的了**
OT 需要一个中心服务器（用于保证正确性），而 CRDT 则可以支持点对点直接传输数据。
在简单的数据模型下（如项目看板、纯文本编辑），CRDT 可以很好满足协同相关的需求，并且实现起来也相对简单。
但在富文本编辑的场景下，CRDT 可能有会有一些问题。一是 undo/redo 方面的支持比较弱，二是算法的时间和空间复杂度可能是呈指数上升的，会带来性能上的挑战。
