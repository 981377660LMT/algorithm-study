// .chunk() 可以用来按一定长度分割数组。

/**
 * @param {any[]} items
 * @param {number} size
 * @returns {any[][]}
 */
function chunk<T = any>(items: T[], size: number): T[][] {
  if (size === 0) return []

  const res: T[][] = []
  for (let i = 0; i < items.length; i += size) {
    res.push(items.slice(i, i + size))
    // res.push(items.slice(i, Math.min(i + size, items.length)))
  }

  return res
}

console.log(chunk([1, 2, 3, 4, 5], 3))
// [[1, 2, 3], [4, 5]]
