class ListNode {
  constructor(public val: number = 0, public next: ListNode | null = null) {}
}

const a = new ListNode(1)
const b = new ListNode(2)
const c = new ListNode(-3)
const d = new ListNode(3)
const e = new ListNode(1)
a.next = b
b.next = c
c.next = d
d.next = e

/**
 * @param {ListNode} head
 * @return {ListNode}
 */
const removeZeroSumSublists = function (head: ListNode): ListNode | null {
  // 想象一下祖玛游戏递归消除的场景:一次只消除一对，然后递归
  const removeTwo = (head: ListNode | null): ListNode | null => {
    // console.dir(head, { depth: null })
    // 因为头节点可能被删 所以要dummy
    const dummy = new ListNode(0, head)
    if (!head) return null
    let headP: ListNode | null = head
    let record = new Map<number, ListNode>([[0, dummy]])
    let curSum = 0

    while (headP) {
      const curVal = headP.val
      curSum += curVal
      if (record.has(curSum)) {
        const preNode = record.get(curSum)!
        preNode.next = headP.next
        return removeTwo(dummy.next)
      } else {
        record.set(curSum, headP)
        headP = headP.next
      }
    }

    return dummy.next
  }
  return removeTwo(head)
}

console.log(removeZeroSumSublists(a))
export {}
