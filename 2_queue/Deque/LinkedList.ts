// 链表是实现的结构不是抽象的结构，
// 由于缓存极不友好，实际表现比从算法复杂度上得出的感觉差很多，
// 实际应用里面但凡能用连续内存（通过将旧内存复制到新内存来扩容）
// 做的都不会用链表这种实现（除非复制的成本非常非常非常高，而这种情况很不常见）。

export { LinkedList } from '../../3_linkedList/LinkedList'

// java查找链表元素：起点折半查找 这样最坏情况也只要找一半就可以了。
// Node<E> node(int index) {
//   // assert isElementIndex(index);

//   if (index < (size >> 1)) {
//       Node<E> x = first;
//       for (int i = 0; i < index; i++)
//           x = x.next;
//       return x;
//   } else {
//       Node<E> x = last;
//       for (int i = size - 1; i > index; i--)
//           x = x.prev;
//       return x;
//   }
// }
