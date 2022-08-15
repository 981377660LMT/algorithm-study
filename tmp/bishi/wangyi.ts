// /**
//  * 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
//  *
//  *
//  * @param str string字符串
//  * @return string字符串
//  */
// export function getLongestPalindrome(str: string): string {
//   let res = ''
//   const n = str.length

//   for (let i = 0; i < n; i++) {
//     const cand1 = expand(i, i)
//     if (cand1.length > res.length) res = cand1
//     const cand2 = expand(i, i + 1)
//     if (cand2.length > res.length) res = cand2
//   }

//   return res

//   function expand(left: number, right: number): string {
//     while (left >= 0 && right < n && str[left] === str[right]) {
//       left--
//       right++
//     }
//     return str.slice(left + 1, right)
//   }
// }

/**
 * 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
 *
 *
 * @param preStr string字符串 先序遍历序列
 * @param midStr string字符串 中序遍历序列
 * @return string字符串
 */
export function getPostOrderOfTree(preStr: string, midStr: string): string {
  // write code here
  const n = preStr.length
  if (n <= 1) return preStr
  const root = preStr[0]
  const rootIndex = midStr.indexOf(root) // 可以用哈希表+dfs传递left,right来做到O(1)查找
  const leftTree = getPostOrderOfTree(preStr.slice(1, rootIndex + 1), midStr.slice(0, rootIndex))
  const rightTree = getPostOrderOfTree(preStr.slice(rootIndex + 1), midStr.slice(rootIndex + 1))
  return `${leftTree}${rightTree}${root}`
}

if (require.main === module) {
  console.log(getPostOrderOfTree('ACDEFHGB', 'DECAHFBG'))
}
