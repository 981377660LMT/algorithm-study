// 尽量左移
function bisectLeft(arr: number[], target: number): number {
  if (arr.length === 0) return 0

  let l = 0
  let r = arr.length - 1

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const midElement = arr[mid]
    if (midElement < target) {
      l = mid + 1
    } else if (midElement >= target) {
      r = mid - 1
    }
  }

  return l
}

if (require.main === module) {
  // const arr = [7, 7, 7, 7, 7, 7]
  const arr = [-3, -1, 1, 3]
  console.log(bisectLeft(arr, 1))
}

export { bisectLeft }
