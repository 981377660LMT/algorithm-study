// 后序遍历顺序：左右根
// https://leetcode-cn.com/problems/er-cha-sou-suo-shu-de-hou-xu-bian-li-xu-lie-lcof/solution/li-yong-hou-xu-bian-li-de-te-xing-lai-ya-2rka/
function verifyPostorder(postorder: number[]): boolean {
  if (postorder.length < 2) return true
  const rootVal = postorder.pop()!
  let index = 0

  // 找到左右子树的分界点
  while (postorder[index] < rootVal) index++

  const leftTree = postorder.slice(0, index)
  const rightTree = postorder.slice(index)
  const isValidRight = rightTree.every(v => v > rootVal)

  return isValidRight && verifyPostorder(leftTree) && verifyPostorder(rightTree)
}

console.log(verifyPostorder([1, 6, 3, 2, 5]))
console.log(verifyPostorder([1, 3, 2, 6, 5]))

export default 1
