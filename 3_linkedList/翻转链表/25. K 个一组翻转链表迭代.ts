class Node {
  value: number
  next: Node | undefined
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(4)
const e = new Node(5)
a.next = b
b.next = c
c.next = d
d.next = e

// 思路和两个一组翻转链表相同
// 每 k 个节点一组进行翻转
// 1.判断存在
// 2.reverseTwo
// 3.递归
const reverseKGroup = (head: Node, k: number): Node => {
  if (k === 1) return head
  const reverse = (head: Node) => {
    let p1: Node | undefined = undefined
    let p2: Node | undefined = head
    while (p2) {
      const tmp: Node | undefined = p2.next
      p2.next = p1
      p1 = p2
      p2 = tmp
    }
    return p1
  }
  const dummy = new Node(0, head)
  let dummyP: Node = dummy

  let p1: Node | undefined // 一段的开头
  let p2: Node | undefined // 一段的结尾

  while (true) {
    p1 = dummyP.next
    p2 = dummyP.next
    let count = 1
    // k个一组
    while (p2 && count < k) {
      p2 = p2.next
      count++
    }

    // [p1,p2]凑满k个节点，反转这一段
    if (p2 && count === k) {
      const tmp: Node | undefined = p2!.next
      p2!.next = undefined
      reverse(p1!)
      p1!.next = tmp

      dummyP.next = p2
      dummyP = p1!
    } else {
      return dummy.next!
    }
  }
}

console.dir(reverseKGroup(a, 3), { depth: null })
// 3,2,1,4,5
export {}
