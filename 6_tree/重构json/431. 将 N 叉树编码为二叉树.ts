import { BinaryTree } from '../力扣加加/Tree'

class Node {
  val: number
  children: Node[]
  constructor(val = 0) {
    this.val = val
    this.children = []
  }
}

// 把N叉树一个节点的第一个孩子都作为二叉树的左节点，
// 然后该节点兄弟挂载在第一个孩子的右节点上
class Codec {
  // Encodes a tree to a binary tree.
  // 第一个兄弟节点先挂在左节点，然后其余的兄弟节点挂在左节点的右节点上
  serialize(root: Node | null): BinaryTree | null {
    if (!root) return null
    const node = new BinaryTree(root.val)
    if (root.children.length) {
      node.left = this.serialize(root.children[0])
    }

    let left = node.left
    if (left) {
      for (const sibling of root.children.slice(1)) {
        left.right = this.serialize(sibling)
        left = left.right!
      }
    }

    return node
  }

  // Decodes your encoded data to tree.
  // 去二叉树左节点上面取回原来的节点
  deserialize(root: BinaryTree | null): Node | null {
    if (!root) return null
    const node = new Node(root.val)
    let left = root.left
    while (left) {
      node.children.push(this.deserialize(left)!)
      left = left.right
    }
    return node
  }
}
