// 带环:两个人同时起跑，速度快的会超过速度慢的;
// 不带环:遍历结束后，没有相逢

class Node {
  value: number | undefined
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(1)
const c = new Node(2)
const d = new Node(3)
a.next = b
b.next = c
c.next = d
d.next = b

// 检测有环
const hasCycle = (node: Node): boolean => {
  let fastNode: Node | undefined = node
  let slowNode: Node | undefined = node

  while (fastNode) {
    fastNode = fastNode.next?.next
    slowNode = slowNode?.next
    if (slowNode === fastNode) {
      return true
    }
  }

  return false
}

console.log(hasCycle(a))
console.dir(a, { depth: null })

// 检测有环
// 返回链表开始入环的第一个节点。 如果链表无环，则返回 null。
// 使用visited记录
// const detectCycle = (node: Node): Node | null => {
//   let fastNode: Node | undefined = node
//   const visited = new Set<Node>()

//   while (fastNode) {
//     if (visited.has(fastNode)) return fastNode
//     fastNode = fastNode.next
//   }

//   return null
// }

// 快慢指针在相遇之处找入口节点
// 相遇时快节点路程为慢节点两倍，比慢节点多绕了n个内圈y
// A-B=ny 且A=2B 得 B=ny 则慢节点再走x即可到达环的起点

// 相遇时将慢节点置为原点  继续走相遇处就是入口
const detectCycle = (node: Node): Node | undefined => {
  let fastNode: Node | undefined = node
  let slowNode: Node | undefined = node

  while (fastNode && fastNode.next) {
    fastNode = fastNode.next.next
    slowNode = slowNode!.next
    if (slowNode === fastNode) {
      slowNode = node
      while (slowNode !== fastNode) {
        slowNode = slowNode?.next
        fastNode = fastNode?.next
      }
      return slowNode
    }
  }

  return undefined
}

console.log(detectCycle(a))

// O(n)
export {}

// 快慢指针还能判断回文/找到中点/找到倒数第k个结点

const aa = new Node(1)
const bb = new Node(2)
const cc = new Node(3)
const dd = new Node(4)
const ee = new Node(5)
aa.next = bb
bb.next = cc
cc.next = dd
dd.next = ee

/**
 *
 * @param head  长为len的链表
 * @returns 返回(index>>1)下标处的节点
 * 例如5个节点时返回(5-1)>>1===2处的节点
 * 例如6个节点时返回(6-1)>>1===2处的节点
 */
const findMid = (head: Node) => {
  let fastNode: Node | undefined = head
  let slowNode: Node | undefined = head
  while (fastNode && fastNode.next && fastNode.next.next) {
    fastNode = fastNode.next.next
    slowNode = slowNode?.next
  }

  return slowNode
}

console.log(findMid(aa))
console.log(findMid(bb))
