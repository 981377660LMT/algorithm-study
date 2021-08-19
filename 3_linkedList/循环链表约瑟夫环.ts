class Node {
  value: number
  next: Node | undefined
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

/**
 * @param {number} n
 * @param {number} m
 * @return {number}
 * 每次从这个圆圈里删除第m个数字（删除后从下一个数字开始计数）。求出这个圆圈里剩下的最后一个数字。
 * 1.循环链表模拟
 */
const lastRemaining = (n: number, m: number): number => {
  const nodeArray = Array.from<number, Node>({ length: n }, (_, i) => new Node(i))
  nodeArray.reduce((pre, cur) => (pre.next = cur))
  nodeArray[nodeArray.length - 1].next = nodeArray[0]

  let headP = nodeArray[0]
  let count = 0
  while (headP.next !== headP) {
    // 删除第m-1个节点, 需要找到第m-2个节点
    if (count === m - 2) {
      headP.next = headP.next!.next
      count = 0
    } else {
      count++
    }

    headP = headP.next!
  }

  // console.dir(nodeArray, { depth: null })
  return headP.value
}

console.log(lastRemaining(5, 3))
// console.log(lastRemaining(10, 17))

export default 1
