// 一个序列，要在线支持任意位置插入一个数和查询一个数的rank，强制在线，序列中的数互不相同
// OrderMatainer
//
// The List Indexing problem I is that of performing the following operations on a linked list:
// Insert(x, y) Insert a new record y immediately after record x.
// Delete(x) Delete record x from the list.
// Index(i) Return the ith element in the list.
// Position(x) Return the position in the list of record x. That is, Position = Index -1.

interface IListIndexing<V> {
  insertBefore(pivotValue: V, newValue: V): void
  insertAfter(pivotValue: V, newValue: V): void
  delete(value: V): void
}

class ListIndexing {}

export { ListIndexing }
