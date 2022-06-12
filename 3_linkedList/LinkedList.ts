import assert from 'assert'
import { LinkedListNode } from './LinkedListNode'

/**
 * @description 链表实现的双端队列
 */
class LinkedList<V = number> {
  /**
   * @description 哨兵
   */
  private readonly _root: LinkedListNode<V>

  /**
   * @description 初始化双向链表，判断节点时 next/pre 若为 root，则表示 next/pre 为空
   */
  constructor(iterable?: Iterable<V>) {
    // @ts-ignore
    this._root = new LinkedListNode<V>(undefined)
    this._root.pre = this._root
    this._root.next = this._root

    for (const item of iterable ?? []) this.push(item)
  }

  unshift(val: V): void {
    this._root.insertAfter(new LinkedListNode(val))
  }

  shift(): V | undefined {
    if (this.isEmpty()) return undefined
    return this._root.next?.remove().value
  }

  push(val: V): void {
    this._root.insertBefore(new LinkedListNode(val))
  }

  pop(): V | undefined {
    if (this.isEmpty()) return undefined
    return this._root.pre?.remove().value
  }

  first(): V | undefined {
    if (this.isEmpty()) return undefined
    return this._root.next?.value
  }

  last(): V | undefined {
    if (this.isEmpty()) return undefined
    return this._root.pre?.value
  }

  isEmpty(): boolean {
    return this._root.next === this._root
  }

  toString(): string {
    return `${[...this]}`
  }

  *[Symbol.iterator](): IterableIterator<V> {
    let node = this._root.next!
    while (node !== this._root) {
      yield node.value
      node = node.next!
    }
  }
}

if (require.main === module) {
  const nums = new LinkedList([1, 2, 3, 4, 5])
  assert.strictEqual(nums.shift(), 1)
  for (const num of nums) console.log(num)
  console.log(nums + '')
}

export { LinkedList }

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
