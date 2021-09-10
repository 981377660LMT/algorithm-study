class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(4)
const e = new Node(5)
const f = new Node(6)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 这里的奇数节点和偶数节点指的是节点编号的奇偶性
// 链表的第一个节点视为奇数节点，第二个节点视为偶数节点
// 请尝试使用原地算法完成。你的算法的空间复杂度应为 O(1)，时间复杂度应为 O(nodes)，nodes 为节点总数。
const oddEvenList = (head: Node) => {
  if (!head || !head.next) return head
  const dummy1 = new Node(0, head)
  const dummy2 = new Node(0, head.next)
  let odd = dummy1.next
  let even = dummy2.next

  while (odd && odd.next && even && even.next) {
    const oddNext = odd.next.next
    const evenNext = even.next.next
    odd.next = oddNext
    even.next = evenNext

    odd = oddNext
    even = evenNext
  }

  odd!.next = dummy2.next

  return dummy1.next
}

console.dir(oddEvenList(a), { depth: null })

export {}
