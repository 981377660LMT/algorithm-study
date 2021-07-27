// 仍然是BST
// 红黑树与2-3树等价
// 所有红色节点向左倾斜
// 根节点和空节点是黑色的

// 红色节点(三节点左边的)的孩子是黑色的

// 从任意一个节点到叶子节点经过的黑色节点是一样的
// 红黑树不是平衡二叉树，而是黑平衡的，最大高度2logn，比avl查找慢一点，但是增删比avl快

// avl适合读，红黑树适合写

import { BST, TreeNode } from './BST'

/**
 * 注意三个子过程:左旋转/右旋转/颜色翻转
 */
class RBT extends BST {
  constructor() {
    super()
  }

  /**
   *
   * @param insertedNode 红黑树添加新元素,根节点需要黑色
   * @returns
   */
  override insert(insertedNode: number) {
    if (this.root == null) {
      this.root = new TreeNode(insertedNode)
      this._size++
    } else {
      this.root = this._insert(this.root, insertedNode)
    }

    this.root!.color = 'BLACK'
    return this
  }

  /**
   *
   * @param node 以node为根的二分搜索树
   * @param val 插入的元素
   * @description
   * 1.向二节点插入 插入元素在右侧时需要左旋 左侧直接插入即可
   * 2.向三节点插入 形成临时四节点 需要拆分成三个二节点(黑色) 变色即可 根节点变红色继续融合
   */
  protected override _insert(node: TreeNode | null, val: number): TreeNode | null {
    if (node == null) return node

    // 递归终止条件
    if (node.value > val) {
      if (!node.left) {
        node.left = new TreeNode(val)
        this._size++
      } else {
        this._insert(node.left, val)
      }
    } else if (node.value < val) {
      if (!node.right) {
        node.right = new TreeNode(val)
        this._size++
      } else {
        this._insert(node.right, val)
      }
    }

    // 维护红黑树性质3步

    // 1.需不需要左旋
    if (this.isRed(node.right) && !this.isRed(node.left)) {
      node = this.leftRotate(node)
    }

    // 2.需不需要右旋(黑节点左侧连续两个红节点)
    if (this.isRed(node.left) && this.isRed(node.left?.left)) {
      node = this.leftRotate(node)
    }

    // 3.需不需要颜色翻转
    if (this.isRed(node.left) && this.isRed(node.right)) {
      this.flipColors(node)
    }

    return node
  }

  /**
   *
   * @param node 将node左旋转
   */
  private leftRotate(node: TreeNode): TreeNode {
    const x = node.right!

    node.right = x.left
    x.left = node

    x.color = node.color
    node.color = 'RED'

    return x
  }

  /**
   *
   * @param node 将node右旋转
   */
  private rightRotate(node: TreeNode): TreeNode {
    const x = node.left!

    node.left = x.right
    x.right = node

    x.color = node.color
    node.color = 'RED'

    return x
  }

  /**
   *
   * @param node 颜色反转 中间黑两边红
   */
  private flipColors(node: TreeNode): void {
    node.color = 'RED'
    if (node.left) node.left.color = 'BLACK'
    if (node.right) node.right.color = 'BLACK'
  }

  private isRed(node: TreeNode | null | undefined): boolean {
    // 空节点是黑色
    if (node == null) return false
    return node.color === 'RED'
  }
}

const rbt = new RBT()
rbt.insert(1).insert(5).insert(3).insert(6)
console.log(rbt)
