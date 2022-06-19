/**
 原地将该数组中出现的每个零都复写一遍，并将其余的元素向右平移。
 Do not return anything, modify arr in-place instead.
 @summary
 倒着处理
 */
function duplicateZeros(arr: number[]): void {
  const n = arr.length
  let offset = arr.filter(num => num === 0).length

  for (let i = n - 1; ~i; i--) {
    // 此时偏移量为前面0的个数
    if (i + offset < n) {
      arr[i + offset] = arr[i]
    }

    // 偏移量减少
    if (arr[i] === 0) {
      offset--
      if (i + offset < n) {
        arr[i + offset] = 0
      }
    }
  }
}

export {}
