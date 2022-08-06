# https://www.zhihu.com/question/52997094

请问编辑器中的 undo 和 redo 操作是如何实现的？ - fleuria 的回答 - 知乎
https://www.zhihu.com/question/52997094/answer/133210061
大致上，文本编辑器的内部数据结构用来管理编辑中的文本序列，需要支持随机插入、删除，也需要实时地检索出当前的文本内容交付 U 渲染，最好也考虑进来 undo/redo。
文章最后的 `piece table 算法`已经属于 immutable 数据结构了，只用一张 piece 表来管理变更，而文件正文可以不读到内存里，省内存，对 undo 操作也很友好。

现在应该用 **Rope** 这个数据结构比较多
Rope 是个 immutable 的数据结构，每次操作后记下树的根，undo 时切一个根就好。

来自知乎用户 LdBeth 的回复
ps:rope 是 vim 用的，piece table 是 VSCode 用的，这两个编辑器的用户吹得多了给你的错觉。实际上没別的会用 rope 和 piece table。

---

请问编辑器中的 undo 和 redo 操作是如何实现的？ - Belleve 的回答 - 知乎
https://www.zhihu.com/question/52997094/answer/560529788
Word 里面是使用 **Piece Table 加 Command queue**

---

请问编辑器中的 undo 和 redo 操作是如何实现的？ - 用心阁的回答 - 知乎
https://www.zhihu.com/question/52997094/answer/133194610

**邻接表+对顶栈+action+部分操作巧妙提示**
**当遇到没有逆操作的操作，或者要记录的状态太大，就可以在操作时提醒无法 undo，执行后清空 undo 的栈。**
有时候一个操作对内容的修改的操作很复杂，要实现一次 undo 就能够完全撤销操作，就需要引入原子操作和复合操作的概念，一个符合操作包含一系列操作（原子或复合)，一个复合操作的逆操作就是逆序组织的每个子操作的逆操作。
这样一次替换全部操作就可以由一系列替换操作组成，一次就可以 undo。

---

**脱离 rope 和 piece table 的朴素的数据结构：每行一个 Array 的数据结构**
所以要视实现而定 - 比如文件的数据结构是一个对应每一行的 immutable array，这样`每次修改一行，就只产生对应行的一个新 copy，其他所有的行都因为 structural sharing 不占用额外空间`。
有点像可持久化线段树，初始化时每一行的内容存储到底部的叶子结点，修改时按照行号二分查找到底部再修改，把新的版本 push 到历史数组里，所有历史版本共享一个引用。

---

https://leetcode.cn/problems/design-a-text-editor/solution/vim-by-981377660lmt-3y3i/
[请问编辑器中的 undo 和 redo 操作是如何实现的？](https://www.zhihu.com/question/52997094/answer/133210061)

最近调查了一下编辑器中 undo 和 redo 的实现逻辑，大概有两种思路

1. 一种是**记录每个版本的状态**，这种思路是使用可持久化数据结构，例如

- **rope** (vim 编辑器内部实现)
- **可持久化线段树** ([每行一个不可变数组的实现](https://www.zhihu.com/question/52997094/answer/133210805)，redux + immutable.js 也有类似实践)
- **piece table** (vscode 编辑器内部实现)

  优点是容易实现，缺点是空间复杂度高，因为需要存储大量的数据(即使是 structural sharing 的可持久化数据结构)。

2. 还有一种思路是**记录每个版本的变化**，用[对顶栈](https://leetcode.cn/problems/design-a-text-editor/solution/by-freeyourmind-kr12/)保存变化的 action ，对顶栈来回倒实现 undo 和 redo
   优点是空间复杂度低，缺点是对**每一个 action，必须要是可撤销的**，这就要求必须每次清晰地写出撤销 action 的逻辑，即 Undo 和 Redo 的实际逻辑应该在每种 Command 中实现

> rope 是一种高效的数据结构，用于存储和操作非常长的可变字符串
> 它减少了应用程序的内存重新分配和数据复制的开销
> 适合的应用场景:将非常长的字符串上分成多个较小的字符串

以下是一个简单的性能测试，在长为 1e6 的字符串中进行 1e4 次插入操作，结果为 rope 耗时 21.849ms，字符串暴力修改 5.334s，可见 rope 的效率非常高，适合用于编辑器**大量文本中插入/删除字符**的场景

![image.png](https://pic.leetcode-cn.com/1659756407-pebGOU-image.png)
