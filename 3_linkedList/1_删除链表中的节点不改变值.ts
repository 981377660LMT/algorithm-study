class Node {
  value: number
  next?: Node
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
a.next = b
b.next = c

// 请你删除链表中所有满足 Node.val == val 的节点，并返回 新的头节点
// 如果连续两个节点都是要删除的节点，这个情况容易被忽略。
const deleteNode = (head: Node | undefined, val: number): Node | undefined => {
  if (!head) return head
  const dummy = new Node(0, head)
  let dummyP: Node | undefined = dummy

  while (dummyP && dummyP.next) {
    let next: Node | undefined = dummyP.next
    if (next.value === val) {
      dummyP.next = next.next
      next = next.next
    }
    // 只有下个节点不是要删除的节点才更新 current
    if (!next || next.value !== val) dummyP = next
  }

  return dummy.next
}

console.dir(deleteNode(a, 2), { depth: null })

export { Node }

// O(1)
