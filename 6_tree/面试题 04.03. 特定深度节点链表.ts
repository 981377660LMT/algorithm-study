import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297.二叉树的序列化与反序列化'

class Node {
  value: number
  next: Node | undefined
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

function listOfDepth(tree: BinaryTree | undefined): Node[] | undefined {
  if (!tree) return tree
  const queue: BinaryTree[] = [tree]
  const res: Node[] = []

  while (queue.length) {
    const levelLength = queue.length
    const dummy = new Node()
    let dummyP = dummy
    // 遍历当前层的所有节点
    for (let i = 0; i < levelLength; i++) {
      const head = queue.shift()!
      dummyP.next = new Node(head.val)
      dummyP = dummyP.next
      head.left && queue.push(head.left)
      head.right && queue.push(head.right)
    }
    res.push(dummy.next!)
  }

  return res
}

console.dir(listOfDepth(deserializeNode([1, 2, 3, 4, 5, null, 7, 8])!))
