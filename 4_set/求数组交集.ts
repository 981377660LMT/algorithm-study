const intersection = (arr1: number[], arr2: number[]) => {
  let s1 = new Set(arr1)
  let s2 = new Set(arr2)
  if (s1.size > s2.size) [s1, s2] = [s2, s1]

  const res: number[] = []
  for (const num of s1) {
    s2.has(num) && res.push(num)
  }

  return res
}

console.log(intersection([1, 2, 2, 3, 4, 4], [2, 2, 4, 5, 5, 6, 2000]))
export {}
// 关键思想是存储大数组 遍历小数组
// 1570. 两个稀疏向量的点积
