/**
 * @param {TreeNode} root
 * @return {number}
 * @description
 * // Time O(n)
   // Space O(d) = O(log n) to keep the recursion stack, where d is a tree depth

 */
var countNodes = function (root) {
  if (root == null) return 0
  return countNodes(root.left) + countNodes(root.right) + 1
}

// 对于没有约束的二叉树而言，可以很简单地想到这个递归的解法：
