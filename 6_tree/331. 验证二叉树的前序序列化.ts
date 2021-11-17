/**
 * @param {string} preorder
 * @return {boolean}
 * 验证它是否是正确的二叉树的前序序列化
 * @summary
 * 类似于括号的题目
 * 计算 `出度-入度`
 * 入度：有多少个节点指向它；
   出度：它指向多少个节点。
   根结点的入度为0出度为2，其他非叶子结点的入度为1出度为2，叶子节点入度为1出度为0。
   因为根节点多出来一个出度，所以初始化度为1，一个非叶子节点时度+1，加入一个空节点（叶子节点）时度-1，如果度为0，即达到`出度入度相等`，已经形成一颗二叉树。
 */
const isValidSerialization = function (preorder: string): boolean {
  let degreeDiff = 1

  for (const node of preorder.split(',')) {
    // 出度-入度为0表示完整的树
    if (degreeDiff <= 0) return false
    if (node === '#') degreeDiff--
    else degreeDiff++
  }

  return degreeDiff === 0
}

console.log(isValidSerialization('9,3,4,#,#,1,#,#,2,#,6,#,#'))
// 其中 # 代表一个空节点
console.log(isValidSerialization('9,#,#,1'))
