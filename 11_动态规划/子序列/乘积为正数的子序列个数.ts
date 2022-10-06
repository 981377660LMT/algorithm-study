// 乘积为正数的子序列数 子序列：取还是不取
function foo(nums: number[]): number {
  let [pos, neg] = [0, 0]

  for (const num of nums) {
    if (num === 0) continue
    else if (num > 0) {
      ;[pos, neg] = [pos * 2 + 1, neg * 2]
    } else {
      ;[pos, neg] = [pos + neg, neg + pos + 1]
    }
  }

  return pos
}

console.log(foo([1, 1, -1, -1]))

export {}
