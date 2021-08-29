class Node {
  value: number
  next?: Node
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(4)
const b = new Node(2)
const c = new Node(1)
const d = new Node(3)
a.next = b
b.next = c
c.next = d

// 需要反转链表的一端节点
// 使用三节点
const insertSort = (head: Node | undefined) => {
  if (!head || !head.next) return head
  const dummy = new Node(-1, head)
  let dummyP: Node = dummy
  let cur: Node | undefined = head

  while (cur) {
    let last: Node | undefined = cur.next
    // 需要三节点法反转p1,p2这一段链表
    if (last && last.value < cur.value) {
      // 从 dummy 到 cur 线性遍历找第一个满足条件的位置并插入
      while (dummyP.next && dummyP.next.value < last.value) {
        dummyP = dummyP.next
      }
      const tmp = dummyP.next
      cur.next = last.next
      last.next = tmp
      dummyP.next = last
      dummyP = dummy
    } else {
      cur = last
    }
  }

  return dummy.next
}

console.dir(insertSort(a), { depth: null })
export {}
