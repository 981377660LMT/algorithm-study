class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(3)
const c = new Node(5)
const d = new Node(4)
const e = new Node(2)
const f = new Node(6)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 思路:
// 1.返回的不是头节点，需要创建虚拟节点
// 2.遍历head,每次从虚拟节点的头开始循环，看dummyP.next和现在的head.val谁小，
// 3. 一直while找插入的位置，大的数插在后面
const insertionSortList = (head: Node) => {
  const dummy = new Node(0)
  let headP: Node | undefined = head
  let dummyP: Node | undefined = dummy

  // 遍历head
  while (headP) {
    const next: Node | undefined = headP.next
    // 每次循环重置指针从头开始向后找,寻找插入的位置
    dummyP = dummy
    while (dummyP.next && dummyP.next.value < headP.value) {
      dummyP = dummyP?.next
    }
    // 插入headP对应的节点
    headP.next = dummyP.next
    dummyP.next = headP
    headP = next
  }

  return dummy.next
}

console.dir(insertionSortList(a), { depth: null })

export {}
