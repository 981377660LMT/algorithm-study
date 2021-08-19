// 双端(循环)链表便于遍历

class Node {
  value: number
  next: Node | undefined
  pre: Node | undefined
  constructor(value: number = 0, next?: Node, pre?: Node) {
    this.value = value
    this.next = next
    this.pre = pre
  }
}

const gen = (n: number): Node => {
  const nodeArray = Array.from<number, Node>({ length: n }, (_, i) => new Node(i))
  nodeArray.reduce((pre, cur) => (pre.next = cur))
  nodeArray.reduceRight((pre, cur) => (pre.pre = cur))
  nodeArray[0].pre = nodeArray[nodeArray.length - 1]
  nodeArray[nodeArray.length - 1].next = nodeArray[0]
  return nodeArray[0]
}

console.dir(gen(4), { depth: null })
export default 1
