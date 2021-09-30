class Node {
  value: number
  next?: Node
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(2)
const e = new Node(1)
a.next = b
b.next = c
c.next = d
d.next = e

// 想一下回文字符串判断(双指针，一个前移，一个后移，判断相等)
// 解法1：遍历链表存val到数组，再用数组的方法判断
// 两个字符串一个头插一个尾插
// const isPalindrome = (head: Node) => {
//   let str1 = ''
//   let str2 = ''
//   while (head) {
//     str1 = str1 + head.value
//     str2 = head.value + str2
//     head = head.next!
//   }

//   return str1 === str2
// }

// 解法2：快慢指针O(n) 时间复杂度和 O(1) 空间复杂度解决
const isPalindrome2 = (head: Node) => {
  const findMid = (head: Node): Node | undefined => {
    let fastNode: Node | undefined = head
    let slowNode: Node | undefined = head
    while (fastNode && fastNode.next && fastNode.next.next) {
      fastNode = fastNode.next.next
      slowNode = slowNode?.next
    }

    return slowNode
  }

  const reverse = (head: Node | undefined): Node | undefined => {
    if (!head) return head
    let pre: Node | undefined = undefined
    let cur: Node | undefined = head
    while (cur) {
      const next: Node | undefined = cur.next
      cur.next = pre
      pre = cur
      cur = next
    }
    return pre
  }

  const mid = findMid(head)
  let head1: Node | undefined = head
  let head2: Node | undefined = reverse(mid?.next)
  while (head1 && head2) {
    if (head1.value !== head2?.value) return false
    head1 = head1.next
    head2 = head2.next
  }

  return true
}

console.dir(isPalindrome(a), { depth: null })

export {}
