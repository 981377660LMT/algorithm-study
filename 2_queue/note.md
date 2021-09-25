双端队列的两种实现方式:LinkedList 双向链表/循环双端数组
Java 集合框架对 Deque 接口有几类实现：
数组（Resizable Array）、链表（Linked List）。

ArrayDeque：数组实现
LinkedList：链表实现

arraydeque 是个可扩容的双向队列,
底层的数组有头尾指针,在作为队列使用时,
效率比 linkedlist 快,这个结论其实 jdk 的描述直接给出来了
