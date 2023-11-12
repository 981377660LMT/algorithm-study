/* eslint-disable prefer-destructuring */
// https://noshi91.github.io/Library/data_structure/pairing_heap.cpp
// https://scrapbox.io/data-structures/Pairing_Heap
// https://oi-wiki.org/ds/pairing-heap/

// 配对堆是一棵满足堆性质的带权多叉树（如下图），即每个节点的权值都小于或等于他的所有儿子
// 一个节点的所有儿子节点形成一个单向链表。每个节点储存第一个儿子的指针，即链表的头节点；和他的右兄弟的指针。

type PNode<T> = {
  value: T
  head: PNode<T> | undefined
  next: PNode<T> | undefined
}

/**
 * 配对堆.
 * @link https://oi-wiki.org/ds/pairing-heap/
 */
class PairingHeap<T> {
  /**
   * 合并两个堆，返回新的堆，原来的堆被破坏。
   * 注意两个堆的比较函数必须相同。
   */
  static meld<T>(heap1: PairingHeap<T>, heap2: PairingHeap<T>): PairingHeap<T> {
    return new PairingHeap(
      heap1._compare,
      PairingHeap._merge(heap1._root, heap2._root, heap1._compare)
    )
  }

  private static _merge<T>(
    heap1: PNode<T> | undefined,
    heap2: PNode<T> | undefined,
    compare: (a: T, b: T) => number
  ): PNode<T> | undefined {
    if (!heap1) return heap2
    if (!heap2) return heap1
    if (compare(heap1.value, heap2.value) > 0) {
      const tmp = heap1
      heap1 = heap2
      heap2 = tmp
    }
    heap2.next = heap1.head
    heap1.head = heap2
    return heap1
  }

  private static _mergeList<T>(
    list: PNode<T> | undefined,
    compare: (a: T, b: T) => number
  ): PNode<T> | undefined {
    if (!list || !list.next) return list
    const next = list.next
    const rem = next.next
    return PairingHeap._merge(
      PairingHeap._merge(list, next, compare),
      PairingHeap._mergeList(rem, compare),
      compare
    )
  }

  private readonly _compare: (a: T, b: T) => number
  private _root: PNode<T> | undefined = undefined

  constructor(compare: (a: T, b: T) => number, root?: PNode<T>) {
    this._compare = compare
    this._root = root
  }

  empty(): boolean {
    return !this._root
  }

  top(): T | undefined {
    return this._root ? this._root.value : undefined
  }

  push(value: T): void {
    this._root = PairingHeap._merge(
      this._root,
      {
        value,
        head: undefined,
        next: undefined
      },
      this._compare
    )
  }

  pop(): T | undefined {
    if (!this._root) return undefined
    const res = this._root.value
    this._root = PairingHeap._mergeList(this._root.head, this._compare)
    return res
  }
}

export { PairingHeap }

if (require.main === module) {
  const heap1 = new PairingHeap((a: number, b: number) => a - b)
  const heap2 = new PairingHeap((a: number, b: number) => a - b)
  heap1.push(2)
  heap1.push(3)
  heap1.push(1)
  heap2.push(5)
  heap2.push(5)
  heap2.push(6)
  const heap3 = PairingHeap.meld(heap1, heap2)
  console.log(heap3.pop())
  console.log(heap3.pop())
  console.log(heap3.pop())
  console.log(heap3.pop())
  console.log(heap3.pop())
  console.log(heap3.pop())
  console.log(heap3.pop())
}
