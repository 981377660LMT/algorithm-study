/**
 * 有序数组原地去重.
 */
function uniqueInplace(sorted: unknown[]): void {
  let slow = 0
  for (let fast = 0; fast < sorted.length; fast++) {
    if (sorted[fast] !== sorted[slow]) sorted[++slow] = sorted[fast]
  }
  sorted.length = slow + 1
}

/**
 * 有序数组原地去重.
 * @param equals 判断两个元素是否相等的函数.
 */
function uniqueInplace2<E>(sorted: E[], equals: (a: E, b: E) => boolean): void {
  let slow = 0
  for (let fast = 0; fast < sorted.length; fast++) {
    if (!equals(sorted[fast], sorted[slow])) sorted[++slow] = sorted[fast]
  }
  sorted.length = slow + 1
}

export { uniqueInplace, uniqueInplace2 }

if (require.main === module) {
  const arr = [1, 2, 2, 3, 3, 3, 4, 4, 4, 4]
  uniqueInplace(arr)
  console.log(arr) // [1, 2, 3, 4]
}
