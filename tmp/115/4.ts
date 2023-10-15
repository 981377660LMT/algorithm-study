export {}

const MOD = 1e9 + 7

// TODO: 1.js卡常技巧
// TODO: 2.为什么排序从小到大比从大到小快(因为从大到小会阻碍缓存命中)
// TODO: 3.多重背包优化的正解
function countSubMultisets(nums: number[], l: number, r: number): number {
  const counter = new Map<number, number>()
  nums.forEach(num => {
    if (counter.has(num)) counter.set(num, counter.get(num)! + 1)
    else counter.set(num, 1)
  })
  const elements: [number, number][] = []
  counter.forEach((count, num) => {
    if (num !== 0) elements.push([num, count])
  })

  // !TODO:为什么
  // 从小到大排序 => 4500ms
  // 从大到小排序 => 9500ms
  // elements.sort((a, b) => b[0] - a[0])
  const m = elements.length
  const sum = nums.reduce((a, b) => a + b, 0)
  const memo = new Int32Array((sum + 1) * m).fill(-1)
  const dfs = (index: number, curSum: number): number => {
    if (curSum > r) return 0
    // @ts-ignore
    if (index === m) return l <= curSum && curSum <= r
    const hash = curSum * m + index
    if (memo[hash] !== -1) return memo[hash]

    let res = 0
    const { 0: curNum, 1: curCount } = elements[index]
    for (let i = 0; i <= curCount; i++) {
      const sum = i * curNum
      if (curSum + sum > r) break
      res += dfs(index + 1, curSum + sum)
      res %= MOD
    }

    memo[hash] = res
    return res
  }

  let res = dfs(0, 0)
  res *= (counter.get(0) || 0) + 1
  return res % MOD
}

// nums = [1,2,2,3], l = 6, r = 6
console.log(countSubMultisets([1, 2, 2, 3], 6, 6))
