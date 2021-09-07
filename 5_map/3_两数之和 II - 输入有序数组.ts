// 无序数组:需要哈希表记录
// 有序数组：双指针
function twoSum(numbers: number[], target: number): number[] {
  let l = 0
  let r = numbers.length - 1
  while (l < r) {
    const sum = numbers[l] + numbers[r]
    if (sum === target) return [l + 1, r + 1]
    else if (sum > target) r--
    else l++
  }
  return [Infinity, Infinity]
}

console.log(twoSum([2, 7, 11, 15], 9))

export {}
