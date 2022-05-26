import { BinaryTree } from '../分类/Tree'
import { deserializeNode } from '../重构json/297.二叉树的序列化与反序列化'

class BSTApi {
  private root: BinaryTree | null

  constructor(root: BinaryTree | null) {
    this.root = root
  }

  get isValidBST(): boolean {
    if (!this.root) return true
    const isValid = (root: BinaryTree | null, min: number, max: number): boolean => {
      if (!root) return true
      if (root.val <= min) return false
      if (root.val >= max) return false
      return isValid(root.left, min, root.val) && isValid(root.right, root.val, max)
    }
    return isValid(this.root, -Infinity, Infinity)
  }

  // 套模板
  insert(val: number) {
    if (!this.root) {
      return (this.root = new BinaryTree(val))
    }
    const insertHelper = (root: BinaryTree | null, val: number): BinaryTree => {
      if (!root) return new BinaryTree(val)
      if (root.val < val) root.right = insertHelper(root.right, val)
      else if (root.val > val) root.left = insertHelper(root.left, val)
      return root
    }
    return insertHelper(this.root, val)
  }

  // 套模板
  search(val: number) {
    if (!this.root) return false
    const searchHelper = (root: BinaryTree | null, val: number): boolean => {
      if (!root) return false
      if (root.val === val) return true
      else if (root.val < val) return searchHelper(root.right, val)
      else if (root.val > val) return searchHelper(root.left, val)
      return false
    }
    return searchHelper(this.root, val)
  }

  // 套模板
  delete(val: number) {
    if (!this.root) return
    const getMin = (root: BinaryTree) => {
      while (root.left) {
        root = root.left
      }
      return root
    }
    const deleteHelper = (root: BinaryTree | null, val: number): BinaryTree | null => {
      if (!root) return root
      if (root.val === val) {
        if (!root.left) return root.right
        if (!root.right) return root.left
        // 后继节点代替(右子树中最小节点)
        // 前驱节点代替(左子树中最大节点)
        // 找到右子树的最小节点
        const successor = getMin(root.right)
        // 把 root 改成 successor
        root.val = successor.val
        // 转而去删除 successor
        root.right = deleteHelper(root.right, successor.val)
      } else if (root.val > val) {
        root.left = deleteHelper(root.left, val)
      } else if (root.val < val) {
        root.right = deleteHelper(root.right, val)
      }
      return root
    }
    return deleteHelper(this.root, val)
  }
}

const bstApi = new BSTApi(deserializeNode([5, 3, 6, 2, 4, null, null, 1])!)
console.log(bstApi.isValidBST)
console.log(bstApi.search(3))
console.dir(bstApi.delete(3), { depth: null })
