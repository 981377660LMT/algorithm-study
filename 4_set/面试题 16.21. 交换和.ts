/**
 * @param {number[]} array1
 * @param {number[]} array2
 * @return {number[]}
 * 请交换一对数值（每个数组中取一个数值），使得两个数组所有元素的和相等。
 * 若无满足条件的数值，返回空数组。
 */
const findSwapValues = function (array1: number[], array2: number[]): number[] {
  const sum = (nums: number[]) => nums.reduce((pre, cur) => pre + cur, 0)
  const sum1 = sum(array1)
  const sum2 = sum(array2)
  let diff = sum1 - sum2
  if (diff & 1) return []
  diff >>= 1

  const set2 = new Set(array2)
  for (const num1 of array1) {
    if (set2.has(num1 - diff)) {
      return [num1, num1 - diff]
    }
  }

  return []
}

console.log(findSwapValues([4, 1, 2, 1, 1, 2], [3, 6, 3, 3]))

export {}
