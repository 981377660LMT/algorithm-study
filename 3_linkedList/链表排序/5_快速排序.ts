// class Node {
//   value: number
//   next?: Node
//   constructor(value: number = 0, next?: Node) {
//     this.value = value
//     this.next = next
//   }
// }

// const a = new Node(4)
// const b = new Node(2)
// const c = new Node(1)
// const d = new Node(3)
// a.next = b
// b.next = c
// c.next = d

// // 需要反转链表的一端节点
// // 使用三节点

// const sortList = (head: Node | undefined): void => {
//   const partition = (head: Node, end: Node | undefined): Node => {
//     //  这里不随机选取pivot了 直接用头
//     const pivotVal = head.value
//   }

//   const quickSort = (head: Node | undefined, end: Node | undefined): void => {
//     if (head !== end) {
//       const pivot = partition(head, end)
//       quickSort(head, pivot)
//       quickSort(pivot.next, end)
//     }
//   }

//   // 是原地快排
//   quickSort(head, undefined)
// }

// console.dir(sortList(a), { depth: null })
// export {}

// 不想看 有点复杂
// 最坏情况也是 n ^ 2 ，因此面试或者竞赛不建议使用
// 推荐归并排序
