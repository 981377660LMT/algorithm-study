# Deque 常见实现方式

- 两个 slice 头对头拼在一起(golang)
  常见实现，动态扩容
- RingBuffer
  环形缓冲区，固定容量
  disruptor 并发消息队列
- LinkedList
  python deque 分块链表

---

虽然但是 stl 里是用 deque 来实现 queue 的（

stl 的 deque 应该也是 linkedlist 的方式，但是 block 的大小是动态的

另外根据需要，block list 可以做成 O(1)插入，O(k)随机访问，也可以反过来做成 O(k)插入，O(1)随机访问，理论上也可以都是 O(lg(k))但是这不如通过合理的动态 block 大小让 k 是 lgN 的
