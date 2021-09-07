class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(1)
const c = new Node(1)
const d = new Node(2)
const e = new Node(2)
const f = new Node(3)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 只保留原始链表中 没有重复出现 的数字。
// 返回同样按升序排列的结果链表。

const deleteDuplicates = function (head: Node) {
  if (!head) return head
  const dummy = new Node(0, head) // 头节点可能被删

  // pre 写(改变next) cur 读
  let pre: Node | undefined = dummy
  let cur: Node | undefined = head

  while (cur) {
    while (cur.next && cur.next.value === cur.value) {
      cur = cur?.next
    }

    // 不用删
    if (pre!.next === cur) {
      pre = pre?.next
      cur = cur.next
    } else {
      // 要删
      // 这一句删除了重复元素
      pre.next = cur.next
      cur = cur.next
    }
  }

  return dummy.next
}

console.dir(deleteDuplicates(a), { depth: null })
export {}
