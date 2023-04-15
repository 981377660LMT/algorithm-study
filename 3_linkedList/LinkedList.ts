/* eslint-disable @typescript-eslint/no-non-null-assertion */

import assert from 'assert'

import { LinkedListNode } from './LinkedListNode'

/**
 * 链表实现的双端队列
 */
class LinkedList<E = number> {
  /** 哨兵 */
  readonly root: LinkedListNode<E>
  private _size = 0

  /**
   * 初始化双向链表，判断节点时 next/pre 若为 root，则表示 next/pre 为空
   */
  constructor(iterable?: Iterable<E>) {
    this.root = new LinkedListNode<E>(undefined)
    this.root.pre = this.root
    this.root.next = this.root

    for (const item of iterable || []) {
      this.push(item)
      this._size++
    }
  }

  unshift(val: E): void {
    this.root.insertAfter(new LinkedListNode(val))
    this._size++
  }

  shift(): E | undefined {
    if (this._isEmpty()) return undefined
    this._size--
    // eslint-disable-next-line prefer-destructuring
    const next = this.root.next
    return next ? next.remove().value : undefined
  }

  push(val: E): void {
    this.root.insertBefore(new LinkedListNode(val))
    this._size++
  }

  pop(): E | undefined {
    if (this._isEmpty()) return undefined
    this._size--
    // eslint-disable-next-line prefer-destructuring
    const pre = this.root.pre
    return pre ? pre.remove().value : undefined
  }

  first(): E | undefined {
    // eslint-disable-next-line prefer-destructuring
    const next = this.root.next
    return next ? next.value : undefined
  }

  last(): E | undefined {
    // eslint-disable-next-line prefer-destructuring
    const pre = this.root.pre
    return pre ? pre.value : undefined
  }

  toString(): string {
    return `${[...this]}`
  }

  insert(node: LinkedListNode<E>, cur: E): LinkedListNode<E> {
    const newNode = new LinkedListNode(cur)
    node.insertBefore(newNode)
    this._size++
    return newNode
  }

  erase(node: LinkedListNode<E>): void {
    node.remove()
    this._size--
  }

  // eslint-disable-next-line generator-star-spacing
  *[Symbol.iterator](): IterableIterator<E> {
    let node = this.root.next!
    while (node !== this.root) {
      yield node.value!
      node = node.next!
    }
  }

  get length(): number {
    return this._size
  }

  private _isEmpty(): boolean {
    return this.root.next === this.root
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
