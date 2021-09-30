type Index = number
type Diff = number
type Count = number
/**
 *
 * @param nums
 * 区别是子序列了
 */
function numberOfArithmeticSlices(nums: number[]): number {
  // dp[i][diff]表示以 diff 为公差的、且以 nums[i] 为结尾元素的等差数列的个数
  // 但是diff过大 不能用数组 只能用map
  const map = {} as Record<Index, Record<Diff, Count>>
  const n = nums.length
  let res = 0
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < i; j++) {
      const diff = nums[i] - nums[j]
    }
  }

  return res
}

console.log(numberOfArithmeticSlices([2, 4, 6, 8, 10]))

export {}

// class DefaultDict<K, V> {
//   private dict: Record<K, V>
//   private defaultFactory: (...args: any[]) => V

//   constructor(defaultFactory: (...args: any[]) => V) {
//     this.dict = Object.create(null)
//     this.defaultFactory = defaultFactory
//   }

//   get(key: K) {
//     if (!this.map.has(key)) return this.defaultFactory()
//     return this.map.get(key)!
//   }

//   set(key: K, value: V) {
//     if (!this.map.has(key)) this.map.set(key, this.defaultFactory())
//     this.map.set(key, value)
//   }
// }
