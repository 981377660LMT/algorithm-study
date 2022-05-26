import { trailingZeroes } from './172. 阶乘后的零'

// f(x) 是 x! 末尾是 0 的数量。
// 我们发现答案只可能是0或5
function preimageSizeFZF(k: number): number {
  let l = 0
  let r = k
  while (l <= r) {
    const mid = (l + r) >> 1
    // 注意这里的优化
    const count = trailingZeroes(mid * 5)
    if (count === k) return 5
    else if (count > k) r = mid - 1
    else l = mid + 1
  }
  return 0
}

console.log(preimageSizeFZF(0))
// 输出：5
// 解释：0!, 1!, 2!, 3!, and 4! 均符合 K = 0 的条件。
