/**
 * @description 二进制版求子集
 * @description 时间复杂度O(n)
 * 位运算版参考`1863. 找出所有子集的异或总和再求和`
 */
const subsets = (nums: number[]) =>
  Array.from({ length: Math.pow(2, nums.length) }, (_, k) => k)
    .map(num => num.toString(2).padStart(nums.length, '0').split(''))
    .map(item =>
      item
        .map((isNeed, index) => (isNeed === '1' ? nums[index] : Infinity))
        .filter(v => v !== Infinity)
    )

export {}

console.log(subsets([1, 2, 3]))
