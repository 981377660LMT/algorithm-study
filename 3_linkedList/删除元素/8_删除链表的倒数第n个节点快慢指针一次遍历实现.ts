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

// 先用快慢指针定位到要删除的节点
// 见删除某一个节点的操作
const removeNthFromEnd = (head: Node | undefined, n: number) => {
  const dummy = new Node(0, head)
  let slow: Node | undefined = dummy
  let fast: Node | undefined = dummy

  for (let index = 0; index < n; index++) {
    fast = fast?.next
  }

  // 注意这里是fast?.next而不是fast 要让 slow成为要删除的节点的前一个结点
  while (fast?.next) {
    fast = fast.next
    slow = slow?.next
  }

  // 删除此节点
  // slow!.value = slow?.next?.value!
  slow!.next = slow?.next?.next

  return dummy.next
}

console.dir(removeNthFromEnd(a, 2), { depth: null })
// console.dir(removeNthFromEnd(d, 1), { depth: null })
// console.dir(removeNthFromEnd(e, 1), { depth: null })

export {}
