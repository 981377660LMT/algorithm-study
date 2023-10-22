// 0 < i < arr.length - 1 条件下，存在 i 使得：
// arr[0] < arr[1] < ... arr[i-1] < arr[i]
// arr[i] > arr[i+1] > ... > arr[arr.length - 1]

/**
 * 有效的山脉数组.
 */
function validMountainArray(
  arr: ArrayLike<number>,
  options?: {
    /** 左侧是否严格递增.默认为true. */
    leftStrict?: boolean
    /** 右侧是否严格递减.默认为true. */
    rightStrict?: boolean
    /** 是否允许某一侧为空.默认为false. */
    allowEmptySide?: boolean
  }
): boolean {
  const { leftStrict = true, rightStrict = true, allowEmptySide = false } = options || {}
  const n = arr.length
  let ptr = 0

  if (leftStrict) {
    while (ptr + 1 < n && arr[ptr] < arr[ptr + 1]) ptr++
  } else {
    while (ptr + 1 < n && arr[ptr] <= arr[ptr + 1]) ptr++
  }

  if (!allowEmptySide && (ptr === 0 || ptr === n - 1)) return false

  if (rightStrict) {
    while (ptr + 1 < n && arr[ptr] > arr[ptr + 1]) ptr++
  } else {
    while (ptr + 1 < n && arr[ptr] >= arr[ptr + 1]) ptr++
  }

  return ptr === n - 1
}

export { validMountainArray, validMountainArray as isMountainArray }

if (require.main === module) {
  console.log(validMountainArray([2, 1])) // false
  console.log(validMountainArray([2, 1], { allowEmptySide: true })) // true
  console.log(validMountainArray([3, 5, 5])) // false
  console.log(validMountainArray([3, 5, 5], { rightStrict: false })) // true
  console.log(validMountainArray([0, 3, 2, 1])) // true
}
