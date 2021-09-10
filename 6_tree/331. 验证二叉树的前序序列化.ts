/**
 * @param {string} preorder
 * @return {boolean}
 * 验证它是否是正确的二叉树的前序序列化
 * @summary
 * 类似于括号的题目
 */
const isValidSerialization = function (preorder: string): boolean {
  let balance = 1
  for (const node of preorder.split(',')) {
    if (balance <= 0) return false
    if (node === '#') balance--
    else balance++
  }
  return balance === 0
}

console.log(isValidSerialization('9,3,4,#,#,1,#,#,2,#,6,#,#'))
// 其中 # 代表一个空节点
console.log(isValidSerialization('9,#,#,1'))
