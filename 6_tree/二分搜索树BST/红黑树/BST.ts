class TreeNode {
  value: number
  left: TreeNode | null
  right: TreeNode | null
  height: number
  color: 'RED' | 'BLACK'

  constructor(value: number = 0) {
    this.value = value
    this.left = null
    this.right = null
    this.height = 1
    // 默认红颜色，为融合准备
    this.color = 'RED'
  }
}

// 我们的二分搜索树不包含重复元素
class BST {
  protected root: TreeNode | null
  protected _size: number
  constructor() {
    this.root = null
    this._size = 0
  }

  toString() {
    const stringRef = { value: '' }
    this.generateBSTString(this.root, 0, stringRef)
    return stringRef.value
  }

  private generateBSTString(root: TreeNode | null, depth: number, stringRef: { value: string }) {
    if (root == null) {
      stringRef.value += this.generateDepthString(depth) + 'null\n'
      return
    }

    stringRef.value += this.generateDepthString(depth) + root.value + '\n'
    this.generateBSTString(root.left, depth + 1, stringRef)
    this.generateBSTString(root.right, depth + 1, stringRef)
  }

  private generateDepthString(depth: number) {
    return '--'.repeat(depth)
  }

  isEmpty(): boolean {
    return this._size === 0
  }

  size(): number {
    return this._size
  }

  /**
   * 前序遍历按插入顺序输出
   */
  preOrder(): void {
    this._preOrder(this.root)
  }

  private _preOrder(root: TreeNode | null): void {
    if (root == null) return
    console.log(root.value)
    this._preOrder(root.left)
    this._preOrder(root.right)
  }

  /**
   * 中序遍历从小到大输出
   */
  inOrder(): void {
    this._inOrder(this.root)
  }

  private _inOrder(root: TreeNode | null): void {
    if (root == null) return
    this._inOrder(root.left)
    console.log(root.value)
    this._inOrder(root.right)
  }

  /**
   * 后序遍历按插入顺序倒序输出
   */
  postOrder(): void {
    this._postOrder(this.root)
  }

  private _postOrder(root: TreeNode | null): void {
    if (root == null) return
    this._postOrder(root.left)
    this._postOrder(root.right)
    console.log(root.value)
  }

  /**
   * @description BST查找元素是否包含
   */
  contains(val: number): boolean {
    return this._contains(this.root, val)
  }

  /**
   *
   * @param node 以node为根的二分搜索树
   * @param val 查询包含的元素
   */
  private _contains(node: TreeNode | null, val: number): boolean {
    if (node == null) return false
    if (node.value === val) {
      return true
    } else if (node.value > val) {
      return this._contains(node.left, val)
    } else {
      return this._contains(node.right, val)
    }
  }

  /**
   * @description BST添加新元素(递归比非递归简单)
   */
  insert(insertedNode: number) {
    if (this.root == null) {
      this.root = new TreeNode(insertedNode)
      this._size++
    } else {
      this._insert(this.root, insertedNode)
    }

    return this
  }

  /**
   *
   * @param node 以node为根的二分搜索树
   * @param val 插入的元素
   * @description 有待改进
   */
  protected _insert(node: TreeNode | null, val: number): void {
    if (node == null) return

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
  }

  /**
   *
   * @returns 返回删除的二分搜索树最小值
   */
  deleteMin() {
    const min = this.findMin()
    this._deleteMin(this.root)
    return min
  }

  /**
   *
   * @returns 返回删除的二分搜索树最大值
   */
  deleteMax() {
    const max = this.findMax()
    this._deleteMax(this.root)
    return max
  }

  /**
   *
   * @param root 删除以root为根的二叉搜索树的最小值，并将最小值的右子树作为新的左子树代替
   */
  private _deleteMin(root: TreeNode | null) {
    if (!root) return
    if (!root.left) {
      const rightNode = root.right
      root.right = null
      this._size--
      return rightNode
    }

    root.left = this._deleteMin(root.left)!
    return root
  }

  /**
   *
   * @param root 删除以root为根的二叉搜索树的最大值，并将最大值的左子树作为新的右子树代替
   */
  private _deleteMax(root: TreeNode | null) {
    if (!root) return null
    if (!root.right) {
      const leftNode = root.left
      root.left = null
      this._size--
      return leftNode
    }

    root.right = this._deleteMax(root.right)!
    return root
  }

  delete(val: number) {
    this._delete(this.root, val)
    return this
  }

  /**
   * 二分搜索树删除节点:1.只有左孩子的节点2.只有右孩子的节点3.左右孩子都有的节点
   * 将删除节点的前驱/后继节点的值顶上来，递归删除原来的前驱/后继节点
   * @summary 记得return root
   */
  protected _delete(root: TreeNode | null, val: number) {
    if (!root) return null

    if (root.value < val) {
      root.right = this._delete(root.right, val)
      return root
    } else if (root.value > val) {
      root.left = this._delete(root.left, val)
      return root
    } else {
      // 待删除结点左子树为空
      if (!root.left) {
        const rightNode = root.right
        root.right = null
        this._size--
        return rightNode
      }
      // 待删除结点右子树为空
      else if (!root.right) {
        const leftNode = root.left
        root.left = null
        this._size--
        return leftNode
      }
      // 待删除结点左右子树都不为空
      //  找后继节点代替
      else {
        let rootP: TreeNode | null = root.right

        while (rootP.left) {
          rootP = rootP?.left
        }
        root.value = rootP.value
        root.right = this._delete(root.right, root.value)
        return root
      }
    }
  }

  private findMin() {
    if (this.isEmpty()) throw new Error('BST is empty')
    let p: TreeNode | null = this.root
    let value: number | null = null
    while (p) {
      value = p.value
      p = p.left
    }
    return value
  }

  private findMax() {
    if (this.isEmpty()) throw new Error('BST is empty')
    let p: TreeNode | null = this.root
    let value: number | null = null
    while (p) {
      value = p.value
      p = p.right
    }
    return value
  }
}

if (require.main === module) {
  const bst = new BST()
  // const a = new TreeNode(1)
  // const b = new TreeNode(3)
  // const c = new TreeNode(2)
  // const d = new TreeNode(4)

  // bst退化为链表
  // Array.from({ length: 5 }, (_, i) => i + 1).forEach(v => bst.insert(v))
  bst.insert(3).insert(1).insert(6).insert(4).insert(5).insert(5)
  console.dir(bst, { depth: null })
  bst.delete(4)
  console.dir(bst, { depth: null })

  // console.log(bst.contains(3))
  // bst.inOrder()
  // console.log(bst.toString())
}

export { BST, TreeNode }
