// 尽量右移
function bisectRight(arr: number[], target: number): number {
  if (arr.length === 0) return 0

  let l = 0
  let r = arr.length - 1

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const midElement = arr[mid]
    if (midElement <= target) {
      l = mid + 1
    } else if (midElement > target) {
      r = mid - 1
    }
  }

  return l
}

if (require.main === module) {
  const arr = [1, 2, 2, 2, 3, 3, 4, 5, 6, 7]
  console.log(bisectRight(arr, 3))
}

export { bisectRight }
