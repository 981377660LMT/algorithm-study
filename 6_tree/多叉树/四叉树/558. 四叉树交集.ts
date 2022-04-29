import { construct } from './427. 建立四叉树'

class Node {
  val: boolean // val：储存叶子结点所代表的区域的值。1 对应 True，0 对应 False
  isLeaf: boolean // 当这个节点是一个叶子结点时为 True，如果它有 4 个子节点则为 False
  topLeft: Node | null
  topRight: Node | null
  bottomLeft: Node | null
  bottomRight: Node | null
  constructor(
    val: boolean,
    isLeaf: boolean,
    topLeft: Node | null,
    topRight: Node | null,
    bottomLeft: Node | null,
    bottomRight: Node | null
  ) {
    this.val = val
    this.isLeaf = isLeaf
    this.topLeft = topLeft
    this.topRight = topRight
    this.bottomLeft = bottomLeft
    this.bottomRight = bottomRight
  }
}

/**
 *
 * @param quadTree1
 * @param quadTree2
 * 返回一个表示 n * n 二进制矩阵的四叉树，它是 quadTree1 和 quadTree2 所表示的两个二进制矩阵进行 按位逻辑或运算 的结果
 */
function intersect(quadTree1: Node | null, quadTree2: Node | null): Node | null {
  return merge(quadTree1, quadTree2)

  function merge(node1: Node | null, node2: Node | null): Node | null {
    // 叶子直接合并
    if (node1?.isLeaf && node2?.isLeaf) {
      return new Node(node1.val || node2.val, true, null, null, null, null)
    }

    // 非叶子节点子节点全是相同叶子节点，合并
    const topLeft = merge(node1?.topLeft || node1, node2?.topLeft || node2)
    const bottomLeft = merge(node1?.bottomLeft || node1, node2?.topLeft || node2)
    const topRight = merge(node1?.topRight || node1, node2?.topRight || node2)
    const bottomRight = merge(node1?.bottomRight || node1, node2?.bottomRight || node2)
    if (
      topLeft?.isLeaf &&
      topRight?.isLeaf &&
      bottomLeft?.isLeaf &&
      bottomRight?.isLeaf &&
      topLeft.val === topRight.val &&
      topLeft.val === bottomLeft.val &&
      topLeft.val === bottomRight.val
    ) {
      return new Node(topLeft.val, true, null, null, null, null)
    }

    return new Node(true, false, topLeft, topRight, bottomLeft, bottomRight)
  }
}

export {}
