class ListNode {
  constructor(public val: number = 0, public next: ListNode | null = null) {}
}

// 将链表从中间断开，然后右侧反转即可
const a = new ListNode(1)
const b = new ListNode(2)
const c = new ListNode(2)
const d = new ListNode(1)

/**
 * @param {ListNode} head
 * @return {boolean}
 */
const isPalindrome = function (head: ListNode): boolean {
  if (!head) return true
  let slow: ListNode | null = head
  let fast: ListNode | null = head

  while (fast && fast.next) {
    slow = slow!.next
    fast = fast.next.next
  }

  let pre = slow
  slow = slow!.next
  pre!.next = null

  while (slow) {
    const tmp = slow.next
    slow.next = pre
    pre = slow
    slow = tmp
  }

  let l: ListNode | null = head
  let r: ListNode | null = pre

  while (r) {
    if (l?.val !== r.val) return false
    l = l.next
    r = r.next
  }

  return true
}

console.log(isPalindrome(a))
export default 1
