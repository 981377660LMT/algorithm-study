/* eslint-disable @typescript-eslint/no-non-null-assertion */

import assert from 'assert'

import { LinkedListNode } from './LinkedListNode'

/**
 * 链表实现的双端队列
 */
class LinkedList<E = number> {
  private _size = 0

  /** 哨兵 */
  private readonly _root: LinkedListNode<E>

  /**
   * 初始化双向链表，判断节点时 next/pre 若为 root，则表示 next/pre 为空
   */
  constructor(iterable?: Iterable<E>) {
    this._root = new LinkedListNode<E>(undefined as unknown as E)
    this._root.pre = this._root
    this._root.next = this._root

    for (const item of iterable ?? []) {
      this.push(item)
      this._size++
    }
  }

  unshift(val: E): void {
    this._root.insertAfter(new LinkedListNode(val))
    this._size++
  }

  shift(): E | undefined {
    if (this._isEmpty()) return undefined
    this._size--
    return this._root.next?.remove().value
  }

  push(val: E): void {
    this._root.insertBefore(new LinkedListNode(val))
    this._size++
  }

  pop(): E | undefined {
    if (this._isEmpty()) return undefined
    this._size--
    return this._root.pre?.remove().value
  }

  first(): E | undefined {
    if (this._isEmpty()) return undefined
    return this._root.next?.value
  }

  last(): E | undefined {
    if (this._isEmpty()) return undefined
    return this._root.pre?.value
  }

  toString(): string {
    return `${[...this]}`
  }

  // eslint-disable-next-line generator-star-spacing
  *[Symbol.iterator](): IterableIterator<E> {
    let node = this._root.next!
    while (node !== this._root) {
      yield node.value
      node = node.next!
    }
  }

  get size(): number {
    return this._size
  }

  private _isEmpty(): boolean {
    return this._root.next === this._root
  }
}

if (require.main === module) {
  const nums = new LinkedList([1, 2, 3, 4, 5])
  assert.strictEqual(nums.shift(), 1)
  for (const num of nums) console.log(num)
  console.log(`${nums}`)
}

export { LinkedList }

// java查找链表元素：起点折半查找 这样最坏情况也只要找一半就可以了。
// Node<E> node(int index) {
//   assert isElementIndex(index);

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
