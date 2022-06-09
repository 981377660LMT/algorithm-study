// 尽量左移
function bisectLeft(arr: ArrayLike<number>, target: number): number {
  if (arr.length === 0) return 0

  let left = 0
  let right = arr.length - 1

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    const midElement = arr[mid]
    if (midElement < target) {
      left = mid + 1
    } else if (midElement >= target) {
      right = mid - 1
    }
  }

  return left
}

// 尽量右移
function bisectRight(arr: ArrayLike<number>, target: number): number {
  if (arr.length === 0) return 0

  let left = 0
  let right = arr.length - 1

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    const midElement = arr[mid]
    if (midElement <= target) {
      left = mid + 1
    } else if (midElement > target) {
      right = mid - 1
    }
  }

  return left
}

function bisectInsort(arr: number[], target: number): void {
  const pos = bisectLeft(arr, target)
  arr.splice(pos, 0, target)
}

if (require.main === module) {
  const arr0 = [-3, -1, 1, 3]
  console.log(bisectLeft(arr0, 1))
  const arr1 = [1, 2, 2, 2, 3, 3, 4, 5, 6, 7]
  console.log(bisectRight(arr1, 3))
  const arr2: number[] = []
  bisectInsort(arr2, 11)
  console.log(arr2)
}

export { bisectLeft, bisectRight, bisectInsort }
