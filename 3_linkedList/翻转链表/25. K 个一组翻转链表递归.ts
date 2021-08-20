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

  // 1.判断连续的k个节点是否存在 否则返回head
  // 边界条件:例如k=2时需要head和head.next都存在,循环完后headP指向k个数最后一个
  if (!head) return head
  let headP = head
  for (let i = 1; i < k; i++) {
    headP = headP.next!
    if (!headP) return head
  }

  // 2. reverseTwo k 这一段 反转后headP在头部
  const tmp = headP.next
  headP.next = undefined
  reverse(head)

  // 3.递归
  head.next = reverseKGroup(tmp!, k)

  return headP
}

console.dir(reverseKGroup(a, 3), { depth: null })
// 3,2,1,4,5
export {}
