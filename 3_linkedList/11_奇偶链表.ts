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
  if (!head) return head
  let odd = head
  // 最后odd.next要接到headNext
  let headNext = head.next
  // 如何确定边界条件 看odd是否能继续接(不管headNext,因为最后是odd接headNext)
  while (odd?.next?.next) {
    const tmp = odd.next
    odd.next = odd.next.next
    odd = tmp.next!
    tmp.next = odd.next
  }

  odd.next = headNext
  return head
}

console.dir(oddEvenList(a), { depth: null })

export {}
