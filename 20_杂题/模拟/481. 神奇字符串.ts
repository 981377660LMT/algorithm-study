/**
 * @param {number} n  给你一个整数 n ，返回在神奇字符串 s 的前 n 个数字中 1 的数目。
 * @return {number}
 * @description 
 * 神奇字符串 s 仅由 '1' 和 '2' 组成
   串联字符串中 '1' 和 '2' 的连续出现次数可以生成该字符串。
 * @summary 有点像外观数组那道题
 */
const magicalString = function (n: number): number {
  const sb: number[] = [1]
  let i = 0

  while (sb.length < n) {
    if (sb[i] === 2) sb.push(sb[sb.length - 1])
    sb.push(sb[sb.length - 1] ^ 3) // 1就2 2就1
    i++
  }

  return sb.slice(0, n).filter(v => v === 1).length
}

console.log(magicalString(6))
// 输出：3
// 解释：神奇字符串 s 的前 6 个元素是 “122112”，它包含三个 1，因此返回 3 。

export default 1

// s 的前几个元素是 s = "1221121221221121122……" 。
// 如果将 s 中连续的若干 1 和 2 进行分组，
// 可以得到 "1 22 11 2 1 22 1 22 11 2 11 22 ......" 。
// 每组中 1 或者 2 的出现次数分别是 "1 2 2 1 1 2 1 2 2 1 2 2 ......" 。
// 上面的出现次数正是 s 自身。
