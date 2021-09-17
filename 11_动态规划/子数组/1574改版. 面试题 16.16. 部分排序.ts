// 找出索引m和n，只要将索引区间[m,n]的元素排好序，整个数组就是有序的。
// 若不存在这样的m和n（例如整个数组是有序的），请返回[-1,-1]。
function subSort(array: number[]): number[] {
  let min = Infinity
  let max = -Infinity
  let l = -1
  let r = -1

  // 默认单增 正着比较最大 负着比较最小
  for (let i = 0; i < array.length; i++) {
    if (array[i] < max) r = i
    max = Math.max(max, array[i])
  }

  for (let i = array.length - 1; ~i; i--) {
    if (array[i] > min) l = i
    min = Math.min(min, array[i])
  }

  return [l, r]
}

console.log(subSort([1, 8, 3, 4, 9, 10, 2]))
