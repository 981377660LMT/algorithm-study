## 文本编辑器中的数据结构

1. 对顶栈+邻接表+command 实现 undo/redo 逻辑(移动光标)
2. 链表 => 插入删除为 O(1)
3. C++的 rope/splay(二叉树,旋转操作保证中序遍历不变) 便于插入和删除 高效地处理字符串
4. piece table(改进版本的 gap buffer)
5. 每行一个 immutable array
6. codemirror: b tree

## 洛谷上的进阶题

[P4567 [AHOI2006]文本编辑器(rope/splay 需要支持字符串反转操作)](https://www.luogu.com.cn/problem/P4567)
[P2201 数列编辑器(对顶栈)](https://www.luogu.com.cn/problem/P2201)
[P4008 [NOI2003] 文本编辑器(对顶栈)](https://www.luogu.com.cn/problem/P4008)
[P5599 【XR-4】文本编辑器](https://www.luogu.com.cn/problem/P5599)
