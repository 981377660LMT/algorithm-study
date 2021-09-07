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

// 请你删除所有重复的元素，使每个元素 只出现一次 。
// const deleteDuplicates = function (head: Node) {
//   if (!head) return null

//   let slow: Node | undefined = head
//   let fast: Node | undefined = head

//   while (fast) {
//     if (fast.value !== slow?.value) {
//       // 注意这里是先前进再赋值
//       slow = slow?.next
//       slow!.value = fast.value
//     }
//     fast = fast.next
//   }

//   slow!.next = undefined
//   return head
// }

// 请你删除所有重复的元素，使每个元素 只出现一次 。
const deleteDuplicates = function (head: Node) {
  if (!head) return head

  let headP: Node | undefined = head

  while (headP) {
    while (headP.next && headP.next.value === headP.value) {
      headP.next = headP.next.next // 删除后面的重复节点
    }

    headP = headP.next
  }

  return head
}

console.dir(deleteDuplicates(a), { depth: null })
export {}
