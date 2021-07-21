// 迭代法反转
// const reverseList = (head: Node, start: number, end: number) => {
//   // 注意p1p2都是node
//   let n1: Node | undefined = head
//   let n2: Node | undefined = head

//   // 起点
//   for (let i = 0; i < start - 1; i++) {
//     n1 = n1?.next
//     n2 = n1?.next
//   }
//   for (let i = 0; i < end - start; i++) {
//     const tmp = n2?.next
//     n2!.next = n1
//     n1 = n2
//     n2 = tmp
//   }

//   console.log(n1, n2)

//   return n1
// }

// console.log(reverseList(a, 2, 4))