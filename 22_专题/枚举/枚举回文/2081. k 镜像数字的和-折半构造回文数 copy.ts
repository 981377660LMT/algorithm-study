// 2081. k 镜像数字的和
// https://leetcode.cn/problems/sum-of-k-mirror-numbers/description/
// 2 <= k <= 9
// 1 <= n <= 30

import { enumeratePalindrome } from './enumeratePalindrome'

function kMirror(k: number, n: number): number {
  const res: number[] = []

  enumeratePalindrome(1, 1e9, p => {
    if (res.length >= n) return true
    const num = Number(p)
    const kDigit = num.toString(k)
    if (kDigit === kDigit.split('').reverse().join('')) {
      res.push(num)
    }
  })

  return res.reduce((pre, cur) => pre + cur, 0)
}

if (require.main === module) {
  console.log(kMirror(3, 7))
}
