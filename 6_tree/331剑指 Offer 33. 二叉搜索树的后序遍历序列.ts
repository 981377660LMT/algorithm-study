// 判断该数组是不是某二叉搜索树的后序遍历结果
// https://leetcode-cn.com/problems/er-cha-sou-suo-shu-de-hou-xu-bian-li-xu-lie-lcof/solution/di-gui-he-zhan-liang-chong-fang-shi-jie-jue-zui-ha/
function verifyPostorder(postorder: number[]): boolean {
  let parent = Infinity
  const stack: number[] = []

  for (let i = postorder.length - 1; ~i; i--) {
    const cur = postorder[i]
    // 如果比栈顶小 表示cur为左节点 需要寻找父节点
    while (stack.length && stack[stack.length - 1] > cur) {
      parent = stack.pop()!
    }
    if (cur > parent) return false
    stack.push(cur) // 大于栈顶表示cur为右节点
  }

  return true
}

console.log(verifyPostorder([1, 6, 3, 2, 5]))
console.log(verifyPostorder([1, 3, 2, 6, 5]))

export default 1
