// 自平衡二叉树结构
// 需要标注节点的高度
// 平衡因子：左右子树高度之差 必须不能超过1
// 什么时候维持平衡:插入/删除节点时 更新父节点平衡因子

import { BST, TreeNode } from './BST'

// AVL 树也是一个二叉搜索树，所以可以直接继承上篇实现的二叉搜索树 BinarySearchTree
class AVLTree extends BST {
  constructor() {
    super()
  }

  /**
   * @description BST添加新元素(递归比非递归简单)
   */
  override insert(val: number): this {
    if (this.root == null) {
      this.root = new TreeNode(val)
      this._size++
    } else {
      this.root = this._insert(this.root, val)
    }

    return this
  }

  /**
   *
   * @param root 中序遍历需要顺序排列
   * @description 中序遍历树看是否递增
   */
  get isBST() {
    const path: number[] = []
    this.avlInOrder(this.root, path)
    for (let i = 1; i < path.length; i++) {
      if (path[i] < path[i - 1]) return false
    }
    return true
  }

  /**
   * @description 是否为平衡二叉树
   */
  get isBalanced(): boolean {
    return this._isBalanced(this.root)
  }

  private avlInOrder(root: TreeNode | null, path: number[]): void {
    if (!root) return
    root.left && this.avlInOrder(root.left, path)
    path.push(root.value)
    root.right && this.avlInOrder(root.right, path)
  }

  private _isBalanced(root: TreeNode | null): boolean {
    if (!root) return true
    const balanceFactor = this.getBalanceFactor(root)
    if (Math.abs(balanceFactor) > 1) return false
    return this._isBalanced(root.left) && this._isBalanced(root.right)
  }

  /**
   *
   * @param node 以node为根的二分搜索树
   * @param insertedNode 插入的元素
   * @description 自底向上组装树，遇到不平衡的节点就旋转
   */
  protected override _insert(node: TreeNode | null, val: number): TreeNode | null {
    if (node == null) return node

    // 递归终止条件
    if (node.value > val) {
      if (!node.left) {
        node.left = new TreeNode(val)
        this._size++
      } else {
        node.left = this._insert(node.left, val)
      }
    } else if (node.value < val) {
      if (!node.right) {
        node.right = new TreeNode(val)
        this._size++
      } else {
        node.right = this._insert(node.right, val)
      }
    }

    node.height = 1 + Math.max(this.getNodeHeight(node.left), this.getNodeHeight(node.right))

    const balanceFactor = this.getBalanceFactor(node)
    if (Math.abs(balanceFactor) > 1) {
      console.log('danger node:', node)
      console.log(balanceFactor)
      console.log(this.getBalanceFactor(node.left))
      console.log(this.getBalanceFactor(node.right))
    }

    // 插入的元素在不平衡节点的左侧的左侧(LL) 右旋转
    if (balanceFactor > 1 && this.getBalanceFactor(node.left) > 0) {
      return this.rightRotate(node)
    }

    // 插入的元素在不平衡节点的右侧的右侧(RR) 左旋转
    if (balanceFactor < -1 && this.getBalanceFactor(node.right) < 0) {
      return this.leftRotate(node)
    }

    // 插入的元素在不平衡节点的左侧的右侧(LR) 左孩子左旋转 不平衡节点右旋转
    if (balanceFactor > 1 && this.getBalanceFactor(node.left) < 0) {
      node.left = this.leftRotate(node.left!)
      return this.rightRotate(node)
    }

    // 插入的元素在不平衡节点的右侧的左侧(RL) 右子右旋转 不平衡节点左旋转
    if (balanceFactor < -1 && this.getBalanceFactor(node.right) > 0) {
      node.right = this.rightRotate(node.right!)
      return this.leftRotate(node)
    }

    return node
  }

  /**
   *
   * @param y 对节点y进行向右旋转操作，返回旋转后新的根节点x
   */
  private rightRotate(y: TreeNode): TreeNode {
    let x = y.left!
    let t3 = x.right

    // 右旋
    x.right = y
    y.left = t3
    // 更新y和x的高度
    y.height = Math.max(this.getNodeHeight(y.left), this.getNodeHeight(y.right)) + 1
    x.height = Math.max(this.getNodeHeight(x.left), this.getNodeHeight(x.right)) + 1

    return x
  }

  /**
   *
   * @param y 对节点y进行向左旋转操作，返回旋转后新的根节点x
   */
  private leftRotate(y: TreeNode): TreeNode {
    let x = y.right!
    let t3 = x.left

    // 左旋
    x.left = y
    y.right = t3
    // 更新y和x的高度
    y.height = Math.max(this.getNodeHeight(y.left), this.getNodeHeight(y.right)) + 1
    x.height = Math.max(this.getNodeHeight(x.left), this.getNodeHeight(x.right)) + 1

    return x
  }

  /**
   * @description 获取节点高度
   */
  private getNodeHeight(node: TreeNode | null) {
    if (node == null) return 0
    return node.height
  }

  /**
   * 获取平衡因子
   */
  private getBalanceFactor(node: TreeNode | null) {
    if (node == null) return 0
    return this.getNodeHeight(node.left) - this.getNodeHeight(node.right)
  }
}

const avl = new AVLTree()

avl.insert(3).insert(6).insert(4)
console.dir(avl, { depth: null })
// console.log(avl.isBST)
console.log(avl.isBalanced)
