interface ITreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

class TreeNode implements ITreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val: number = 0) {
    this.val = val
    this.left = null
    this.right = null
  }
}

// ps:中序遍历二分搜索树就得到从小到大的有序数组
// 关键:根节点是中位数
// 选择中点作为根节点，根节点左侧的作为左子树，右侧的作为右子树即可。原因很简单，这样分配可以保证左右子树的节点数目差不超过 1。因此高度差自然也不会超过 1 了。
const sortedArrayToBST = (nums: number[]): TreeNode | null => {
  if (!nums.length) return null
  const mid = Math.floor(nums.length / 2)
  const root = new TreeNode(nums[mid])

  root.left = sortedArrayToBST(nums.slice(0, mid))
  root.right = sortedArrayToBST(nums.slice(mid + 1))

  return root
}

console.dir(sortedArrayToBST([-10, -3, 0, 5, 9]), { depth: null })
// 输出：[0,-3,9,-10,null,5]
export {}
