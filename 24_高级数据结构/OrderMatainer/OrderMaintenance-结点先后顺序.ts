/* eslint-disable max-len */
/* eslint-disable prefer-destructuring */
// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/data_structure/order_maintenance.cc#L1
//
// Order Maintenance
//
// - alloc(): return new node x
// - insertAfter(x, y): insert node y after node x
// - erase(x): erase node x from the list
// - order(x, y): return true if x is before y
//
// Running Time:
//   worst case O(1) for create_node, erase, and order.
//   amortized O(log n) for insert; very small constant.
//
// Reference:
//   P. Dietz and D. Sleator (1987):
//   "Two algorithms for maintaining order in a list".
//   STOC.
//

import { shuffle } from 'lodash'

// !API:
//  Alloc() *Node
//  Build(nums []int) // 用nums中的元素构建链表.
//  InsertAfter(x, y *Node) // 将y插入到x后面.
//  IsBefore(x, y *Node) bool // 判断x是否在y前面.
//  Erase(x *Node) // 删除x.

// !维护元素顺序的链表/带插入全序集维护.
// 用于维护元素的先后顺序, 以及判断元素是否在另一个元素的前面.

// 还有一种 AVL方法：https://stackoverflow.com/questions/32839578/order-maintenance-data-structure-in-c

interface INode {
  id: number
  /**
   * 排序值.
   */
  orderedKey: number
  prev: INode
  next: INode
}

interface IOrderMaintainer {
  insertBefore(pivotId: number, newId: number): boolean
  insertAfter(pivotId: number, newId: number): boolean
  insertFirst(newId: number): void
  insertLast(newId: number): void
  delete(id: number): boolean

  /**
   * 获取排序值.
   */
  getOrderedKey(id: number): number
  resetOrderedKey(): void
}

const GAP = 2 ** 36 //
const MAX_KEY = GAP * (1e5 + 10) // 允许的最大排序值

class OrderMaintainer implements IOrderMaintainer {
  private readonly _head = this._createNode(-1)
  private readonly _idToNode = new Map<number, INode>()
  private _length = 0

  constructor(initIds: number[] = []) {
    let pre = this._head
    initIds.forEach(id => {
      const cur = this._createNode(id)
      this._idToNode.set(id, cur)
      this._insertAfter(pre, cur)
      pre = cur
    })
    this._length = initIds.length
  }

  resetOrderedKey(): void {
    let cur = this._head.next
    let count = 1
    while (cur !== this._head) {
      cur.orderedKey = GAP * count++
      cur = cur.next
    }
  }

  insertBefore(pivotId: number, newId: number): boolean {
    if (!this._idToNode.has(pivotId)) {
      console.warn(`pivotId ${String(pivotId)} not found in OrderMaintainer`)
      return false
    }
    const pivotNode = this._idToNode.get(pivotId)!
    const newNode = this._setDefault(newId)
    this._insertBefore(pivotNode, newNode)
    return true
  }

  insertAfter(pivotId: number, newId: number): boolean {
    if (!this._idToNode.has(pivotId)) {
      console.warn(`pivotId ${String(pivotId)} not found in OrderMaintainer`)
      return false
    }
    const pivotNode = this._idToNode.get(pivotId)!
    const newNode = this._setDefault(newId)
    this._insertAfter(pivotNode, newNode)
    return true
  }

  insertFirst(id: number): void {
    const newNode = this._setDefault(id)
    this._insertAfter(this._head, newNode)
  }

  insertLast(id: number): void {
    const newNode = this._setDefault(id)
    this._insertBefore(this._head, newNode)
  }

  delete(id: number): boolean {
    if (!this._idToNode.has(id)) return false
    const node = this._idToNode.get(id)!
    this._delete(node)
    this._idToNode.delete(id)
    return true
  }

  getOrderedKey(id: number): number {
    const node = this._idToNode.get(id)
    if (node === undefined) {
      console.warn(`id ${String(id)} not found in OrderMaintainer`)
      return 0
    }
    return node.orderedKey
  }

  toString(): string {
    const res: Record<number, number> = {}
    let cur = this._head.next
    while (cur !== this._head) {
      res[cur.id] = cur.orderedKey
      cur = cur.next
    }
    return JSON.stringify(res, null, 2)
  }

  get length(): number {
    return this._length
  }

  private _insertBefore(pivot: INode, newValue: INode): void {
    this._insertAfter(pivot.prev, newValue)
  }

  private _insertAfter(pivot: INode, newValue: INode): void {
    const pivotKey = pivot.orderedKey
    if (this._distOnCycle(pivot, pivot.next) <= 1) {
      this._adjustKey(pivot)
    }
    newValue.orderedKey = pivotKey + Math.floor(this._distOnCycle(pivot, pivot.next) / 2)
    newValue.next = pivot.next
    newValue.prev = pivot
    newValue.next.prev = newValue
    newValue.prev.next = newValue
  }

  private _createNode(id: number): INode {
    const res = { id, orderedKey: 0 } as INode
    res.prev = res
    res.next = res
    return res
  }

  private _setDefault(id: number): INode {
    if (this._idToNode.has(id)) return this._idToNode.get(id)!
    const res = this._createNode(id)
    this._idToNode.set(id, res)
    this._length++
    return res
  }

  private _delete(x: INode): void {
    x.prev.next = x.next.prev
    x.next.prev = x.prev.next
    this._length--
  }

  // TODO: 逻辑检查
  private _adjustKey(start: INode): void {
    let next1 = start.next
    let next2 = next1.next
    let adjustNodeCount = 3
    while (next2 !== start && this._distOnCycle(start, next2) <= 4 * this._distOnCycle(start, next1)) {
      adjustNodeCount++
      next2 = next2.next
      if (next2 === start) break
      adjustNodeCount++
      next2 = next2.next
      next1 = next1.next
    }
    let gap = Math.floor((start === next2 ? MAX_KEY : this._distOnCycle(start, next2)) / adjustNodeCount)
    let baseKey = next2.orderedKey
    while (true) {
      if (next2 === this._head) baseKey += MAX_KEY
      next2 = next2.prev
      if (next2 === start) break
      baseKey -= gap
      next2.orderedKey = baseKey
    }
  }

  private _distOnCycle(x: INode, y: INode): number {
    const diff = Math.abs(x.orderedKey - y.orderedKey)
    return Math.min(diff, MAX_KEY - diff)
  }
}

export {}

if (require.main === module) {
  const order = Array.from({ length: 1e5 }, (_, i) => i)
  shuffle(order)

  // console.timeEnd('aaa')

  console.time('OrderMaintainer')
  const T = new OrderMaintainer(order)
  console.timeEnd('OrderMaintainer')
  // console.log(T.toString(), T.length)

  T.resetOrderedKey()
}
