文本编辑器中的数据结构

1. 对顶栈+邻接表+command 实现 undo/redo 逻辑
2. 链表 => 插入删除为 O(1)
3. C++的 rope/splay(二叉树,旋转操作保证中序遍历不变) 便于插入和删除 高效地处理字符串
4. piece table(改进版本的 gap buffer)
5. 每行一个 immutable array
