// 其实用map更好O(n)

// 时间复杂度O(n^2)
const intersection = (arr1: number[], arr2: number[]) =>
  [...new Set(arr1)].filter(ele => arr2.includes(ele))

console.log(intersection([1, 2, 3], [1]))
export {}
