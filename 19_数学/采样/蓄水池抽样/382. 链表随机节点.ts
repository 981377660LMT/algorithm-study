import { randint } from '../randint'

class ListNode {
  val: number
  next: ListNode | null
  constructor(val?: number, next?: ListNode | null) {
    this.val = val === undefined ? 0 : val
    this.next = next === undefined ? null : next
  }
}

class Solution {
  private head: ListNode | null
  constructor(head: ListNode | null) {
    this.head = head
  }

  /**
   * 随机选择链表的一个节点，并返回相应的节点值。
   * 保证每个节点被选的概率一样。
   * @summary
   * 蓄水池抽样
   * 被选中=被选了*不被替换
   */
  getRandom(): number {
    let headP = this.head
    let res = 0
    let count = 0

    while (headP) {
      count++
      const rand = randint(1, count)
      if (rand === count) res = headP.val
      headP = headP.next
    }

    return res
  }
}

export {}
