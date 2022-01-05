/**
 * @param {string} s s.length 是 4 的倍数,s 中只含有 'Q', 'W', 'E', 'R' 四种字符
 * @return {number}
 * @description 假如在该字符串中，这四个字符都恰好出现 n/4 次，那么它就是一个「平衡字符串」。
 * 请返回待替换子串的最小可能长度。
 * 注意：必须是替换一段连续的子串
 * @summary 我们窗口内的元素是多出来的元素，我们是把多的元素放到窗口中，那么窗口外的元素每种就肯定都是小于等于N/4的了
 */
const balancedString = function (s: string): number {
  const len = s.length
  // 记录滑动窗口外的对应关系
  const countMap = new Map<string, number>()
  for (const letter of s) {
    countMap.set(letter, (countMap.get(letter) || 0) + 1)
  }
  if (Array.from(countMap.values()).every(v => v === len / 4)) return 0
  let l = 0
  let res = len

  for (let r = 0; r < len; r++) {
    countMap.set(s[r], countMap.get(s[r])! - 1)
    // 满足条件时,找到最佳窗口以保证窗口外的字母数量都小于等于n//4(这样子串就可以将多出来的分给不足的了)
    while (l <= r && Array.from(countMap.values()).every(v => v <= len / 4)) {
      res = Math.min(res, r - l + 1)
      l++
      countMap.set(s[l - 1], countMap.get(s[l - 1])! - 1)
    }
  }
  return res
}

console.log(balancedString('QWER'))
// @ts-ignore
// console.log(~~undefined)
console.log('asas'.split(''))
