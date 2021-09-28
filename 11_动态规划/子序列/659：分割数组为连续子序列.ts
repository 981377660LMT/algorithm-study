// 1296. 划分数组为连续数字的集合
// 给你一个按升序排序的整数数组 num（可能包含重复数字）， 1 <= nums.length <= 10000
// 请你将它们分割成一个或多个长度至少为 3 的子序列
/**
 * @summary
 * 每遍历一个数，将该数加入能加入的长度最短的序列中，不能加入序列则新建一个序列；然后更新字典
 * @link https://leetcode-cn.com/problems/split-array-into-consecutive-subsequences/solution/zui-jian-dan-de-pythonban-ben-by-semirondo/
 * 直观的解法：打扑克！O(n^2)
 */
function isPossible1(nums: number[]): boolean {
  const n = nums.length
  if (n < 3) return false

  // python中这一段可用for else 简化
  const res: number[][] = []
  for (const num of nums) {
    let shouldInsert = true

    for (let i = res.length - 1; i >= 0; i--) {
      const arr = res[i]
      if (num === arr[arr.length - 1] + 1) {
        arr.push(num)
        shouldInsert = false
        break
      }
    }

    shouldInsert && res.push([num])
  }

  return res.every(arr => arr.length >= 3)
}

// 更好的解法 O(n)
// 贪心
// 应该尽量避免新建短的子序列。
// 使用两个哈希表，第一个哈希表存储数组中的每个数字的剩余次数，
// 第二个哈希表记录哪些元素可以被接到其他子序列后面
type NeedNext = number
type Count = number
function isPossible2(nums: number[]): boolean {
  const n = nums.length
  if (n < 3) return false

  const counter = new Map<number, Count>()
  const need = new Map<NeedNext, Count>()

  for (const num of nums) {
    counter.set(num, (counter.get(num) || 0) + 1)
  }

  for (const num of nums) {
    const count = counter.get(num) || 0
    if (count === 0) continue

    // 先判断 v 是否能接到其他子序列后面
    if (need.has(num) && need.get(num)! > 0) {
      // 可以接到之前的某个序列后面
      counter.set(num, count - 1)
      need.set(num, need.get(num)! - 1)
      need.set(num + 1, (need.get(num + 1) || 0) + 1)
      // 如何避免重复计算counter.get(num + 2)? :python中使用海象运算符:=
    } else if ((counter.get(num + 1) || 0) > 0 && (counter.get(num + 2) || 0) > 0) {
      // 可以有新的序列
      counter.set(num, (counter.get(num) || 0) - 1)
      counter.set(num + 1, (counter.get(num + 1) || 0) - 1)
      counter.set(num + 2, (counter.get(num + 2) || 0) - 1)
      need.set(num + 3, (need.get(num + 3) || 0) + 1)
    } else {
      return false
    }
  }

  return true
}
// console.log(isPossible([1, 2, 3, 3, 4, 4, 5, 5]))
console.log(isPossible2([1, 2, 3, 4, 4, 5]))

export default 1
// 斗地主里面顺子至少要 5 张连续的牌，我们这道题只计算长度最小为 3 的子序列，怎么办？
// 很简单，把我们的 else if 分支修改一下，连续判断 v 之后的连续 5 个元素就行了。

// 如果我想要的不只是一个布尔值，我想要你给我把子序列都打印出来，怎么办？
// 只要修改 need，不仅记录对某个元素的需求个数，而且记录具体是哪些子序列产生的需求：
// eg:记录哪两个子序列需要 6
type NeedByArray = number[][]
const need = new Map<NeedNext, NeedByArray>()
// 我们记录具体子序列的需求也实现了。
