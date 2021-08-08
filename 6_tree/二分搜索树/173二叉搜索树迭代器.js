/**
 * @param {TreeNode} root
 */
var BSTIterator = function (root) {
  this.stack = []
  // 维护一个栈，栈顶的节点都是已经遍历完左子树的节点，
  // 那么我们此时就只用将栈顶结点弹出并返回即可
  while (root) {
    this.stack.push(root)
    root = root.left
  }
}

/**
 * @return {number}
 */
BSTIterator.prototype.next = function () {
  if (this.stack.length) {
    const root = this.stack.pop()
    // 输出完当前弹出的栈顶节点后，不是继续弹出栈顶节点，因为弹出节点的右子树仍未遍历
    let r = root.right
    while (r) {
      this.stack.push(r)
      r = r.left
    }
    return root.val
  } else {
    return null
  }
}

/**
 * @return {boolean}
 */
BSTIterator.prototype.hasNext = function () {
  return this.stack.length > 0
}
