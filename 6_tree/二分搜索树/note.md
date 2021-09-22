平衡二叉树(AVL)指的是：一个二叉树每个节点的左右两个子树的高度差的绝对值不超过 1。

- 如果需要让你判断一个树是否是平衡二叉树，只需要死扣定义，然后用递归即可轻松解决。
- 如果需要你将一个数组或者链表（逻辑上都是线性的数据结构）转化为平衡二叉树，只需要随便选一个节点，并分配一半到左子树，另一半到右子树即可。
- 同时，如果要求你转化为**平衡二叉搜索树**，则可以选择排序数组(或链表)的**中点**，左边的元素为左子树， 右边的元素为右子树即可。

习惯：链表的根写 head,树的根写 root
chang'yong'fang

常用方法:

1. 中序遍历看大小
2. 中序遍历加 pre
   `面试题 04.06. 后继者`
   `2_验证二分搜索树`

```JS
const isValidBST = (root: TreeNode) => {
  if (!root) return true

  let pre: TreeNode | null = null
  const inorder = (root: TreeNode | null): boolean => {
    if (!root) return true
    if (!inorder(root.left)) return false
    if (pre && pre.val >= root.val) return false
    // pre最开始是在最左下角
    pre = root
    if (!inorder(root.right)) return false
    return true
  }

  return inorder(root)
}
```

BST 基本操作
https://labuladong.gitbook.io/algo/mu-lu-ye-1/mu-lu-ye-1/er-cha-sou-suo-shu-cao-zuo-ji-jin
