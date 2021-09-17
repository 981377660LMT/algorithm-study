// 给定一个有序整数数组，编写一种方法找出魔术索引
// 魔术索引，满足条件A[i] = i
// 如果没有，则返回-1。若有多个魔术索引，返回索引值最小的一个。

// 如果不存在重复值：二分O(logn)
// 如果存在重复值:直接遍历O(n)
function findMagicIndex(nums: number[]): number {
  return nums.findIndex((v, i) => v === i)
}

console.log(findMagicIndex([0, 2, 3, 4, 5]))
console.log(findMagicIndex([0, 0, 2]))
