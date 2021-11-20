/**
 Do not return anything, modify arr in-place instead.
 @summary
 倒着处理
 */
function duplicateZeros(arr: number[]): void {
  const n = arr.length
  let zeroCount = arr.filter(num => num === 0).length

  for (let i = n - 1; ~i; i--) {
    // 此时偏移量为前面0的个数
    if (i + zeroCount < n) {
      arr[i + zeroCount] = arr[i]
    }

    // 偏移量减少
    if (arr[i] === 0) {
      zeroCount--
      if (i + zeroCount < n) {
        arr[i + zeroCount] = 0
      }
    }
  }
}

const arr = [1, 0, 2, 3, 0, 4, 5, 0]

duplicateZeros(arr)
console.log(arr)

export {}

function duplicateZeros2(arr: number[]): void {
  const n = arr.length
  let i = 0

  while (i < n) {
    if (arr[i] === 0) {
      arr.splice(i, 0, 0)
      arr.pop()
      i += 2
    } else {
      i++
    }
  }
}
