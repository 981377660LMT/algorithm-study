// 其实用map更好O(n)

// 时间复杂度O(n^2)
const intersection = (arr1: number[], arr2: number[]) =>
  [...new Set(arr1)].filter(ele => arr2.includes(ele))

console.log(intersection([1, 2, 3], [1]))
export {}
// 关键思想是存储大数组 遍历小数组
// 1570. 两个稀疏向量的点积
