const LOWERCASE = 'abcdefghijklmnopqrstuvwxyz'
const UPPERCASE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
const DIGITS = '0123456789'
const MOD = 1e9 + 7
const EPS = 1e-8
const DIRS4 = [
  [-1, 0],
  [0, 1],
  [1, 0],
  [0, -1]
]
const DIRS8 = [
  [-1, 0],
  [-1, 1],
  [0, 1],
  [1, 1],
  [1, 0],
  [1, -1],
  [0, -1],
  [-1, -1]
]
const INF = 2e15

// 给你一个下标从 0 开始的整数数组 nums ，它表示一个班级中所有学生在一次考试中的成绩。老师想选出一部分同学组成一个 非空 小组，且这个小组的 实力值 最大，如果这个小组里的学生下标为 i0, i1, i2, ... , ik ，那么这个小组的实力值定义为 nums[i0] * nums[i1] * nums[i2] * ... * nums[ik​] 。
// 请你返回老师创建的小组能得到的最大实力值为多少。
function maxStrength(nums: number[]): number {
  let res = -INF
  enumerateSubset(nums, sub => {
    if (sub.length > 0) {
      let prod = 1
      for (const num of sub) {
        prod *= num
      }
      res = Math.max(res, prod)
    }
  })
  return res
}

function enumerateSubset<T>(nums: ArrayLike<T>, callback: (subset: T[]) => void): void {
  const n = nums.length
  for (let state = 0; state < 1 << n; state++) {
    const cands: T[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) cands.push(nums[j])
    }
    callback(cands)
  }
}
