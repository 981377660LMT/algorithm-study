// 给定一个含有各种数据的数组，请实现一个deduplicate()来去除重复的元素。
// 请直接修改传入的数组
/**
 * @param {any[]} arr
 */
function deduplicate(arr: any[]) {
  const s = new Set(arr)
  arr.length = 0
  arr.push(...s)
  return arr
}

console.log(deduplicate([1, 2, 3, 2, 1, 34, 2]))
