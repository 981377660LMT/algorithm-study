/**
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

const arr = [1, 0, 2, 3, 0, 4, 5, 0]

duplicateZeros(arr)
console.log(arr)

export {}

function duplicateZeros2(arr: number[]): void {
  for (let i = 0; i < arr.length; i++) {
    if (arr[i] === 0) {
      arr.splice(i, 0, 0)
      arr.pop()
      i++
    }
  }
}
